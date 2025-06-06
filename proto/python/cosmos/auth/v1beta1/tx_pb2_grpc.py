# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from cosmos.auth.v1beta1 import tx_pb2 as cosmos_dot_auth_dot_v1beta1_dot_tx__pb2


class MsgStub(object):
    """Msg defines the x/auth Msg service.
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.UpdateParams = channel.unary_unary(
                '/cosmos.auth.v1beta1.Msg/UpdateParams',
                request_serializer=cosmos_dot_auth_dot_v1beta1_dot_tx__pb2.MsgUpdateParams.SerializeToString,
                response_deserializer=cosmos_dot_auth_dot_v1beta1_dot_tx__pb2.MsgUpdateParamsResponse.FromString,
                )


class MsgServicer(object):
    """Msg defines the x/auth Msg service.
    """

    def UpdateParams(self, request, context):
        """UpdateParams defines a (governance) operation for updating the x/auth module
        parameters. The authority defaults to the x/gov module account.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_MsgServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'UpdateParams': grpc.unary_unary_rpc_method_handler(
                    servicer.UpdateParams,
                    request_deserializer=cosmos_dot_auth_dot_v1beta1_dot_tx__pb2.MsgUpdateParams.FromString,
                    response_serializer=cosmos_dot_auth_dot_v1beta1_dot_tx__pb2.MsgUpdateParamsResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'cosmos.auth.v1beta1.Msg', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Msg(object):
    """Msg defines the x/auth Msg service.
    """

    @staticmethod
    def UpdateParams(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cosmos.auth.v1beta1.Msg/UpdateParams',
            cosmos_dot_auth_dot_v1beta1_dot_tx__pb2.MsgUpdateParams.SerializeToString,
            cosmos_dot_auth_dot_v1beta1_dot_tx__pb2.MsgUpdateParamsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
