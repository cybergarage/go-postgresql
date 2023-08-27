# go-postgresql

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-postgresql)
[![test](https://github.com/cybergarage/go-postgresql/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-postgresql/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-postgresql.svg)](https://pkg.go.dev/github.com/cybergarage/go-postgresql)

The go-postgresql is a database framework for implementing a [PostgreSQL](https://www.postgresql.org/)-compatible server using Go easily.

## What is the go-postgresql?

The go-postgresql handles [PostgreSQL protocol](https://dev.postgresql.org/doc/dev/postgresql-server/latest/) and interprets the major messages automatically so that all developers can develop PostgreSQL-compatible servers easily. Since the go-postgresql handles all startup and system commands automatically, developers can easily implement their PostgreSQL-compatible server only by simply handling DDL (Data Definition Language) and DML (Data Manipulation Language) query commands.
 
![](doc/img/framework.png)

## Table of Contents

- [Getting Started](doc/getting-started.md)

## Examples

- [Examples](doc/examples.md)
  - [go-postgresqld](examples/go-postgresqld)

## References

- [PostgreSQL](https://www.postgresql.org/)
  - [PostgreSQL: Documentation: Developer's Guide](https://www.postgresql.org/docs/16/part-developer.htm)
    - [PostgreSQL: Documentation: command: Overview](https://www.postgresql.org/docs/16/protocol-overview.htm)
    - [PostgreSQL: Documentation: command: Message Flow](https://www.postgresql.org/docs/16/protocol-flow.html)[PostgreSQL: Documentation: command: Message Formats](https://www.postgresql.org/docs/16/protocol-message-formats.html)
    - [PostgreSQL: Documentation: command: Error and Notice Message Fields](https://www.postgresql.org/docs/16/protocol-error-fields.html)
- [PostgreSQL wiki](https://wiki.postgresql.org/wiki/Main_Page)
  - [What information is available to learn PostgreSQL internals? - Developer FAQ - PostgreSQL wiki](https://wiki.postgresql.org/wiki/Developer_FAQ#What_information_is_available_to_learn_PostgreSQL_internals.3F)
- [PostgreSQL Internals Through Picture](https://momjian.us/main/writings/pgsql/internalpics.pdf)
