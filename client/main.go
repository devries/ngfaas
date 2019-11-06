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

	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
)

const (
	// address = "ngfaas-j6z4gxi7tq-uc.a.run.app:443"
	address = "localhost:50051"
)

func main() {
	nf := flag.Int64("n", 5, "number of fucks to get")

	flag.Parse()

	certificates, err := tls.LoadX509KeyPair("127.0.0.1/cert.pem", "127.0.0.1/key.pem")
	if err != nil {
		log.Fatalf("Unable to load client certificate and key: %s", err)
	}

	pool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("minica.pem")
	if err != nil {
		log.Fatalf("unable to read ca certificate: %s", err)
	}

	ok := pool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append ca certificate to pool")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{certificates},
		RootCAs:      pool,
	})

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(transportCreds))
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
