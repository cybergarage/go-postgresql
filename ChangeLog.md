# Changelog

## v1.6.0 (2024-xx-xx)
- Improved
  - Support more data types
    -  .... 
  - SELECT
    - Added tests for ORDER BY

## v1.5.4 (2024-06-28)
- Add connection manager

## v1.5.3 (2024-05-25)
- Updated authenticator interface

## v1.5.2 (2024-05-18)
- Updated 
  - TLS settings to allow binary certificates

## v1.5.1 (2024-05-14)
- Updated 
  - TLS settings to add enable/disable options

## v1.5.0 (2024-05-12)
- New Features
  - Update configuration to support TLS settings

## v1.4.1 (2024-03-20)
- Fixed lint warnings

## v1.4.0 (2023-12-03)
- New Features
  - Added authencator interface
    - Added support for password authentication 
- Improvements
  - Supported empty query for ping queries

## v1.3.5 (2023-11-02)
- Updated
  - Protocol Messager Reader
    - Updated to wait until the specified number of bytes are read

## v1.3.4 (2023-10-27)
- Updated
  - Connection
    - Add deadline methods to Conn

## v1.3.3 (2023-10-03)
- Updated
  - Connection
    - Add Conn.UUID() method to retrieve the UUID of a connection

## v1.3.2 (2023-09-30)
- Improvements
  - Tracer
    - Inserted more fine spans
  
## v1.3.1 (2023-09-29)
- Improvements
  - Example Server (go-postgresqld)
    - ALTER
      - Added support for ALTER TABLE ADD and DROP COLUMN

## v1.3.0 (2023-09-28)
- Improvements
  - Query Executor Interfaces
    - Added a system query executor interface for implementing system tables
    - Updated the extended query executor interface to handle extended query messages
    - Updated the bulk executor interface to handle additional COPY messages
  - Support for New Data Types
    - Timestamp
  - Example Server (go-postgresqld)
    - Added support for pgbench workload

## v1.2.1 (2023-09-19)
- Supported
  - Sync message
- Improved
  - UPDATE
    - Supported arithmetic operations
  - Example server (go-postgresqld)
    - SELECT
      - Supported LIMIT clause
    - UPDATE
      - Supported arithmetic operations

## v1.2.0 (2023-09-18)
- Added
  - Query executor interfaces
    - ALTER, TRUNCATE, VACCUM
  - Transaction executor interface
    - BEGIN, COMMIT, ROLLBACK
- Improved
  - Update aggreation functions for empty result set
    - COUNT, SUM, AVG, MIN and MAX
  
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
