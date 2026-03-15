# go-postgresql

`go-postgresql` is a framework for building PostgreSQL-compatible servers in Go. It implements the PostgreSQL wire protocol, handling the complex message flow (startup, authentication, parameter/status, etc.), allowing developers to focus on implementing SQL execution logic.

## Project Overview

- **Core Functionality:** Implements the PostgreSQL Frontend/Backend protocol.
- **Modular Design:** Provides interfaces for DDL (Data Definition), DML (Data Manipulation), TCL (Transaction Control), and system queries.
- **Integrations:** Designed to work with other Cybergarage projects:
  - `go-sqlparser`: For SQL parsing and execution.
  - `go-authenticator`: For password and TLS-based authentication.
  - `go-tracing`: For distributed tracing.
  - `go-logger`: For logging.
  - `go-sqltest`: For testing and verification.

## Key Architecture

- **`postgresql/`**: Contains the core server interfaces and implementation.
  - `Server`: The main server interface.
  - `QueryExecutor`: Interface for handling simple and extended queries.
  - `DDOExecutor` / `DMOExecutor`: Interfaces for DDL and DML operations.
- **`postgresql/protocol/`**: Handles the low-level PostgreSQL wire protocol messages.
- **`postgresql/query/`**: Manages SQL statements, result sets, and data types.
- **`postgresql/auth/`**: Manages authentication mechanisms.
- **`examples/go-postgresqld/`**: A reference implementation of a minimal in-memory PostgreSQL-compatible server.

## Building and Running

The project uses a `Makefile` for common tasks:

- **Build all binaries:**
  ```bash
  make build
  ```
- **Run the example server (go-postgresqld):**
  ```bash
  make run
  ```
- **Run tests with coverage:**
  ```bash
  make test
  ```
- **Run specific tests (e.g., sysbench):**
  ```bash
  make sysbench
  ```
- **Generate certificates for testing:**
  ```bash
  make certs
  ```

## Development Conventions

- **Interface-Driven:** Most core components are defined as interfaces (e.g., `Server`, `QueryExecutor`) to allow for easy extensibility and mocking.
- **Quality Checks:** Before submitting code, ensure it passes linting and vetting:
  ```bash
  make lint
  make vet
  ```
- **Versioning:** The project uses a `version.gen` script to generate `version.go`. Use `make version` to update it.
- **Testing:** Tests are located in the `postgresqltest/` directory and use a combination of Go's `testing` package and integration tests with tools like `psql` and `sysbench`.

## Project Structure

- `bin/`: Utility tools (e.g., `pgpcapdump`).
- `doc/`: Documentation, diagrams, and references.
- `examples/`: Example implementations.
- `postgresql/`: Core framework source code.
- `postgresqltest/`: Comprehensive test suite and test data.
- `scripts/`: Helper scripts for benchmarking and setup.
