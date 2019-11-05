package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/devries/ngfaas/api"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterNgFaaSServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}

type server struct{}

func (s *server) GetFucks(ctx context.Context, in *api.FuckNumber) (*api.FuckBox, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Received Request from %s", p.Addr)
	} else {
		log.Printf("Received Request")
	}

	if in.Number < 0 {
		retErr := status.Errorf(codes.InvalidArgument, "Negative fucks are not allowed")
		log.Printf("Error: Asked for negative fucks")
		return nil, retErr
	}

	if in.Number > 500 {
		retErr := status.Errorf(codes.InvalidArgument, "%d is too many fucks to give", in.Number)
		log.Printf("Error: Asked for too many fucks")
		return nil, retErr
	}

	contentBox := make([]string, in.Number)

	for i := int64(0); i < in.Number; i++ {
		contentBox[i] = "fuck"
	}

	return &api.FuckBox{Contents: contentBox}, nil
}
