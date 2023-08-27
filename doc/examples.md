# Examples

## go-postgresqld

The go-postgresqld is a simple PostgreSQL-compatible implementation using the go-postgresqld. The sample implementation is a in-memory PostgreSQL-compatible server, and it supports only a table and do not support any JOIN queries.
```
go-postgresqld is an example of a compatible PostgreSQL server implementation using go-postgresql.
	NAME
	 go-postgresqld

	SYNOPSIS
	 go-postgresqld [OPTIONS]

	OPTIONS
	-debug  : Enable verbose output.
	-profile: Enable profiling.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
```

To install the binary, use the following command. The install command installs the utility programs into `GO_PATH/bin`.

```
make install
```

The profile option enables pprof serves of Go which has the HTTP interface to observe go-postgresqld profile data.

- [The Go Programming Language - Package pprof](https://golang.org/pkg/net/http/pprof/)
- [The Go Blog - Profiling Go Programs](https://blog.golang.org/profiling-go-programs)
