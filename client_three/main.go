package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/devries/ngfaas/api"
	"google.golang.org/grpc"

	"crypto/x509"
	"google.golang.org/grpc/credentials"
)

const (
	address = "localhost:50051"
)

func main() {
	nf := flag.Int64("n", 5, "number of fucks to get")

	flag.Parse()

	pool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("minica.pem")
	if err != nil {
		log.Fatalf("unable to load ca certificate: %s", err)
	}

	ok := pool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append ca certificate to pool")
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
