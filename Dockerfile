FROM alpine:3.21
RUN apk update && apk add git go

USER root

COPY . /go-postgresql
WORKDIR /go-postgresql

RUN go build -o /go-postgresqld github.com/cybergarage/go-postgresql/examples/go-postgresqld

ENTRYPOINT ["/go-postgresqld"]
