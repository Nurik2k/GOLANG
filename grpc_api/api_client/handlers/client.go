package handlers

import (
	"context"
	pb "grpc_api/proto/protobuf"
	"log"
)

func Get(client pb.ApiClient) {

	getRequest := &pb.GetRequest{
		UserId:   "1",
		UserName: "Nurzhan",
		Password: "123",
	}
	getResponse, err := client.Get(context.Background(), getRequest)
	if err != nil {
		log.Fatalf("Error calling Get: %v", err)
	}
	log.Printf("Get Response: %v \n", getResponse)
}

func GetById(client pb.ApiClient) {
	request := &pb.PostIdRequest{
		PostId: "1",
	}

	getByIdResponse, err := client.GetById(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Get by id: %v", getByIdResponse)
}

func Post(client pb.ApiClient) {
	postRequest := &pb.PostRequest{
		UserName: "Nurzhan",
		Password: "123",
	}

	postResponse, err := client.Post(context.Background(), postRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Post Response: %v \n", postResponse)
}

func Put(client pb.ApiClient) {
	putRequest := &pb.PutRequest{
		PostId: "1",
		Post: &pb.Put{
			UserName: "Alisa",
			Password: "1234",
		},
	}

	putResponse, err := client.Put(context.Background(), putRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Put: %v", putResponse)
}

func Delete(client pb.ApiClient) {
	deleteRequest := &pb.PostIdRequest{
		PostId: "1",
	}

	deleteResponse, err := client.Delete(context.Background(), deleteRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Delete: %v", deleteResponse)
}
