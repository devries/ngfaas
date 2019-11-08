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
	"google.golang.org/grpc/metadata"
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

	// In this case we will read the token directly from the context, but it makes sense
	// to use an Interceptor to validate the token and potentially add role information
	// to the context for use by the service functions.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("Received no metadata")
		return nil, status.Errorf(codes.Unauthenticated, "Unable to access metadata")
	} else {
		tokens := md.Get("authorization")
		if len(tokens) < 1 || tokens[0] != "Bearer HelloWorld" {
			log.Printf("Invalid or Missing token")
			return nil, status.Errorf(codes.Unauthenticated, "Invalid or missing authorization token")
		}
		log.Printf("Received valid authorization token")
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
