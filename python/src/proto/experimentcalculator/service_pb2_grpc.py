# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from proto.experimentcalculator import service_pb2 as proto_dot_experimentcalculator_dot_service__pb2


class ExperimentCalculatorServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.CalcExperiment = channel.unary_unary(
                '/bucketeer.experimentcalculator.ExperimentCalculatorService/CalcExperiment',
                request_serializer=proto_dot_experimentcalculator_dot_service__pb2.BatchCalcRequest.SerializeToString,
                response_deserializer=proto_dot_experimentcalculator_dot_service__pb2.BatchCalcResponse.FromString,
                )


class ExperimentCalculatorServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def CalcExperiment(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ExperimentCalculatorServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'CalcExperiment': grpc.unary_unary_rpc_method_handler(
                    servicer.CalcExperiment,
                    request_deserializer=proto_dot_experimentcalculator_dot_service__pb2.BatchCalcRequest.FromString,
                    response_serializer=proto_dot_experimentcalculator_dot_service__pb2.BatchCalcResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'bucketeer.experimentcalculator.ExperimentCalculatorService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class ExperimentCalculatorService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def CalcExperiment(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.experimentcalculator.ExperimentCalculatorService/CalcExperiment',
            proto_dot_experimentcalculator_dot_service__pb2.BatchCalcRequest.SerializeToString,
            proto_dot_experimentcalculator_dot_service__pb2.BatchCalcResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
