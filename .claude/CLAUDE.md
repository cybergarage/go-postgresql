# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Project Is

`go-postgresql` is a Go framework for building PostgreSQL-compatible servers. It handles the full PostgreSQL wire protocol internally and exposes interfaces for plugging in custom SQL execution, authentication, and error handling. It does **not** implement a full SQL engine — consumers implement the executor interfaces to provide actual data storage and query processing.

## Commands

```bash
# Format + version update
make format

# Static analysis
make vet
make lint

# Full test suite (lint → format → vet → tests with coverage)
make test

# Tests only (skip linting)
make test_only

# Generate TLS test certificates (required before first test run)
make certs

# Build example binaries
make build

# Run example server with debug output
make run
```

Run a single test:
```bash
go test -v -run ^TestName ./postgresqltest/server/
```

Run tests for a specific package:
```bash
go test -v ./postgresql/...
go test -v ./postgresqltest/...
```

CI runs: `make certs && make test`

## Architecture

### Layer Overview

```
Applications (go-postgresqld, go-sqlserver, PuzzleDB)
    ↓ implement executor interfaces
postgresql/ (framework)
    ↓ delegates protocol details
postgresql/protocol/ (wire protocol)
    ↓ uses
postgresql/query/, postgresql/system/, postgresql/auth/, postgresql/stmt/
```

### Core Interfaces (`postgresql/`)

- **`server.go`** — `Server` interface; attach executors via setters (`SetQueryExecutor`, `SetBulkQueryExecutor`, etc.)
- **`executor.go`** — Executor interfaces:
  - `DDOExecutor` / `DDOExExecutor` — CREATE/ALTER/DROP DATABASE/TABLE/INDEX
  - `DMOExecutor` / `DMOExExecutor` — INSERT/SELECT/UPDATE/DELETE/VACUUM/TRUNCATE
  - `TCOExecutor` — BEGIN/COMMIT/ROLLBACK
  - `BulkQueryExecutor` — COPY operations
- **`server_impl.go`** — Concrete server; embeds `protocol.Server`, manages executor instances
- **`server_query_handler.go`** — Routes incoming queries to the appropriate executor
- **`server_startup_handler.go`** — Handles PostgreSQL startup/authentication handshake

### Protocol Layer (`postgresql/protocol/`)

~58 files implementing the PostgreSQL binary wire protocol. Key concepts:
- Message reading/writing: `message.go`, `reader.go`, `writer.go`
- Request handling: `startup.go`, `query.go`, `bind.go`, `execute.go`, `parse.go`
- Response generation: `response.go`, `responses.go`, `data_row.go`
- Connection management: `conn.go`, `conn_impl.go`
- Auth flow: `authentication.go`, `authentication_ok.go`, `authentication_required.go`

### Supporting Packages

| Package | Purpose |
|---|---|
| `postgresql/query/` | SQL parsing wrapper around `go-sqlparser`; statement/resultset abstractions |
| `postgresql/system/` | PostgreSQL system catalogs (`pg_catalog`, `information_schema`), OID mappings, data types |
| `postgresql/auth/` | Authentication manager wrapping `go-authenticator` |
| `postgresql/stmt/` | Prepared statement cache per connection |
| `postgresql/net/` | Network connection abstraction |
| `postgresql/encoding/bytes/` | Binary encoding for integers/floats |
| `postgresql/errors/` | PostgreSQL-compatible error generation |

### Test Infrastructure (`postgresqltest/`)

Tests use a complete in-memory PostgreSQL server implemented in `postgresqltest/server/`. This test server implements all executor interfaces with an in-memory store and is used as a reference implementation for integration tests.

- `postgresqltest/server/` — Test server with in-memory data store
- `postgresqltest/protocol/` — Wire protocol unit tests
- `postgresqltest/sqltest/` — SQL conformance tests via `go-sqltest`
- `postgresqltest/pgbench/`, `sysbench/`, `ycsb/` — Performance benchmark tests

### Reference Implementation (`examples/go-postgresqld/`)

Minimal in-memory PostgreSQL-compatible server demonstrating the framework. See `examples/go-postgresqld/server/` for how to implement the executor interfaces.

## Key External Dependencies

- `github.com/cybergarage/go-sqlparser` — SQL parsing
- `github.com/cybergarage/go-authenticator` — Authentication (SCRAM, MD5, certificate)
- `github.com/cybergarage/go-sqltest` — SQL test suite
- `github.com/jackc/pgx/v5` — PostgreSQL driver used in tests
- `github.com/google/gopacket` — Used in `bin/pgpcapdump` for PCAP analysis
