# Changelog

## v1.2.0 (2023-xx-xx)
- Support more data types
  - Timestamp, Year, .... 

## v1.1.2 (2023-xx-xx)
- Improved
  - Example server (go-postgresqld)
  - Support ALTER queries
  - Support pgbench workload

## v1.1.1 (2023-09-12)
- Improved
  - Example server (go-postgresqld)
    - SELECT
      - Supported basic aggregate functions
        - COUNT, SUM, AVG, MIN and MAX
      - Supported basic math functions
        - ABS, CEIL and FLOOR

## v1.1.0 (2023-09-10)
- Improved
  - Executor
    - Added ErrorHandler to handle unsupported queries
  - Query parser
    - Supported
      - SELECT
        - functions
      - ALTER DATABASE
      - ALTER TABLE 
        - ADD, RENAME and DROP COLUMN

## v1.0.2 (2023-09-05)
- Improved
  -  Message response compatibility
     -  Update RowDescription response using pg_type table
     -  Update DataRow response

## v1.0.1 (2023-08-30)
- Improved
  - Added support for psql command
    - Addde support for SSLRequest startup
    - Fixed ParameterStatus responses
  - Updated go-postgresqld to set verbose logger as default

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