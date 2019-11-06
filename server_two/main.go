package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"

	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/devries/ngfaas/api"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "50051"
	}

	certificate, err := tls.LoadX509KeyPair("localhost/cert.pem", "localhost/key.pem")
	if err != nil {
		log.Fatalf("could not load server key pair: %s", err)
	}

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("minica.pem")
	if err != nil {
		log.Fatalf("failed to read ca certificate: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append ca certificate to certificate pool")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on port %s", port)

	tlsConfig := &tls.Config{
		ClientAuth:   tls.VerifyClientCertIfGiven,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
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
