package grpcservice

import (
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type grpcserver struct {
	lis net.Listener
}

func NewServerFromPort(port int) (*grpcserver, error) {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("grpcserver.NewServerFromPort() failed to listen: %v", err)
		return nil, err
	}
	return &grpcserver{lis: lis}, err
}

func NewServerListner(lis net.Listener) *grpcserver {
	return &grpcserver{lis: lis}
}

// TODO migrate all clients to use GetServiceTCPAddr
func (gs grpcserver) GetServicePort() (int, error) {

	addr, err := gs.GetServiceTCPAddr()

	if err != nil {
		return -1, err
	}

	return addr.Port, nil
}

func (gs grpcserver) GetServiceTCPAddr() (*net.TCPAddr, error) {

	addr := gs.lis.Addr()
	switch t := addr.(type) {
	case (*net.TCPAddr):
		return addr.(*net.TCPAddr), nil
	default:
		msg := fmt.Sprintf("grpcserver server Listner address expected net.TCPAddr, was %T", t)
		return nil, errors.New(msg)
	}
}

func (gs grpcserver) Serve(s *grpc.Server) error {

	log.Printf("server listening at %v", gs.lis.Addr())
	if err := s.Serve(gs.lis); err != nil {
		log.Printf("greeterserver.ServeListener() failed to listen: %v", err)
		return err
	}
	return nil
}
