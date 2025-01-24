# Dockerizer CLI

Dockerizer CLI is a command-line tool that automatically containerizes software projects. It detects different programming languages and frameworks to generate optimized Dockerfile and docker-compose.yml files.

## Features

- 🔍 Automatic language and framework detection
- 📦 Optimized Dockerfile generation with multi-stage builds
- 🛠 Docker Compose support with database and cache services
- 🚀 Interactive setup process
- 💡 Smart defaults with customization options

## Supported Technologies

All supported technologies are defined in `config/supported_tech.yaml`:

- **Languages**: Node.js, Python, Go, PHP, Ruby, Java
- **Frameworks**: Next.js, React, Angular, Express, NestJS, Django, Flask, FastAPI, Gin, Fiber, Echo, Laravel, Symfony, Rails, Spring Boot
- **Databases**: PostgreSQL, MySQL, MongoDB
- **Cache**: Redis

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
   make install
   ```

### Binary Installation

1. Download the appropriate binary for your operating system from the [Releases](https://github.com/username/dockerizer-cli/releases) page
2. Extract the downloaded archive
3. Add the binary to your PATH or move it to an appropriate location

## Usage

Simply navigate to your project directory and run:

```bash
dockerize init
```

The tool will:
1. Analyze your project structure
2. Detect the programming language and framework
3. Ask for confirmation or alternative selection
4. Offer database integration options
5. Generate optimized Docker files

## Example

```bash
$ cd my-project
$ dockerize init

✔ Detected Node.js as the project language. Is this correct? [Y/n] y
✔ Detected React as the framework. Is this correct? [Y/n] y
✔ Does your project need a database? [Y/n] y
✔ Select database type: PostgreSQL

Successfully generated Docker files!

Next steps:
1. Review the generated files
2. Build and run your container:
   docker-compose up --build
```

## Requirements

- Go 1.16 or higher (for building from source)
- Docker
- Docker Compose

## License

MIT License - See [LICENSE](LICENSE) file for details. 