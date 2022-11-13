package greeterserver

import (
	"fmt"
	"log"
	"net"
)

// TODO feels like we should be able to combine the async and normal modes of this

// we return the lis out to get the port in cases where it
// was passed in as 0, ie as in tests
func ServePort(port int) (net.Listener, error) {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("greeterserver.ServePort() failed to listen: %v", err)
		return lis, err
	}
	return lis, ServeListener(lis)
}

func ServePortAsync(port int) (net.Listener, func(), error) {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("greeterserver.ServePort() failed to listen: %v", err)
		return lis, func() {}, err
	}
	return lis, ServeListenerAsync(lis), err
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

func ServeListenerAsync(lis net.Listener) (cancel func()) {

	// TODO prob should use a context for this, no?

	flag := make(chan (struct{}))
	cancel = func() { flag <- struct{}{} }
	closed := bool(false)

	s := NewHelloServer()

	// this is our channel close listener, it takes a ref to the server and
	// our flag and Stop()'s it when it get a message on the flag channel
	go func(chan (struct{}), *bool) {
		log.Printf("ServeListenerAsync cancel is closing")
		<-flag
		s.Stop()
		if !closed {
			close(flag)
			closed = true // do we have to syncronize this critical path?
		} else {
			log.Printf("called a second time!")
		}
	}(flag, &closed)

	// this is the routine to serve requests
	go func(cancel func()) {
		defer cancel()
		log.Printf("ServeListenerAsync() listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Printf("ServeListenerAsync() failed to listen: %v", err)
		}
	}(cancel)

	return cancel
}
