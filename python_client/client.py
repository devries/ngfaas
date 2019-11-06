import logging
import argparse

import grpc

import ngfucks_pb2
import ngfucks_pb2_grpc

#address = 'ngfaas-j6z4gxi7tq-uc.a.run.app:443'
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
        creds = grpc.ssl_channel_credentials(f.read())

    with grpc.secure_channel(address, creds) as channel:
    #with grpc.insecure_channel(address) as channel:
        try:
            stub = ngfucks_pb2_grpc.NgFaaSStub(channel)
            r = get_fucks(stub, args.number)
            logging.info("Fucks: "+', '.join(r))
        except grpc.RpcError as e:
            logging.error(f"{e.code()}: {e.details()}")

if __name__=='__main__':
    main()
