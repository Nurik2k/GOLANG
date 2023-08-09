package handlers

import (
	"context"
	pb "grpc_api/proto/protobuf"
)

type Server struct {
	pb.UnimplementedApiServer
}

func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	mockResponse := &pb.GetResponse{
		Resp: []*pb.PostResponse{
			{
				UserId:   in.GetUserId(),
				UserName: in.GetUserName(),
				Password: in.GetPassword(),
			},
		},
	}

	return mockResponse, nil
}

func (s *Server) GetById(ctx context.Context, in *pb.PostIdRequest) (*pb.PostResponse, error) {
	var err error
	if in.PostId == "1" {
		return &pb.PostResponse{
			Id:       "1",
			UserId:   "1",
			UserName: "Nurhzan",
			Password: "123",
		}, nil
	}

	return nil, err
}

func (s *Server) Post(ctx context.Context, in *pb.PostRequest) (*pb.PostResponse, error) {
	newPost := &pb.PostResponse{
		Id:       "1",
		UserId:   "1",
		UserName: in.GetUserName(),
		Password: in.GetPassword(),
	}
	return newPost, nil
}

func (s *Server) Put(ctx context.Context, in *pb.PutRequest) (*pb.PostResponse, error) {
	var err error

	if in.PostId == "1" {
		return &pb.PostResponse{
			Id:       "1",
			UserId:   "1",
			UserName: in.Post.GetPassword(),
			Password: in.Post.GetUserName(),
		}, nil
	}

	return nil, err
}

func (s *Server) Delete(ctx context.Context, in *pb.PostIdRequest) (*pb.EmptyResponse, error) {
	var err error

	if in.PostId == "1" {
		return &pb.EmptyResponse{}, nil
	}

	return nil, err
}
