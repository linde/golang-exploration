package greeterserver

import (
	"fmt"
	"log"
	"net"
)

func ServePort(port int) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("greeterserver.ServePort() failed to listen: %v", err)
		return err
	}
	return ServeListener(lis)
}

func ServeListener(lis net.Listener) error {

	s := NewHelloServer()

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Printf("greeterserver.ServeListener() failed to listen: %v", err)
		return err
	}
	return nil
}
