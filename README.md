Order of Operations:

## No security
- client_one/main.go
- python_client_one/client.py
- server_one/main.go

## Deploy to Cloud Run for Server TLS
- client_two/main.go
- python_client_two/client.py
- server_one/main.go (TLS is external to server)

## Use Private CA for Server TLS
- client_three/main.go
- python_client_three/client.py
- server_two/main.go

## Use Client Certificate AUTH over TLS with Private CA
- client_four/main.go
- python_client_four/client.py
- server_three/main.go

## Use Client Token AUTH over TLS with Private CA
- client_five/main.go
- python?
- server_four/main.go

