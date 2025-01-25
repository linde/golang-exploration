
# Overview 

this is a set of explorations of golang, grpc, protobuf, unit/e2e testing in golang, etc. stuff i prob should already know, but dont yet.

to build, you need to first install some things:

```bash
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

go mod tidy
```

next, obtain the protobuf comiler, [protoc](https://grpc.io/docs/protoc-installation/), and use `go generate` to compile the service protobuf and then run tests.

```bash
go generate ./... 
go test -v ./... 
```

# Tring it out

## GRPC 

Then run a server in another terminal or in the background via `&`:

```bash
go run main.go server
```

and try some client calls:

```bash
# default without params
go run main.go client

# specify whom to greet
go run main.go client --name=dolly

# get empahtic and do it a few times
go run main.go client --name=dolly --times=3

# get impatient and specify a timeout to stop any exuberance
go run main.go client --name=dolly --times=100000000 --timeout=1

# forground and ctrl-c your server

```

## REST Gateway 

You can also run a REST gateway to handle REST client calls and proxy them to the gateway.  To do this, run the server with the `--rest [port #]` parameter and use curl or another http client:

```bash

export REST_GW_PORT=8888
go run main.go server --rest=${REST_GW_PORT}  &

curl "http://localhost:${REST_GW_PORT}/v1/helloservice/sayhello?name=dolly&times=1"



```


## clean-up

```bash

# remove the protoc generated go files

find greeter -name *.pb.go  | xargs rm
find greeter -name *.pb.gw.go  | xargs rm

```

Collected Tasks

* TODO: figure out what's up with `greeterserver.ServeListener() failed to listen: <nil>` in the test output
* TODO: clean up output and logger use and use [logging constants](https://pkg.go.dev/log#pkg-constants)
* TODO: think about moving proto to a `proto/` dir
* TODO: generate the openapi v2 schema for the rest gw via this [example](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/examples/internal/proto/examplepb/a_bit_of_everything.proto#L219)
* TODO: consider using docker images for the protobuf binaries, etc, to forgo installing





