package helloworld

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
	UnimplementedHelloWorldServiceServer
}

func (s *Server) SayHello(ctx context.Context, in *HelloWorld) (*HelloWorld, error) {
	log.Printf("Received: %v", in.GetMessage())
	return &HelloWorld{Message: "Hello " + in.GetMessage()}, nil
}
