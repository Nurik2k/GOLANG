package main

import (
	cl "grpc_api/api_client/handlers"
	pb "grpc_api/proto/protobuf"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial("localhost:5025", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewApiClient(conn)

	// Call the Get method
	cl.Get(client)

	// Call the GetById method
	cl.GetById(client)

	// Call the Post method
	cl.Post(client)

	// Call the Put method
	cl.Put(client)

	// Call the Delete method
	cl.Delete(client)
}
