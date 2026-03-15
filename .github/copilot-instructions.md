# Copilot Instructions for go-postgresql

## Project Overview

`go-postgresql` is a **Go framework for building PostgreSQL-compatible servers**. It implements the PostgreSQL Frontend/Backend wire protocol, handling the startup/auth/message-flow machinery so that implementors only need to supply SQL execution logic via executor interfaces.

## Build, Test, and Lint Commands

```bash
make format      # regenerates postgresql/version.go, runs gofmt -s
make vet         # go vet (after format)
make lint        # golangci-lint across postgresql/, postgresqltest/, examples/
make test        # lint → full test suite with coverage → coverage HTML report
make test_only   # test suite only, skips lint
make build       # builds example daemon (go-postgresqld)
make run         # installs and starts go-postgresqld --debug locally
make certs       # regenerates TLS test certificates in postgresqltest/certs/
```

**Run a single test or package:**
```bash
go test -v -run ^TestFoo ./postgresql/...
go test -v -run ^TestFoo ./postgresqltest/server/...
# Sysbench integration test:
go test -v -p 1 -run ^TestSysbench ./postgresqltest/sysbench/...
```

> `make test` also runs `chmod og-rwx` on `postgresqltest/certs/client-key.pem` before executing. If running `go test` directly on TLS-related tests, ensure that permission is set first.

## Architecture

```
postgresql/          # Core framework (the importable library)
  protocol/          # Low-level wire protocol: message types, reader/writer, server loop
  query/             # SQL statement types, ResultSet→Responses conversion, data types
  auth/              # Authentication manager and mechanisms
  errors/            # Shared sentinel errors (ErrNotImplemented, ErrNotExist, etc.)
  encoding/          # Binary encoding helpers
  net/               # Network primitives
  stmt/              # Prepared statement / portal management
  system/            # System catalog query helpers
  server.go          # Public Server interface
  executor.go        # Public executor interfaces (DDO, DMO, TCO, Bulk, System)
  server_impl.go     # Concrete server implementation (unexported)

postgresqltest/      # Integration and regression test suites
  server/            # Test server wrapping examples/go-postgresqld/server
  protocol/          # Wire protocol tests
  sqltest/           # go-sqltest based SQL conformance tests
  sysbench/          # sysbench benchmark integration
  pgbench/           # pgbench integration
  ycsb/              # YCSB benchmark integration
  certs/             # TLS certificates for tests

examples/
  go-postgresqld/    # Minimal reference in-memory PostgreSQL-compatible server
    server/          # Concrete executor implementation to study as a pattern
      executor.go    # Copy/CopyData handlers
      handler.go     # SQL executor implementation
      auth.go        # Auth credential store
      store/         # In-memory table/row store
```

### How the Framework Fits Together

1. **`protocol.Server`** handles the TCP listener, connection lifecycle, and wire framing.
2. **`postgresql.Server`** (the public interface) wraps `protocol.Server` and adds executor registration.
3. Implementors call `postgresql.NewServer()` and register their logic via `Set*Executor` / `SetSQLExecutor`.
4. The framework routes parsed SQL to the appropriate executor method (e.g., `Insert`, `Select`, `Begin`).
5. Executor methods return `(protocol.Responses, error)`. Use helpers in `protocol/responses.go` and `query/resultset.go` to build responses.

```go
// Typical server setup pattern (see examples/go-postgresqld/server/server.go)
srv := postgresql.NewServer()
srv.SetSQLExecutor(myExecutor)   // highest-level: go-sqlparser based
srv.SetQueryExecutor(myExecutor) // or implement DDO/DMO/TCO interfaces directly
srv.SetBulkQueryExecutor(myExecutor)
srv.SetErrorHandler(myExecutor)
srv.Start()
```

## Key Interfaces

| Interface | Purpose |
|-----------|---------|
| `postgresql.QueryExecutor` | Composes `TCOExecutor + DDOExecutor + DMOExecutor` |
| `postgresql.ExQueryExecutor` | Extended ops: `CreateIndex`, `DropIndex`, `Vacuum`, `Truncate` |
| `postgresql.SystemQueryExecutor` | System catalog `SELECT` handling |
| `postgresql.BulkQueryExecutor` | `COPY` / `CopyData` streaming |
| `postgresql.SQLExecutor` | Alias for `go-sqlparser`'s `sql.Executor` (highest-level integration) |
| `postgresql.ErrorHandler` | Custom `ParserError` handler |

All executor methods return `(protocol.Responses, error)`.

## Key Conventions

- **Interface-first design**: core components are defined as interfaces in `postgresql/` and implemented separately. Prefer adding to an interface rather than exposing concrete types.
- **Null/default executors**: the framework ships `NewDefaultQueryExecutor()`, `NewNullBulkExecutor()`, etc. These are used as fallbacks when no executor is registered — don't return `nil` from `New*` constructors.
- **Response construction**: use `protocol.NewInsertCompleteResponsesWith(n)`, `protocol.NewSelectCompleteResponsesWith(n)`, etc., or `query.NewResponseFromResultSet(rs)` for SELECT results. Never construct raw byte slices.
- **Error sentinel values**: use `errors.NewErrNotImplemented(msg)`, `errors.NewErrNotExist(v)`, etc. from `postgresql/errors`. Wrap with `%w` for `errors.Is` compatibility.
- **Linting**: `golangci-lint` is configured with `default: all` minus a specific set of disabled linters (see `.golangci.yaml`). The `gci`, `gofmt`, and `goimports` formatters are enforced. Run `make lint` before any PR.
- **Commit style**: short imperative subject with optional scope, e.g. `fix(system/schema): ...`, `refactor: ...`. Keep to one line.
- **version.go** is generated: edit `postgresql/version.gen`, then run `make format` (which calls `./version.gen`). Never hand-edit `version.go`.
- **Test location**: unit tests live beside the package (`_test.go`); integration tests that need the full server stack go in `postgresqltest/`.
- **Certificate permissions**: `postgresqltest/certs/client-key.pem` must be `og-rwx` before TLS tests run.
