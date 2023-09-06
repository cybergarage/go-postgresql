# Changelog

## v1.2.0 (2023-xx-xx)
- Support more data types
  - Timestamp, Year, .... 

## v1.1.0 (2023-xx-xx)
- Update
  - Add ErrorHandler to handler unsupported queries
- Improve
  - Support pgbench

## v1.0.2 (2023-09-05)
- Improve message response compatibility
  -  Update RowDescription response using pg_type table
  -  Update DataRow response

## v1.0.1 (2023-08-30)
- Add support for psql command
  - Add support for SSLRequest startup
  - Fix ParameterStatus responses
- Update go-postgresqld to set verbose logger as default

## v1.0.0 (2023-08-27)
- Fix executor interface
  - StartupHandler
  - Authenticator
  - QueryExecutor
  - BulkExecutor

## v0.9.1 (2023-08-18)
- Add new exeutor interfaces
  - BulkExecutor for COPY messages

## v0.9.0 (2023-07-28)
- Initial public release  
- Support frontend messages
  - Start-up messages
  - Parse
  - Bind
  - Query
  - Terminate
- Add initial executor interfaces
  - CREATE DATABASE
  - CREATE TABLE
  - CREATE INDEX
  - DROP DATABASE
  - DROP TABLE
  - INSERT
  - SELECT (without subqueries)
  - UPDATE
  - DELETE
- Add examples
  - go-postgresqld