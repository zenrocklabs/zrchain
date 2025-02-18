# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from ibc.core.channel.v2 import tx_pb2 as ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2


class MsgStub(object):
    """Msg defines the ibc/channel/v2 Msg service.
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.SendPacket = channel.unary_unary(
                '/ibc.core.channel.v2.Msg/SendPacket',
                request_serializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgSendPacket.SerializeToString,
                response_deserializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgSendPacketResponse.FromString,
                )
        self.RecvPacket = channel.unary_unary(
                '/ibc.core.channel.v2.Msg/RecvPacket',
                request_serializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgRecvPacket.SerializeToString,
                response_deserializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgRecvPacketResponse.FromString,
                )
        self.Timeout = channel.unary_unary(
                '/ibc.core.channel.v2.Msg/Timeout',
                request_serializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgTimeout.SerializeToString,
                response_deserializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgTimeoutResponse.FromString,
                )
        self.Acknowledgement = channel.unary_unary(
                '/ibc.core.channel.v2.Msg/Acknowledgement',
                request_serializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgAcknowledgement.SerializeToString,
                response_deserializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgAcknowledgementResponse.FromString,
                )


class MsgServicer(object):
    """Msg defines the ibc/channel/v2 Msg service.
    """

    def SendPacket(self, request, context):
        """SendPacket defines a rpc handler method for MsgSendPacket.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def RecvPacket(self, request, context):
        """RecvPacket defines a rpc handler method for MsgRecvPacket.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Timeout(self, request, context):
        """Timeout defines a rpc handler method for MsgTimeout.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Acknowledgement(self, request, context):
        """Acknowledgement defines a rpc handler method for MsgAcknowledgement.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_MsgServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'SendPacket': grpc.unary_unary_rpc_method_handler(
                    servicer.SendPacket,
                    request_deserializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgSendPacket.FromString,
                    response_serializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgSendPacketResponse.SerializeToString,
            ),
            'RecvPacket': grpc.unary_unary_rpc_method_handler(
                    servicer.RecvPacket,
                    request_deserializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgRecvPacket.FromString,
                    response_serializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgRecvPacketResponse.SerializeToString,
            ),
            'Timeout': grpc.unary_unary_rpc_method_handler(
                    servicer.Timeout,
                    request_deserializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgTimeout.FromString,
                    response_serializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgTimeoutResponse.SerializeToString,
            ),
            'Acknowledgement': grpc.unary_unary_rpc_method_handler(
                    servicer.Acknowledgement,
                    request_deserializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgAcknowledgement.FromString,
                    response_serializer=ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgAcknowledgementResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'ibc.core.channel.v2.Msg', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Msg(object):
    """Msg defines the ibc/channel/v2 Msg service.
    """

    @staticmethod
    def SendPacket(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/ibc.core.channel.v2.Msg/SendPacket',
            ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgSendPacket.SerializeToString,
            ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgSendPacketResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def RecvPacket(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/ibc.core.channel.v2.Msg/RecvPacket',
            ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgRecvPacket.SerializeToString,
            ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgRecvPacketResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Timeout(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/ibc.core.channel.v2.Msg/Timeout',
            ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgTimeout.SerializeToString,
            ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgTimeoutResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Acknowledgement(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/ibc.core.channel.v2.Msg/Acknowledgement',
            ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgAcknowledgement.SerializeToString,
            ibc_dot_core_dot_channel_dot_v2_dot_tx__pb2.MsgAcknowledgementResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
