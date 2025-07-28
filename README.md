# go-postgresql

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-postgresql)
[![test](https://github.com/cybergarage/go-postgresql/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-postgresql/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-postgresql.svg)](https://pkg.go.dev/github.com/cybergarage/go-postgresql) [![codecov](https://codecov.io/gh/cybergarage/go-postgresql/graph/badge.svg?token=IN4V9KDK69)](https://codecov.io/gh/cybergarage/go-postgresql)

The go-postgresql is a database framework for implementing a [PostgreSQL](https://www.postgresql.org/)-compatible server using Go easily.

## What is the go-postgresql?

The go-postgresql handles [PostgreSQL protocol](https://dev.postgresql.org/doc/dev/postgresql-server/latest/) and interprets the major messages automatically so that all developers can develop PostgreSQL-compatible servers easily. 
 
![](doc/img/framework.png)

Since the go-postgresql handles all startup and system commands automatically, developers can easily implement their PostgreSQL-compatible server only by simply handling DDL (Data Definition Language) and DML (Data Manipulation Language) query commands.

## Table of Contents

- [Getting Started](doc/getting-started.md)

## Examples

- [Examples](doc/examples.md)
	- [go-postgresqld](examples/go-postgresqld) [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/go-postgresqld)](https://hub.docker.com/repository/docker/cybergarage/go-postgresqld/)
	- [go-sqlserver](https://github.com/cybergarage/go-sqlserver) [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/go-sqlserver)](https://hub.docker.com/repository/docker/cybergarage/go-sqlserver/)
	- [PuzzleDB](https://github.com/cybergarage/puzzledb-go) [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/puzzledb)](https://hub.docker.com/repository/docker/cybergarage/puzzledb/)

# Related Projects

The go-postgresql is developed in collaboration with the following Cybergarage projects:

-   [go-logger](https://github.com/cybergarage/go-logger) ![go logger](https://img.shields.io/github/v/tag/cybergarage/go-logger)
-   [go-safecast](https://github.com/cybergarage/go-safecast) ![go safecast](https://img.shields.io/github/v/tag/cybergarage/go-safecast)
-   [go-sqlparser](https://github.com/cybergarage/go-sqlparser) ![go sqlparser](https://img.shields.io/github/v/tag/cybergarage/go-sqlparser)
-   [go-tracing](https://github.com/cybergarage/go-tracing) ![go tracing](https://img.shields.io/github/v/tag/cybergarage/go-tracing)
-   [go-authenticator](https://github.com/cybergarage/go-authenticator) ![go authenticator](https://img.shields.io/github/v/tag/cybergarage/go-authenticator)
-   [go-sasl](https://github.com/cybergarage/go-sasl) ![go sasl](https://img.shields.io/github/v/tag/cybergarage/go-sasl)
-   [go-sqltest](https://github.com/cybergarage/go-sqltest) ![go sqltest](https://img.shields.io/github/v/tag/cybergarage/go-sqltest)
-   [go-pict](https://github.com/cybergarage/go-pict) ![go pict](https://img.shields.io/github/v/tag/cybergarage/go-pict)

## References

- [PostgreSQL](https://www.postgresql.org/)
  - [PostgreSQL: Documentation: Frontend/Backend Protocol](https://www.postgresql.org/docs/current/protocol.html)
