# Contributing to 1merge

Thank you for your interest in contributing to 1merge! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful and constructive in all interactions. We're committed to maintaining a welcoming environment for all contributors.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- 1Password CLI (`op`) installed and authenticated
- Git

### Development Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/1merge.git
   cd 1merge
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/jstillwa/1merge.git
   ```
4. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests (requires 1Password CLI)
go test -v ./internal/items -run Integration
```

### Building

```bash
# Build the binary
go build -o 1merge

# Run the binary
./1merge --help
```

## Making Changes

### Code Style

- Follow Go conventions and idioms
- Use `go fmt` for formatting
- Keep functions focused and testable
- Add comments for exported functions and complex logic

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):

- **feat**: New feature
- **fix**: Bug fix
- **refactor**: Code restructuring
- **test**: Adding or improving tests
- **docs**: Documentation changes
- **chore**: Build, dependency, or maintenance tasks
- **perf**: Performance improvements

Example:
```
feat(merger): add support for merging custom fields

Improve the merge logic to properly handle custom 1Password fields
while preserving field metadata.
```

### Pull Request Process

1. Ensure all tests pass: `go test ./...`
2. Update documentation if needed
3. Push to your fork and create a Pull Request
4. Provide a clear description of changes
5. Link any related issues
6. Be responsive to review feedback

### Before Submitting

- [ ] Tests pass locally (`go test ./...`)
- [ ] Code is properly formatted (`go fmt ./...`)
- [ ] Commit messages follow Conventional Commits format
- [ ] Changes are documented (code comments, README updates)
- [ ] No hardcoded credentials or sensitive data

## Testing Guidelines

- Write tests for new features
- Update tests when modifying existing functionality
- Aim for clear, descriptive test names
- Test both happy path and error cases
- Use `--dry-run` mode for testing 1Password CLI interactions

## Documentation

- Update README.md if user-facing behavior changes
- Add inline comments for complex logic
- Keep documentation accurate and up to date
- Use clear, concise language

## Licensing

By contributing to 1merge, you agree that your contributions will be licensed under the Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License (CC BY-NC-SA 4.0).

This means:
- Your contributions can be used for non-commercial purposes
- Derivative works must be shared under the same license
- Attribution to you should be maintained

## Questions?

- Check existing issues for answers
- Open a discussion or issue if you have questions
- Review the README.md for project overview

Thank you for contributing to 1merge!
