package restservice

import (
	"context"
	"fmt"
	"log"
	"myapp/greeter"
	"myapp/grpcservice"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// This is a rest gateway serving on restGatewayPort that proxies
// to the rpc endpoint from rpcAddr. access it with a URL like:
// http://0.0.0.0:{restGatewayPort}/v1/helloservice/sayhello?name=dolly&times=15
func NewRestGateway(restGatewayPort int, rpcAddr *net.TCPAddr, gwAddrChan chan<- string) {

	// TODO should this not take a cancelable context?
	conn, err := grpcservice.NewNetClientConn(context.Background(), rpcAddr.String())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	gwmux := runtime.NewServeMux()
	err = greeter.RegisterGreeterHandler(context.Background(), gwmux, conn.GetClientConn())
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwTargetStr := fmt.Sprintf(":%d", restGatewayPort)
	listener, err := net.Listen("tcp", gwTargetStr)
	if err != nil {
		log.Fatalf("could not create REST gateway listener: %v", err)
	}

	gwAddrChan <- listener.Addr().String()
	log.Printf("Serving gRPC-Gateway on %s\n", listener.Addr().String())

	servingErr := http.Serve(listener, gwmux)
	if servingErr != nil {
		log.Fatalf("REST gateway had error serving: %v", err)
	}
}
