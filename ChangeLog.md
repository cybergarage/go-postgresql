# Changelog

## v1.x.0 (2023-xx-xx)
- Update authenticator interface

## v1.1.0 (2023-xx-xx)
- Support more data types
  - Timestamp, Year, .... 

## v1.0.1 (2023-09-xx)
- Support SSLRequest startup
- Fix ParameterStatus message
- go-postgresqld/
  - Update to set logger

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