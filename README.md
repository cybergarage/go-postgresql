# go-postgresql

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-postgresql)
[![test](https://github.com/cybergarage/go-postgresql/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-postgresql/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-postgresql.svg)](https://pkg.go.dev/github.com/cybergarage/go-postgresql) [![codecov](https://codecov.io/gh/cybergarage/go-postgresql/graph/badge.svg?token=IN4V9KDK69)](https://codecov.io/gh/cybergarage/go-postgresql)

go-postgresql is a framework for building [PostgreSQL](https://www.postgresql.org/) compatible servers in Go.

## Overview

go-postgresql implements the [PostgreSQL wire protocol](https://www.postgresql.org/docs/current/protocol.html) and automatically interprets the major message flow (startup, authentication, parameter/status, etc.). This lets you focus on implementing your own DDL (Data Definition Language) and DML (Data Manipulation Language) logic instead of re‑implementing the protocol machinery.

![](doc/img/framework.png)

All startup and system messages are handled for you. Provide implementations for the SQL execution interfaces and you have a functioning PostgreSQL‑compatible server.

## Table of Contents

- [Getting Started](doc/getting-started.md)
- [Examples](doc/examples.md)
- [References](doc/references.md)

## Examples

Representative example projects built with go-postgresql:

- [go-postgresqld](examples/go-postgresqld) – Minimal in‑memory server example. [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/go-postgresqld)](https://hub.docker.com/repository/docker/cybergarage/go-postgresqld/)
- [go-sqlserver](https://github.com/cybergarage/go-sqlserver) – Alternative SQL server implementation. [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/go-sqlserver)](https://hub.docker.com/repository/docker/cybergarage/go-sqlserver/)
- [PuzzleDB](https://github.com/cybergarage/puzzledb-go) – Pluggable multi‑model database. [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/puzzledb)](https://hub.docker.com/repository/docker/cybergarage/puzzledb/)

## Related Projects

go-postgresql is developed in collaboration with other Cybergarage libraries:

- [go-logger](https://github.com/cybergarage/go-logger) ![go logger](https://img.shields.io/github/v/tag/cybergarage/go-logger)
- [go-safecast](https://github.com/cybergarage/go-safecast) ![go safecast](https://img.shields.io/github/v/tag/cybergarage/go-safecast)
- [go-sqlparser](https://github.com/cybergarage/go-sqlparser) ![go sqlparser](https://img.shields.io/github/v/tag/cybergarage/go-sqlparser)
- [go-tracing](https://github.com/cybergarage/go-tracing) ![go tracing](https://img.shields.io/github/v/tag/cybergarage/go-tracing)
- [go-authenticator](https://github.com/cybergarage/go-authenticator) ![go authenticator](https://img.shields.io/github/v/tag/cybergarage/go-authenticator)
- [go-sasl](https://github.com/cybergarage/go-sasl) ![go sasl](https://img.shields.io/github/v/tag/cybergarage/go-sasl)
- [go-sqltest](https://github.com/cybergarage/go-sqltest) ![go sqltest](https://img.shields.io/github/v/tag/cybergarage/go-sqltest)
- [go-pict](https://github.com/cybergarage/go-pict) ![go pict](https://img.shields.io/github/v/tag/cybergarage/go-pict)

## References

- [PostgreSQL](https://www.postgresql.org/)
  - [Frontend/Backend Protocol](https://www.postgresql.org/docs/current/protocol.html)
  - Additional links in [doc/references.md](doc/references.md)
