package restserver

import (
	"context"
	"fmt"
	"log"
	"myapp/greeter"
	"myapp/grpcservice"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type restgatewayserver struct {
	listener   net.Listener
	gwmux      *runtime.ServeMux
	grpcClient *grpcservice.Clientconn
}

// This is a rest gateway serving on restGatewayPort that proxies
// to the rpc endpoint from rpcAddr. access it with a URL like:
// http://0.0.0.0:{restGatewayPort}/v1/helloservice/sayhello?name=dolly&times=15
func NewRestGateway(restGatewayPort int, rpcAddr *net.TCPAddr) restgatewayserver {

	conn, err := grpcservice.NewNetClientConn(context.Background(), rpcAddr.String())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

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

	return restgatewayserver{listener, gwmux, conn}

}

func (gw restgatewayserver) Serve() error {

	log.Printf("Serving gRPC-Gateway on %s\n", gw.listener.Addr().String())

	servingErr := http.Serve(gw.listener, gw.gwmux)
	log.Fatalf("REST gateway had error serving: %v", servingErr)
	return servingErr
}

func (gw restgatewayserver) Close() {
	if gw.grpcClient != nil {
		gw.grpcClient.Close()
	}
}

func (gw restgatewayserver) GetRestGatewayAddr() net.Addr {
	return gw.listener.Addr()
}
