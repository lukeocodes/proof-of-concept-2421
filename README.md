# Node Transcription Starter

Get started using Deepgram's Transcription with this Node demo app. This starter app demonstrates how to transcribe audio using Deepgram's API.

## Prerequisites

Before you begin, ensure you have:

- Node.js (check package.json for minimum version requirement)
- NPM (Node Package Manager)

## Environment Setup

1. Create a `.env` file by copying the contents of `sample.env`:
   ```bash
   cp sample.env .env
   ```

2. Replace the placeholder in `.env` with your Deepgram API key:
   ```
   DEEPGRAM_API_KEY=your_api_key_here
   ```

   Don't have an API key? Visit [Deepgram's Console](https://console.deepgram.com) to get one.

## Installation

Install all required dependencies:

```bash
npm install
```

## Usage

### Command Line Usage

By default, the app will transcribe our sample audio file (<https://dpgr.am/spacewalk.wav>).

```bash
npm start
```

You can also specify your own audio file:

```bash
# Transcribe from a URL
npm start -- --url=https://example.com/audio.wav

# Transcribe from a local file
npm start -- --path=./path/to/audio.wav
```

The response will be pretty-printed and word-wrapped JSON in this format:
```json
{
  "transcript": "transcription result from Deepgram goes here"
}
```

### Web Server Usage

Start the app as a web server:

```bash
npm start -- --serve
```

Send POST requests to the server with a JSON body containing the audio file URL:

```bash
curl -X POST http://localhost:3000 \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/audio.wav"}'
```

The server will respond with JSON in this format:
```json
{"transcript":"transcription result from Deepgram goes here"}
```

## Contributing

We welcome contributions to make this project better! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on how to get started.

## Security

For information about our security policy and how to report security issues, please see our [Security Policy](SECURITY.md).

## Code of Conduct

This project and everyone participating in it are governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold this code.

## License

This project is licensed under the ISC License - see the [LICENSE](LICENSE) file for details.

## Getting Help

Need assistance? We're here to help!

- Join our [Discord community](https://discord.gg/deepgram) for support
- Found a bug? [Open an issue](https://github.com/deepgram/deepgram-starters/issues/new)
- Have a feature request? [Open an issue](https://github.com/deepgram/deepgram-starters/issues/new)
