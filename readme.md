
this is a set of explorations of golang, grpc, protobuf, unit/e2e testing in golang, etc. stuff i prob should already know, but dont yet.

to build, you need to first install some things:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

next, obtain the [protobuf comiler, protoc](https://grpc.io/docs/protoc-installation/), and use `go generate` to compile the service protobuf and then run tests.

```bash
$ go generate ./... 
$ go test -v ./... 
```

Then run a server in the background or another terminal, optionally in the background with `&`:

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
```

clean-up

```bash
# forground and ctrl-c your server

# remove the protoc generated go files
find greeter -name *.pb.go  | xargs rm
```

Collected Tasks

* TODO: clean up output and logger use and use [logging constants](https://pkg.go.dev/log#pkg-constants)
* TODO: think about moving proto to a `proto/` dir







