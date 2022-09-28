import grpc
from proto.environment import service_pb2_grpc


def create_stub(addr, cert_path, service_token_path):
    with open(cert_path, "rb") as f:
        credentials = grpc.ssl_channel_credentials(f.read())
    with open(service_token_path, "rb") as f:
        access_token = f.read().decode("utf-8")
    call_credentials = grpc.access_token_call_credentials(access_token)
    composite_credentials = grpc.composite_channel_credentials(
        credentials, call_credentials
    )
    channel = grpc.secure_channel(addr, composite_credentials)
    return service_pb2_grpc.EnvironmentServiceStub(channel)
