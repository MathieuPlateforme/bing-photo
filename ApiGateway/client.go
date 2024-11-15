package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"api-gateway/helloworld"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := helloworld.NewHelloWorldServiceClient(conn)

	message := helloworld.HelloWorld{
		Message: "Hello World",
	}

	response, err := c.SayHello(context.Background(), &message)

	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from server: %s", response.Message)
}
