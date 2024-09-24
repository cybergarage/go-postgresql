# Getting Started

This section describes how to implement your PostgreSQL-compatible server using the go-postgresql, and see  [Examples](doc/examples.md) about the sample implementation.

## STEP1: Inheriting Server

The go-postgresql offers a core server, [postgresql.Server](../postgresql/server.go), and so inherit the core server in your instance as the following.

```
import (
	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

type MyServer struct {
	*postgresql.Server
}

func NewMyServer() *MyServer {
	return &MyServer{
		Server: postgresql.NewServer(),
	}
}
```

## STEP2: Preparing Query Handler

To handle queries to the your server, prepare a query handler according to [postgresql.QueryExecutor](../postgresql/executor.go) interface.

```
func NewMyServer() *MyServer {
	myserver := &MyServer{
		Server: postgresql.NewServer(),
	}
    myserver.SetQueryExecutor(myserver)
    return myserver
}

func (server *MyServer) Insert(conn *postgresql.Conn, q *query.Insert) (protocol.Responses, error) {
    .....
}

....
```

The go-postgresql offers the stub query executor, [postgresql.BaseExecutor](../postgresql/executor_base.go) which returns a success status for any query requests.
To inheriting the stub executor, you can start to implement only minimum query handle functions such as INSERT and SELECT.

## STEP3: Starting Server 

After implementing the query handler, start your server using  [postgresql.Server::Start()](../postgresql/server.go).

```
server := NewServer()

err := server.Start()
if err != nil {
	t.Error(err)
	return
}
defer server.Stop()

.... 
```
