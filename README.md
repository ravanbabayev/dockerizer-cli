# Dockerizer CLI

Dockerizer CLI is a command-line tool that automatically containerizes software projects. It detects different programming languages and frameworks to generate optimized Dockerfile and docker-compose.yml files.

## Features

- 🔍 Automatic language and framework detection
- 📦 Optimized Dockerfile generation with multi-stage builds
- 🛠 Docker Compose support with database and cache services
- 🚀 Interactive setup process
- 💡 Smart defaults with customization options

## Supported Technologies

All supported technologies are defined in `supported/*.yaml`:

- **Languages**: Node.js, Python, Go, PHP, Ruby, Java
- **Frameworks**: Next.js, React, Angular, Express, NestJS, Django, Flask, FastAPI, Gin, Fiber, Echo, Laravel, Symfony, Rails, Spring Boot
- **Databases**: PostgreSQL, MySQL, MongoDB
- **Cache**: Redis

## Installation

### Windows
```powershell
# Run in PowerShell as Administrator
irm https://raw.githubusercontent.com/USERNAME/dockerizer-cli/main/install.ps1 | iex
```

### Linux/macOS
```bash
curl -fsSL https://raw.githubusercontent.com/USERNAME/dockerizer-cli/main/install.sh | bash
```

### Manual Installation
1. Download the appropriate binary for your operating system from the [Releases](https://github.com/USERNAME/dockerizer-cli/releases) page
2. Extract the archive
3. Add the binary to your PATH

## Usage

Simply navigate to your project directory and run:

```bash
dockerizer init
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
$ dockerizer init

✔ Detected Node.js as the project language. Is this correct? [Y/n] y
✔ Detected React as the framework. Is this correct? [Y/n] y
✔ Default port for React is 3000. Would you like to use a different port? [y/N] n
✔ Does your project need a database? [Y/n] y
✔ Select database type: PostgreSQL
✔ Default port for PostgreSQL is 5432. Would you like to use a different port? [y/N] n

📦 Generating Docker files...
✅ Successfully generated Docker files!

Next steps:
1. Review the generated files
2. Build and run your container:
   docker-compose up --build
```

## Requirements

- Docker
- Docker Compose

## Development

### Building from Source
```bash
# Clone the repository
git clone https://github.com/USERNAME/dockerizer-cli.git
cd dockerizer-cli

# Build
make build

# Install
make install
```

### Running Tests
```bash
make test
```

## License

MIT License - See [LICENSE](LICENSE) file for details. 