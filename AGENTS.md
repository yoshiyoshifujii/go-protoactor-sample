# Repository Guidelines

## Project Structure & Module Organization
- `go.mod` defines the module (`yoshiyoshifujii/go-protoactor-sample`) and Go version.
- `cmd/` contains entry points (currently `cmd/main.go`).
- `internal/domain`, `internal/usecase`, `internal/interface_adaptor` contain Clean Architecture layers.
- Add tests alongside code using `_test.go` files in the same package.

## Build, Test, and Development Commands
- `make build` builds all packages and outputs `bin/go-protoactor-sample`.
- `make run` runs the sample locally.
- `make test` runs all tests across the module.
- `make fmt` formats Go files with `gofmt`.
- `make tidy` tidies module dependencies.

## Coding Style & Naming Conventions
- Follow standard Go formatting with `gofmt` (tabs for indentation; run `gofmt -w .`).
- File names use `snake_case.go` when needed; package names are short, lowercase, no underscores.
- Exported identifiers use `CamelCase`; unexported use `camelCase`.

## Testing Guidelines
- Use Go’s built-in `testing` package.
- Test files must end with `_test.go`; test functions start with `Test` (for example, `TestActorLifecycle`).
- Keep unit tests fast and deterministic; prefer table-driven tests for multiple cases.

## Commit & Pull Request Guidelines
- No explicit commit conventions are established in this repository yet.
- Use clear, imperative commit messages (for example, `Add actor supervisor sample`).
- PRs should include a brief description, the reason for change, and any test commands run.

## Agent-Specific Instructions
- 対話は日本語で行うこと。

## Security & Configuration Tips
- Avoid committing secrets or environment-specific configuration.
- If adding configuration files, document required environment variables in `README.md` or a new `docs/` entry.
