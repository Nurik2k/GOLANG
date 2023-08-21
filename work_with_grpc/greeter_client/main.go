package main

import (
	"context"
	"flag"
	"log"
	"time"
	pb "workwithgRPC/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	world = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", world, "Name to greet")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: world})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Greeting: %s", r.GetResult())

}
