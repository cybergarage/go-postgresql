FROM golang:1.20-alpine

USER root

COPY . /go-postgresql
WORKDIR /go-postgresql

RUN go build -o /go-postgresqld github.com/cybergarage/go-postgresql/examples/go-postgresqld

ENTRYPOINT ["/go-postgresqld"]
