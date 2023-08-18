# Changelog

## v1.1.0 (2023-xx-xx)
- Update authenticator interface

## v1.0.0 (2023-xx-xx)
- Fix executor interface
- Support CopyData message

## v0.9.1 (2023-08-18)
- Add new exeutor interfaces
  - COPY

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