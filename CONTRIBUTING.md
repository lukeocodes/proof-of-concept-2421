# Contributing Guidelines

## Description
We welcome contributions! Before adding new functionality, open an issue first. Bug reports, fixes, and feedback are always appreciated.

## References
- Follow the [GitHub Flow](https://guides.github.com/introduction/flow/index.html)  
- Please review our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing

## Code of Conduct
This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## First-Time Contributors
Check issues labeled <code>beginner</code> and <code>help-wanted</code> to get started.

## Reporting Bugs

### Before Submitting
- Search existing issues and comment if one exists instead of creating a duplicate.

### Submitting
- Use a clear title.  
- List the exact steps to reproduce the issue.  
- Provide examples, links, or code snippets.  
- Describe observed vs. expected behavior.  
- Include screenshots or GIFs (you can use tools like [LICEcap](https://www.cockos.com/licecap/) or [Silentcast](https://github.com/colinkeenan/silentcast)).  
- Mention if the issue is consistent or intermittent and share environment details.

## Suggesting Enhancements

### Before Submitting
- Search existing suggestions and comment on one instead of creating a duplicate.

### Submitting
- Use a clear title.  
- Describe the enhancement step-by-step.  
- Provide examples or code snippets.  
- Explain current vs. expected behavior and its benefits.

## How to Contribute
1. Fork the repository  
2. Create a new branch for your feature or bug fix  
3. Make your changes  
4. Write or update tests as needed  
5. Run the test suite  
6. Submit a pull request

## Pull Requests

### Steps
1. Use the Pull Request template (see [PULL_REQUEST_TEMPLATE/PULL_REQUEST_TEMPLATE.md](PULL_REQUEST_TEMPLATE/PULL_REQUEST_TEMPLATE.md)).  
2. Follow the Code of Conduct.  
3. Ensure all status checks pass before review ([GitHub status checks](https://help.github.com/articles/about-status-checks/)).  

### Review Policy
Reviewers may request additional changes before merging. The PR will be merged once you have the sign-off of at least one maintainer and all conditions are met.

## Development Setup
1. Install Go 1.21 or later  
2. Clone the repository  
3. Run <code>go mod download</code>  
4. Build the project with <code>go build ./cmd/main</code>

## Testing
Run the tests with:
```bash
go test ./...
```

## Questions
If you have questions, feel free to contact the devrel team in any of these formats:  
- GitHub Discussions: https://github.com/orgs/deepgram/discussions  
- Discord: https://discord.gg/deepgram  
- Bluesky: https://bsky.app/profile/deepgram.com  

## Feature Requests
We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:
- Reporting a bug  
- Discussing the current state of the code  
- Submitting a fix  
- Proposing new features  
- Becoming a maintainer  

## License
By contributing, you agree that your contributions will be licensed under the same license as the project.