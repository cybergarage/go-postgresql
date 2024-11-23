# Getting Started

This section describes how to implement your postgresql-compatible server using `go-postgresql`. See [Examples](examples.md) for sample implementations.

## Introduction

Although `go-postgresql` provides various overrideable interfaces for handling postgresql protocol messages, developers typically only need to implement a `go-sqlparser`-based `SQLExecutor` to build your postgresql-compatible server.

![](img/executor.png)

The message executors are implemented by default and generally does not need to be overridden. If authentication is required, an AuthHandler should be implemented. Error handlers are provided for parsing SQL queries (e.g. recovering from parsing errors), but are not detailed in this chapter as they do not normally need to be implemented.

## STEP 1: Inheriting the Server

The `go-postgresql` library provides a core server, [`postgresql.Server`](../postgresql/server.go), which is responsible for handling postgresql protocol messages. To implement your postgresql-compatible server, you should inherit the core postgresql server in your own instance, as shown below:

```go
import (
	"github.com/cybergarage/go-postgresql/postgresql"
)

type MyServer struct {
	postgresql.Server
}

func NewMyServer() *MyServer {
	return &MyServer{
		Server: postgresql.NewServer(),
	}
}
```

The inherited server instance handles postgresql protocol messages. While the default message executors are implemented in the server instance, you will need to provide a custom SQL executor in the next step to handle SQL queries.

## STEP 2: Preparing the Query Handler

The inherited server instance processes postgresql protocol messages and comes with default message executors, but it does not include an SQL executor. 

The SQL executor is defined in the [`go-sqlparser`](https://github.com/cybergarage/go-sqlparser) library as the [`sql.Executor`](https://github.com/cybergarage/go-sqlparser/blob/master/sql/executor.go) interface. It has no dependencies on `go-postgresql` and is also compatible with [`go-mysql`](https://github.com/cybergarage/go-mysql). The `Executor` interface is defined as follows:

```go
type Executor interface {
	Begin(Conn, Begin) error
	Commit(Conn, Commit) error
	Rollback(Conn, Rollback) error	
	CreateDatabase(Conn, CreateDatabase) error
	CreateTable(Conn, CreateTable) error
	AlterDatabase(Conn, AlterDatabase) error
	AlterTable(Conn, AlterTable) error
	DropDatabase(Conn, DropDatabase) error
	DropTable(Conn, DropTable) error
	Insert(Conn, Insert) error
	Select(Conn, Select) (ResultSet, error)
	Update(Conn, Update) (ResultSet, error)
	Delete(Conn, Delete) (ResultSet, error)	
	SystemSelect(Conn, Select) (ResultSet, error)
	Use(Conn, Use) error
}
```

To handle SQL queries on your server, implement a query handler that conforms to the [`sql.Executor`](https://github.com/cybergarage/go-sqlparser/blob/master/sql/executor.go) interface. Then, set the SQL executor on the server instance using [`postgresql.Server::SetSQLExecutor`](../postgresql/server.go) as shown below:

```go
func NewMyServer() *MyServer {
	myServer := &MyServer{
		Server: postgresql.NewServer(),
	}
	myServer.SetSQLExecutor(myServer)
	return myServer
}

func (server *MyServer) Select(conn Conn, stmt Select) (ResultSet, error) {
	// Implement query logic here
	return nil, nil
}
...
```

While it is possible to replace all the default message executors with your own implementations, this guide focuses on implementing only the SQL executor.

## STEP 3: Starting the Server

After implementing the query handler, start your server using [`postgresql.Server::Start()`](../postgresql/server.go):

```go
server := NewMyServer()

err := server.Start()
if err != nil {
	t.Error(err)
	return
}
defer server.Stop()

// Additional logic here
```

To stop the server, use [`postgresql.Server::Stop()`](../postgresql/server.go).

## Next Steps

This guide provides a basic overview of how to implement your postgresql-compatible server using `go-postgresql`. For more detailed examples, see [Examples](examples.md).
