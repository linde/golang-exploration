
this is a set of explorations of golang, grpc, protobuf, unit/e2e testing in golang, etc. stuff i prob should already know, but dont yet.

## TODO, figure out
to get going, you need to first install some things:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

next, obtain the [protobuf comiler, protoc](https://grpc.io/docs/protoc-installation/), and compile the service protobuf with args to use the go plugin from above:

```bash
protoc --go_out=.                               \
        --go_opt=paths=source_relative          \
        --go-grpc_out=.                         \
        --go-grpc_opt=paths=source_relative     \  
        helloservice/helloservice.proto
```

then run a server in the background or another terminal:

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

# get impatient and specify a timeout to stop the exuberance
go run main.go client --name=dolly --times=100000000 --timeout=1
```

clean-up

```bash
# forground and ctrl-c your server

# remove the protoc generated go files
find helloservice -name *.pb.go  | xargs rm
```




