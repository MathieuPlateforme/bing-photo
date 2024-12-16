package main

import (
	"log"
	"net"

	"api-gateway/helloworld"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server listening at :9000")

	s := grpc.NewServer()
	helloworld.RegisterHelloWorldServiceServer(s, &helloworld.Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
