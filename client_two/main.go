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
	address = "ngfaas.unnecessary.tech:443"
)

func main() {
	nf := flag.Int64("n", 5, "number of fucks to get")

	flag.Parse()

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("unable to load certificate pool: %s", err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(pool, "")))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := api.NewNgFaaSClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := c.GetFucks(ctx, &api.FuckNumber{Number: *nf})
	if err != nil {
		log.Fatalf("could not get fucks: %s", err)
	}

	log.Printf("Fucks: %s", strings.Join(r.Contents, ", "))
}
