# 10x Barry

## An Repo Maintaining AI application

A Go-based tool for processing files using AI capabilities. This tool reads a directory structure, processes files according to `.cursor/rules` markdown component rules, and interacts with AI services.

## Prerequisites

- Go 1.18+ installed
- A valid internet connection

## Installation

To build the binary:

```bash
go mod download
```

### Building for your current platform

```bash
go build -o ./.bin/baz ./cmd/main
```

### Building for Linux (required for GitHub Actions)

```bash
GOOS=linux GOARCH=amd64 go build -o ./.bin/baz.linux ./cmd/main
```

> **Note**: Always build both your platform-specific binary and the Linux binary. The Linux binary is required for GitHub Actions workflows.

## Usage

This starter app can be used in two ways:

1. CLI Usage:
   - Run the tool from the terminal to process your files interactively.
   - For example:  

     ```bash
     ./.bin/baz
     ```

2. Library Usage:
   - Integrate the core functionality of this app within your own Go application.
   - Import the relevant packages from this repo into your code and call the exported functions.

## Configuration

### .bazignore

The `.bazignore` file uses patterns similar to `.gitignore`. Example:

```sh
.cursor/
.git/
*.tmp
```

### Rules

Place your rule files in `.cursor/rules/`. Each rule file should contain instructions for processing specific types of files.

## Development

### Project Structure

```sh
.
├── cmd/
│   └── main/
│       └── main.go
├── pkg/
│   ├── config/
│   ├── filetree/
│   ├── rules/
│   └── tools/
├── .bazignore
├── go.mod
└── README.md
```

## Getting Help

If you have any questions or need help using this application, join our [Discord](https://discord.gg/deepgram) community.

## Reporting Issues and Feature Requests

If you encounter any bugs, or have a feature request, please open an issue in this repository.

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Security

For security concerns, please see our [Security Policy](SECURITY.md).

## Code of Conduct

Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for our code of conduct guidelines.
