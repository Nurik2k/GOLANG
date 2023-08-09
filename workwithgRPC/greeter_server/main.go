package main

import (
	"context"
	"flag"
	"log"
	"net"
	pb "workwithgRPC/helloworld"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "The server port")

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Result: "Hello " + r.GetName()}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
