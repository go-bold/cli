# Go Bold CLI

🚀 **Go Bold - framework for Go**

Go Bold CLI is the command-line interface for the Go Bold framework - a rapid development framework that combines developer experience with Go's performance.

## Installation

### From Source
```bash
go install github.com/gobold/cli@latest
```

### From Binary
Download the latest release from GitHub releases.

## Quick Start

```bash
# Create a new Go Bold application
bold new myapp

# Navigate to your project
cd myapp

# Install dependencies
go mod tidy

# Start development server
bold serve
```

## Commands

### `bold new [project-name]`
Creates a new Go Bold application with the complete project structure:

### `bold serve`
Starts the development server with:
- Hot reload (optional)

Options:
- `--port, -p`: Specify port (default: 8000)
- `--hot`: Enable/disable hot reload (default: true)

## Philosophy

### Zero Boilerplate
No 100+ lines of setup code. Get started immediately with sensible defaults.

### Convention over Configuration
Reasonable defaults out of the box. Override only when needed.

### Developer Experience First
Focus on business logic, not infrastructure.

### Rapid Development
Build MVPs in weeks, not months.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
