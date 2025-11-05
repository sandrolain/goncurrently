# Copilot Instructions

## Context

This is a Go package focused on simplifying concurrent programming by providing high-level abstractions and utilities for managing goroutines, channels, and synchronization.

## Instructions

1. **Dependency updates**: Update dependencies in the `go.mod` file to use the latest versions of libraries, maintaining compatibility with existing code.
2. **Test writing**: Add tests for new features and ensure existing tests are updated and working.
3. **Documentation**: Update documentation to reflect changes made and provide clear instructions on how to use the updated package.
4. **Backward compatibility**: Ensure that any changes made do not break existing functionality. If breaking changes are necessary, document them clearly.
5. **Code quality**: Follow Go best practices for code quality, including proper error handling, code organization, and commenting.
6. **Performance**: Optimize the code for performance where possible without sacrificing readability or maintainability.
7. **Security**: Review the code for potential security vulnerabilities and address any issues found.
8. **Versioning**: Follow semantic versioning principles when making
9. **AI**: Do not cite or mention AI tools in the code, comments, or documentation.

## Approach

- Analyze the existing codebase to understand its structure and functionality.
- Identify areas that require updates or improvements based on the instructions provided.
- Be critical and suggest changes that enhance the overall quality of the code while adhering to the instructions.

## Context directory

- Track changes in markdown files into the file `context/changes/{YYYY-MM-DD}.md` with a short description of the change and the date of the change.
- Use a file to preserve context and description of the project into `context/project.md`.
- Use the `/context` directory to store additional context files if needed.

Do not refer to this directory in the code or documentation.

## Code and language

The code and the texts in the code and comments must be written in English.

## Commit messages

For commits use conventional commits format.

## After every change or integration

- run the specific tests to ensure the code works as expected and a minium coverage of 80% is reached.
- assure that inline docs are updated and accurate.
- run `golangci-lint run ./src/...` to ensure no issues are present.

## tests

- benchmark tests must be created in a separated file respecting the unit tests file, and with the suffix `_bench_test.go`.

# GIT

Do not perform git operations, only suggest code changes.