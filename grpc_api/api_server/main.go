package main

import (
	"grpc_api/api_server/handlers"
	pb "grpc_api/proto/protobuf"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()
	srv := &handlers.Server{}

	lis, err := net.Listen("tcp", ":5025")
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterApiServer(s, srv)
	log.Println("Server is listen on port 5025")

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
