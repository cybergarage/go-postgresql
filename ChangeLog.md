# Changelog

## v1.8.0 (2025-xx-xx)
- Improved:
  - Support for more data types.
  - SELECT:
    - Added tests for ORDER BY.
  - Support for SCRAM-SHA-256.

## v1.6.5 (2025-06-07)
- Improved:
  - Extended query executor for prepared statements.
  - Example (go-postgresql):
    - Enhanced support for `sysbench` workload.
    - Enabled `sysbench` tests as default.

## v1.6.4 (2025-05-30)
- Improved:
  - ReadyForQuery returns transaction status correctly.
  - Ignore TLS requests when TLS configuration is null
  - Example (go-postgresql)
    - Refactor to share code with go-mysql 
    - Support for statement commands
      - Math and Aggregate functions

## v1.6.3 (2025-01-18)
- Updated authenticator interface

## v1.6.2 (2024-12-28)
- New Features:
  - Supported for TLS connections.
  - Supported certificate authentication for TLS connection
- Updated authenticator interface
  - Updated password authenticator 

## v1.6.1 (2024-12-11)
- New Features:
  - Supported converting `CREATE INDEX` and `DROP INDEX` commands to `ALTER TABLE`.
- Added:
  - Enabled secondary index tests of `go-sqltest` package.

## v1.6.0 (2024-11-24)
- Introduced protocol server.
- Updated go-sqlparser` package.
- Updated examples to share a common SQL executor with `go-mysql`.

## v1.5.9 (2024-11-16)
- Updated:
  - `go-sqlparser` package.

## v1.5.8 (2024-11-11)
- Updated
  - `go-sqlparser` package.

## v1.5.7 (2024-10-27)
- Added:
  - error package.
- Updated 
  - `go-sqlparser` package.

## v1.5.6 (2024-10-13)
- Improved:
  - internal executor interface.

## v1.5.5 (2024-10-02)
- Improved:
  - protocol package.

## v1.5.4 (2024-06-28)
- Added:
  - Connection manager.

## v1.5.3 (2024-05-25)
- Updated:
  - Authenticator interface.

## v1.5.2 (2024-05-18)
- Updated:
  - TLS settings to allow binary certificates.

## v1.5.1 (2024-05-14)
- Updated:
  - TLS settings to add enable/disable options.

## v1.5.0 (2024-05-12)
- Added:
  - Support for TLS settings in the configuration.

## v1.4.1 (2024-03-20)
- Fixed:
  - lint warnings.

## v1.4.0 (2023-12-03)
- Added an authenticator interface:
  - Added support for password authentication.
- Supported empty queries for ping operations.

## v1.3.5 (2023-11-02)
- Updated:
  - Protocol Messager Reader:
    - Updated to wait until the specified number of bytes are read.

## v1.3.4 (2023-10-27)
- Updated:
  - Connection:
    - Added deadline methods to `Conn`.

## v1.3.3 (2023-10-03)
- Updated:
  - Connection:
    - Added `Conn.UUID()` method to retrieve the UUID of a connection.

## v1.3.2 (2023-09-30)
- Improved:
  - Tracer:
    - Inserted more granular spans.

## v1.3.1 (2023-09-29)
- Improved:
  - Example Server (`go-postgresqld`):
    - ALTER:
      - Added support for ALTER TABLE ADD and DROP COLUMN.

## v1.3.0 (2023-09-28)
- Improved:
  - Query Executor Interfaces:
    - Added a system query executor interface for implementing system tables.
    - Updated the extended query executor interface to handle extended query messages.
    - Updated the bulk executor interface to handle additional COPY messages.
  - Support for new data types:
    - Added support for Timestamps.
  - Example Server (`go-postgresqld`):
    - Added support for `pgbench` workload.

## v1.2.1 (2023-09-19)
- Supported:
  - Sync message.
- Improved:
  - UPDATE:
    - Added support for arithmetic operations.
  - Example Server (`go-postgresqld`):
    - SELECT:
      - Added support for the LIMIT clause.
    - UPDATE:
      - Added support for arithmetic operations.

## v1.2.0 (2023-09-18)
- Added:
  - Query executor interfaces:
    - ALTER, TRUNCATE, VACUUM.
  - Transaction executor interface:
    - BEGIN, COMMIT, ROLLBACK.
- Improved:
  - Updated aggregation functions for empty result sets:
    - COUNT, SUM, AVG, MIN, and MAX.

## v1.1.1 (2023-09-12)
- Improved:
  - Example Server (`go-postgresqld`):
    - SELECT:
      - Added support for basic aggregate functions:
        - COUNT, SUM, AVG, MIN, and MAX.
      - Added support for basic math functions:
        - ABS, CEIL, and FLOOR.

## v1.1.0 (2023-09-10)
- Improved:
  - Executor:
    - Added an ErrorHandler to handle unsupported queries.
  - Query Parser:
    - Added support for:
      - SELECT functions.
      - ALTER DATABASE.
      - ALTER TABLE:
        - ADD, RENAME, and DROP COLUMN.

## v1.0.2 (2023-09-05)
- Improved:
  - Message response compatibility:
    - Updated RowDescription response using `pg_type` table.
    - Updated DataRow response.

## v1.0.1 (2023-08-30)
- Improved:
  - Added support for the `psql` command:
    - Added support for SSLRequest startup.
    - Fixed ParameterStatus responses.
  - Updated `go-postgresqld` to set verbose logging as the default.

## v1.0.0 (2023-08-27)
- Fixed executor interfaces:
  - StartupHandler.
  - Authenticator.
  - QueryExecutor.
  - BulkExecutor.

## v0.9.1 (2023-08-18)
- Added new executor interfaces:
  - BulkExecutor for COPY messages.

## v0.9.0 (2023-07-28)
- Initial public release.
- Supported frontend messages:
  - Start-up messages.
  - Parse.
  - Bind.
  - Query.
  - Terminate.
- Added initial executor interfaces:
  - CREATE DATABASE.
  - CREATE TABLE.
  - CREATE INDEX.
  - DROP DATABASE.
  - DROP TABLE.
  - INSERT.
  - SELECT (without subqueries).
  - UPDATE.
  - DELETE.
- Added examples:
  - `go-postgresqld`.
