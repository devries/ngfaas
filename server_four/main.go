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

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)), grpc.UnaryInterceptor(AuthenticationInterceptor))
	api.RegisterNgFaaSServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}

// Let's create a Context Key type with a String() method
type contextKey string

func (c contextKey) String() string {
	return "grpc server context key " + string(c)
}

var contextKeyAuthorized = contextKey("authorized")

// This function returns authorized from the context
func Authorized(ctx context.Context) bool {
	authorized, ok := ctx.Value(contextKeyAuthorized).(bool)
	if !ok {
		return false
	}

	return authorized
}

// This is a UnaryServerInterceptor type which is a function with the signature below.
func AuthenticationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// We can perform a simple logging function here, getting the client address.
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Received Request from %s", p.Addr)
	} else {
		log.Printf("Received Request")
	}

	// We will get the token from the context
	md, ok := metadata.FromIncomingContext(ctx)
	authorized := false
	if !ok {
		log.Printf("Received no metadata")
		authorized = false
	} else {
		tokens := md.Get("authorization")
		if len(tokens) < 1 || tokens[0] != "Bearer HelloWorld" {
			log.Printf("Invalid or Missing token")
			authorized = false
		} else {
			log.Printf("Received valid authorization token")
			authorized = true
		}
	}

	ctx = context.WithValue(ctx, contextKeyAuthorized, authorized)

	h, err := handler(ctx, req)
	return h, err
}

type server struct{}

func (s *server) GetFucks(ctx context.Context, in *api.FuckNumber) (*api.FuckBox, error) {
	// The interceptor places a boolean in the context to let us know if the client is authorized.
	authorized := Authorized(ctx)

	// Next we check to see if it is false.
	if !authorized {
		log.Printf("Unauthorized client")
		return nil, status.Errorf(codes.Unauthenticated, "Invalid or missing authorization token")
	}
	log.Printf("Authorized client")

	// Finally we handle the logic of the server, checking if the number of fucks requested is between
	// 0 and 500 (inclusive), and returning the requested number.
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

	log.Printf("Returning %d fucks", in.Number)

	contentBox := make([]string, in.Number)

	for i := int64(0); i < in.Number; i++ {
		contentBox[i] = "fuck"
	}

	return &api.FuckBox{Contents: contentBox}, nil
}
