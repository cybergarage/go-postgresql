# Repository Guidelines

## Project Structure & Module Organization
Core library code lives in `postgresql/`, including protocol handling, encoding helpers, and server primitives. Integration and protocol regression tests live in `postgresqltest/`, with focused suites for `protocol/`, `server/`, `sqltest/`, `sysbench/`, `pgbench/`, `ycsb/`, and certificate fixtures under `certs/`. Runnable examples are in `examples/go-postgresqld/`. Utility binaries live in `bin/`, and longer-form documentation is under `doc/`.

## Build, Test, and Development Commands
Use the top-level `Makefile` as the default workflow:

- `make format`: regenerates `postgresql/version.go` and runs `gofmt -s`.
- `make vet`: runs `go vet` after formatting.
- `make lint`: runs `golangci-lint` across `postgresql/`, `postgresqltest/`, and `examples/`.
- `make test`: applies linting, runs the full Go test suite with coverage, and writes `postgresql-cover.out` and `postgresql-cover.html`.
- `make test_only`: runs the test phase without lint.
- `make build`: builds the example daemon.
- `make run`: installs and starts `go-postgresqld` locally.

For targeted work, use standard Go commands such as `go test ./postgresql/... ./postgresqltest/...` or `go test -run TestSQLTestSuite ./postgresqltest/sqltest`.

## Coding Style & Naming Conventions
Follow standard Go formatting with tabs and `gofmt` output. Keep package names lowercase and short; export only API types and functions that belong in the public protocol/server surface. Use descriptive test names like `TestNewServer` and `TestSQLTestSuite`. The repo enforces `golangci-lint` with `gofmt`, `gci`, and `goimports`, so run `make lint` before opening a PR.

## Testing Guidelines
Add `_test.go` files next to the package they validate unless the test depends on the shared harness in `postgresqltest/`. Prefer table-driven unit tests for protocol and encoding behavior; use `postgresqltest/server` helpers for end-to-end server flows. `make test` also fixes certificate permissions before execution, so use it when touching TLS-related tests.

## Commit & Pull Request Guidelines
Recent history favors short imperative subjects, often with conventional prefixes such as `fix(system/schema): ...`, `refactor: ...`, or `style: ...`. Keep commit titles under a single line and scope them when useful. PRs should include a brief problem statement, the main behavioral change, and the exact verification performed (for example, `make lint` and `make test`). Link related issues when applicable; screenshots are only relevant for docs or example output changes.
