# Examples

## go-postgresqld

`go-postgresqld` is a minimal, in‑memory PostgreSQL‑compatible server built with `go-postgresql`. It intentionally keeps the surface area small: a single table, no JOIN support—ideal for studying the protocol and executor structure.

```
NAME
  go-postgresqld

SYNOPSIS
  go-postgresqld [OPTIONS]

OPTIONS
  -debug    Enable verbose diagnostic output
  -profile  Enable Go pprof profiling endpoints

EXIT STATUS
  0 on success, non‑zero on failure
```

### Installation

Install the example binary (and any helper tools) into `$GOPATH/bin` via:

```
make install
```

### Profiling

The `-profile` option exposes Go's built‑in pprof HTTP endpoints for runtime inspection.

- [Package pprof](https://golang.org/pkg/net/http/pprof/)
- [Profiling Go Programs](https://blog.golang.org/profiling-go-programs)
