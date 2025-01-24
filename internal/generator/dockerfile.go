package generator

import (
	"fmt"
	"os"
	"text/template"

	"dockerizer-cli/internal/analyzer"
)

// DockerfileTemplate represents the basic structure for a Dockerfile
const DockerfileTemplate = `{{ if eq .Language "nodejs" }}
# Build stage
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
{{ if eq .Framework "nextjs" }}
RUN npm run build
{{ else if eq .Framework "react" }}
RUN npm run build
{{ else if eq .Framework "angular" }}
RUN npm run build --prod
{{ end }}

# Production stage
FROM node:18-alpine
WORKDIR /app
{{ if eq .Framework "nextjs" }}
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules
{{ else if eq .Framework "react" }}
COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./
RUN npm install --production
{{ else if eq .Framework "angular" }}
COPY --from=builder /app/dist ./dist
RUN npm install -g serve
{{ else }}
COPY --from=builder /app .
{{ end }}

{{ if eq .Framework "nextjs" }}
CMD ["npm", "start"]
{{ else if eq .Framework "react" }}
CMD ["npm", "start"]
{{ else if eq .Framework "angular" }}
CMD ["serve", "-s", "dist"]
{{ else }}
CMD ["node", "index.js"]
{{ end }}

{{ else if eq .Language "python" }}
# Build stage
FROM python:3.9-slim AS builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --user -r requirements.txt

# Production stage
FROM python:3.9-slim
WORKDIR /app
COPY --from=builder /root/.local /root/.local
COPY . .
ENV PATH=/root/.local/bin:$PATH

{{ if eq .Framework "django" }}
EXPOSE 8000
CMD ["python", "manage.py", "runserver", "0.0.0.0:8000"]
{{ else if eq .Framework "flask" }}
EXPOSE 5000
CMD ["flask", "run", "--host=0.0.0.0"]
{{ else if eq .Framework "fastapi" }}
EXPOSE 8000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
{{ else }}
CMD ["python", "app.py"]
{{ end }}

{{ else if eq .Language "go" }}
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
{{ range .Ports }}
EXPOSE {{ . }}
{{ end }}
CMD ["./main"]

{{ else if eq .Language "java" }}
# Build stage
FROM maven:3.8.4-openjdk-17-slim AS builder
WORKDIR /app
COPY pom.xml .
COPY src ./src
RUN mvn clean package -DskipTests

# Production stage
FROM openjdk:17-slim
WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar
EXPOSE 8080
CMD ["java", "-jar", "app.jar"]

{{ else if eq .Language "php" }}
# Build stage
FROM composer:latest AS builder
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-dev --optimize-autoloader

# Production stage
FROM php:8.2-fpm
WORKDIR /var/www/html
COPY --from=builder /app/vendor ./vendor
COPY . .
{{ if eq .Framework "laravel" }}
RUN chown -R www-data:www-data storage bootstrap/cache
EXPOSE 9000
CMD ["php-fpm"]
{{ else }}
EXPOSE 9000
CMD ["php-fpm"]
{{ end }}

{{ else if eq .Language "ruby" }}
# Build stage
FROM ruby:3.2-alpine AS builder
WORKDIR /app
COPY Gemfile* ./
RUN apk add --no-cache build-base
RUN bundle install --without development test

# Production stage
FROM ruby:3.2-alpine
WORKDIR /app
COPY --from=builder /usr/local/bundle /usr/local/bundle
COPY . .
{{ if eq .Framework "rails" }}
EXPOSE 3000
CMD ["rails", "server", "-b", "0.0.0.0"]
{{ else }}
CMD ["ruby", "app.rb"]
{{ end }}

{{ end }}`

// GenerateDockerfile creates a Dockerfile based on project analysis
func GenerateDockerfile(project *analyzer.ProjectType, outputPath string) error {
	tmpl, err := template.New("dockerfile").Parse(DockerfileTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath + "/Dockerfile")
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, project)
	if err != nil {
		return err
	}

	fmt.Println("Successfully generated Dockerfile with multi-stage build support")
	return nil
} 