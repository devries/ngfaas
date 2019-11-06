import logging
import argparse

import grpc

import ngfucks_pb2
import ngfucks_pb2_grpc

address = 'localhost:50051'

def get_fucks(stub, nfucks):
    request = ngfucks_pb2.FuckNumber(number=nfucks)
    response = stub.GetFucks(request)

    return response.contents

def main():
    logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(message)s', datefmt='%Y/%m/%d %H:%M:%S')
    parser = argparse.ArgumentParser(description="Get some fucks")
    parser.add_argument('-n', dest='number', default=5, type=int)
    args = parser.parse_args()

    with open('minica.pem','rb') as f:
        root_cert = f.read() 

    with open('127.0.0.1/cert.pem','rb') as f:
        cert = f.read()

    with open('127.0.0.1/key.pem', 'rb') as f:
        private_key = f.read()

    creds = grpc.ssl_channel_credentials(certificate_chain=cert,
            private_key=private_key,
            root_certificates=root_cert)

    with grpc.secure_channel(address, creds) as channel:
        try:
            stub = ngfucks_pb2_grpc.NgFaaSStub(channel)
            r = get_fucks(stub, args.number)
            logging.info("Fucks: "+', '.join(r))
        except grpc.RpcError as e:
            logging.error(f"{e.code()}: {e.details()}")

if __name__=='__main__':
    main()
