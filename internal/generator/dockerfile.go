package generator

import (
	"fmt"
	"os"
	"text/template"

	"dockerizer-cli/internal/analyzer"
)

// DockerfileTemplate represents the basic structure for a Dockerfile
const DockerfileTemplate = `{{ if eq .Language "Node.js" }}
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

{{ else if eq .Language "PHP" }}
# Build stage
FROM composer:latest AS builder
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-dev --optimize-autoloader

# Production stage
FROM php:8.2-fpm
WORKDIR /var/www/html

# Install system dependencies
RUN apt-get update && apt-get install -y \
    git \
    curl \
    libpng-dev \
    libonig-dev \
    libxml2-dev \
    zip \
    unzip

# Clear cache
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Install PHP extensions
RUN docker-php-ext-install pdo_mysql mbstring exif pcntl bcmath gd

# Copy composer dependencies
COPY --from=builder /app/vendor ./vendor

# Copy application files
COPY . .

{{ if eq .Framework "laravel" }}
# Set Laravel storage permissions
RUN chown -R www-data:www-data \
    storage \
    bootstrap/cache \
    vendor

# Set Laravel environment
ENV APP_ENV=production
ENV APP_DEBUG=false

# Expose port
EXPOSE {{ index .Ports 0 }}

# Start PHP-FPM
CMD ["php-fpm"]
{{ else }}
EXPOSE 9000
CMD ["php-fpm"]
{{ end }}

{{ else if eq .Language "Python" }}
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

{{ else if eq .Language "Go" }}
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
{{ end }}`

// GenerateDockerfile creates a Dockerfile based on project analysis
func GenerateDockerfile(project *analyzer.ProjectType, outputPath string) error {
	// Validate project configuration
	if project.Language == "" {
		return fmt.Errorf("language not detected")
	}

	// Check if language is supported
	supportedLanguages := map[string]bool{
		"Node.js": true,
		"PHP":     true,
		"Python":  true,
		"Go":      true,
	}

	if !supportedLanguages[project.Language] {
		return fmt.Errorf("unsupported language: %s", project.Language)
	}

	// For PHP, ensure framework is Laravel
	if project.Language == "PHP" && project.Framework != "laravel" {
		return fmt.Errorf("unsupported PHP framework: %s", project.Framework)
	}

	tmpl, err := template.New("dockerfile").Parse(DockerfileTemplate)
	if err != nil {
		return err
	}

	// Create Dockerfile
	file, err := os.Create(outputPath + "/Dockerfile")
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute template with project data
	err = tmpl.Execute(file, project)
	if err != nil {
		// If template execution fails, remove the empty or partial Dockerfile
		os.Remove(outputPath + "/Dockerfile")
		return fmt.Errorf("failed to generate Dockerfile: %w", err)
	}

	// Verify the file is not empty
	fileInfo, err := file.Stat()
	if err != nil || fileInfo.Size() == 0 {
		os.Remove(outputPath + "/Dockerfile")
		return fmt.Errorf("generated Dockerfile is empty, template conditions not met")
	}

	fmt.Println("Successfully generated Dockerfile with multi-stage build support")
	return nil
} 