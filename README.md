# Next Generation Fucks as a Service: A gRPC sample

In order to learn and understand [gRPC](https://grpc.io), I decided to
implement my [wildly successful FaaS API](https://faas.unnecessary.tech) as a
gRPC service. The service is a simple unary remote procedure call with one
method to get some number of fucks, which are returned by the server. The
[service definition](api/ngfucks.proto) is in the `api` directory and
reproduced below.

```protobuf
syntax = "proto3";

package api;

service NgFaaS {
  rpc GetFucks(FuckNumber) returns (FuckBox) {}
}

message FuckNumber {
  int64 number = 1;
}

message FuckBox {
  repeated string contents = 1;
}
```

This protocol defines a service, `NgFaaS` which has one method, `GetFucks` and
takes a `FuckNumber` message as a request, returning a `FuckBox` message as a
response. The `FuckNumber` message contains one 64 bit integer, while the
`FuckBox` message contains an array of strings.

I write the server in Go and have examples of clients written in both python
and Go. I have a total of five samples. They are:

- Clients and server with no encryption or authentication
    - [python client](python_client_one/client.py) (python_client_one)
    - [Go client](client_one/main.go) (client_one)
    - [Go server](server_one/main.go) (server_one)

- Clients with TLS, and a server behind a Cloud Run proxy
    - [python client](python_client_two/client.py) (python_client_two)
    - [Go client](client_two/main.go) (client_two)
    - [Go server](server_one/main.go) (server_one) - Note this is unchanged from the
      previous sample.

- Clients and server use TLS with a private CA
    - [python client](python_client_three/client.py) (python_client_three)
    - [Go client](client_three/main.go) (client_three)
    - [Go server](server_two/main.go) (server_two)

- Client and server use TLS with a private CA, and clients authenticate with a
  certificate during TLS negotiation
    - [python client](python_client_four/client.py) (python_client_four)
    - [Go client](client_four/main.go) (client_four)
    - [Go server](server_three/main.go) (server_three)

- Client and server use TLS with a private CA, and clients authenticate with a
  token.
    - [python client](python_client_five/client.py) (python_client_five)
    - [Go client](client_five/main.go) (client_five)
    - [Go server](server_four/main.go) (server_four)

## Compiling the code

I use Go modules in this code, which means I use a `go.mod` file. Given this
file, when you build the client and server, the `google.golang.org/grpc`
library should automatically be installed. See the [gRPC](https://grpc.io)
page for more information about this.

To begin, install protocol buffers v3 from the [github project release
page](https://github.com/google/protobuf/releases). You will then need to
install the `protoc` plugin for Go using the command:

```sh
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

Make sure the `protoc-gen-go` binary is within your `PATH`.

For the python side, set up a virtual environment with the command:

```sh
$ python -m venv venv
```

Enter the environment and install the gRPC libraries with the commands:

```sh
$ source venv/bin/activate
$ pip install grpcio
$ pip install grpcio-tools
```

In the root directory of the repository, you can generate the required Go code
with the command:

```sh
$ protoc api/ngfucks.proto -I api/ --go_out=plugins=grpc:api
```

This will write the appropriate go file in the `api` directory where it can be
loaded. For the python client, I found it easier to run the following command
multiple times from within each python client directory:

```sh
$ python -m grpc_tools.protoc -I ../api --python_out=. --grpc_python_out=. ../api/ngfucks.proto
```

This will generate the appropriate python files in each directory where they
can be easily imported by the client software.

The server can be built with the command:

```sh
$ go build -o ngfaas_server server_one/main.go
```

The Go client can be built with the command:

```sh
$ go build -o ngfaas_client client_one/main.go
```

## Running the Server and Clients

From the root directory, run the client using the command

```sh
$ ./ngfaas_server
```

Then in another window or shell you can run the clients. The clients take one
optional argument `-n` followed by a number of fucks to get. By default they
will request 5 fucks. 

To request 20 fucks with the Go client run

```sh
$ ./ngfaas_client -n 20
```

from the root directory of the repository.

To use the python client to request 80 fucks, from the root of the repository
run  the command:

```sh
$ python python_client_one/client.py -n 80
```

You can try very large or negative numbers too in order to see what an error
response looks like.

## Running TLS Clients with Cloud Run

I have put up a server at `ngfaas.unnecessary.tech:443` running on Cloud Run.
The clients are automatically set up to query that server. Note that the Go
client uses the correct system certificates, but the python client is
programmed to look for certificates in `/usr/local/etc/openssl/cert.pem` which
is only available on Linux or if you have installed OpenSSL using
[homebrew](https://brew.sh/) on the Mac. 

## Private Certificate Authority

For all the private certificate authority examples, the clients and servers
are hardcoded to look for a file called `minica.pem` as the root certificate
in the current directory. The clients uses a certificate in
`127.0.0.1/cert.pem` and a key in `127.0.0.1/key.pem` while the server uses a
certificate in `localhost/cert.pem` and a key in `localhost/key.pem`. All the
certificates should be signed by the root certificate.

This can easily be set up using the [minica](https://github.com/jsha/minica)
mini certificate authority, which is also available via
[homebrew](https://brew.sh/) on the Mac. Minica is a good certificate
authority for testing TLS enabled services and clients, but in production I
would recommend using something like [HashiCorp
Vault](https://www.vaultproject.io/) which can create short-lived certificates
on the fly in a secure manner. 

To set up the needed development certificates, run the following commands:

```sh
$ minica -domains localhost
$ minica -ip-addresses 127.0.0.1
```

This will create your root certificate and the server and client certificates.
