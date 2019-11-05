package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/devries/ngfaas/api"
	"google.golang.org/grpc"

	"crypto/x509"
	"google.golang.org/grpc/credentials"
)

const (
	address = "ngfaas-j6z4gxi7tq-uc.a.run.app:443"
	// address = "localhost:50051"
)

func main() {
	nf := flag.Int64("n", 5, "number of fucks to get")

	flag.Parse()

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("unable to load certificate pool: %v", err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(pool, "")))
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewNgFaaSClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := c.GetFucks(ctx, &api.FuckNumber{Number: *nf})
	if err != nil {
		log.Fatalf("could not get fucks: %v", err)
	}

	log.Printf("Fucks: %s", strings.Join(r.Contents, ", "))
}
