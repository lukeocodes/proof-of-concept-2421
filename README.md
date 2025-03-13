# Deepgram Luke

## An Repo Maintaining AI application

A Go-based tool for processing files using AI capabilities. This tool reads a directory structure, processes files according to `.cursor/rules` markdown component rules, and interacts with AI services.

## Installation

```bash
go mod download
go build -o ./.bin/processor ./cmd/main
```

## Usage

1. Create a `.lukeignore` file to specify which files/directories to ignore
2. Place your rules in the `cursor/rules` directory
3. Run the tool:

```bash
./.bin/processor
```

## Configuration

### .lukeignore

The `.lukeignore` file uses patterns similar to `.gitignore`. Example:

```
.cursor/
.git/
*.tmp
```

### Rules

Place your rule files in `cursor/rules/`. Each rule file should contain instructions for processing specific types of files.

## Development

### Project Structure

```
.
├── cmd/
│   └── main/
│       └── main.go
├── pkg/
│   ├── config/
│   ├── filetree/
│   ├── rules/
│   └── tools/
├── .lukeignore
├── go.mod
└── README.md
```

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Security

For security concerns, please see our [Security Policy](SECURITY.md).
