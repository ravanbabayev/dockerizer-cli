# Dockerizer CLI

Dockerizer CLI is a command-line tool that automatically containerizes software projects. It detects different programming languages and frameworks to generate optimized Dockerfile and docker-compose.yml files.

## Features

- 🔍 Automatic language and framework detection
- 📦 Optimized Dockerfile generation with multi-stage builds
- 🛠 Docker Compose support with database and cache services
- 🚀 Easy to use and quick to start

## Supported Languages and Frameworks

- **Node.js**: Next.js, React, Angular, Express, NestJS
- **Python**: Django, Flask, FastAPI
- **Go**: Gin, Fiber, Echo
- **PHP**: Laravel, Symfony
- **Ruby**: Rails
- **Java**: Spring Boot

## Installation

### From Source (Using Go)

1. Install Go (1.16 or higher):
   - https://golang.org/dl/

2. Clone the repository:
   ```bash
   git clone https://github.com/username/dockerizer-cli.git
   cd dockerizer-cli
   ```

3. Build and install:
   ```bash
   go install
   ```

### Binary Installation

1. Download the appropriate binary for your operating system from the [Releases](https://github.com/ravanbabayev/dockerizer-cli/releases) page
2. Extract the downloaded archive
3. Add the binary to your PATH or move it to an appropriate location

## Usage

1. Analyze your project:
   ```bash
   dockerizer init --path /project/directory
   ```

2. Generate Docker files:
   ```bash
   dockerizer generate --path /project/directory --output /output/directory
   ```

3. Clean up generated files:
   ```bash
   dockerizer clean --output /output/directory
   ```

## Example Usage

For a React project:
```bash
# Navigate to project directory
cd my-react-app

# Analyze the project
dockerizer init

# Generate Docker files
dockerizer generate

# Start the container
docker-compose up
```

## Requirements

- Go 1.16 or higher
- Docker
- Docker Compose

## License

MIT License - See [LICENSE](LICENSE) file for details. 