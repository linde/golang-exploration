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

// TODO should we return the cancel, err?
func ServeListenerAsync(lis net.Listener) func() {

	// TODO prob should use a context for this, no?

	flag := make(chan (struct{}))
	cancel := func() { flag <- struct{}{} }

	s := NewHelloServer()

	go func(chan (struct{})) {
		<-flag
		s.Stop()
		close(flag)
	}(flag)

	go func(cancel func()) {
		defer cancel()
		log.Printf("ServeListenerAsync() listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Printf("ServeListenerAsync() failed to listen: %v", err)
		}
	}(cancel)

	// TODO should be able to get error as an return from the go rountine
	return cancel
}
