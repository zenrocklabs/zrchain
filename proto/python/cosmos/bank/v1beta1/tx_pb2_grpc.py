# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from cosmos.bank.v1beta1 import tx_pb2 as cosmos_dot_bank_dot_v1beta1_dot_tx__pb2


class MsgStub(object):
    """Msg defines the bank Msg service.
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Send = channel.unary_unary(
                '/cosmos.bank.v1beta1.Msg/Send',
                request_serializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSend.SerializeToString,
                response_deserializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSendResponse.FromString,
                )
        self.MultiSend = channel.unary_unary(
                '/cosmos.bank.v1beta1.Msg/MultiSend',
                request_serializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgMultiSend.SerializeToString,
                response_deserializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgMultiSendResponse.FromString,
                )
        self.UpdateParams = channel.unary_unary(
                '/cosmos.bank.v1beta1.Msg/UpdateParams',
                request_serializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgUpdateParams.SerializeToString,
                response_deserializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgUpdateParamsResponse.FromString,
                )
        self.SetSendEnabled = channel.unary_unary(
                '/cosmos.bank.v1beta1.Msg/SetSendEnabled',
                request_serializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSetSendEnabled.SerializeToString,
                response_deserializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSetSendEnabledResponse.FromString,
                )


class MsgServicer(object):
    """Msg defines the bank Msg service.
    """

    def Send(self, request, context):
        """Send defines a method for sending coins from one account to another account.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def MultiSend(self, request, context):
        """MultiSend defines a method for sending coins from some accounts to other accounts.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def UpdateParams(self, request, context):
        """UpdateParams defines a governance operation for updating the x/bank module parameters.
        The authority is defined in the keeper.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SetSendEnabled(self, request, context):
        """SetSendEnabled is a governance operation for setting the SendEnabled flag
        on any number of Denoms. Only the entries to add or update should be
        included. Entries that already exist in the store, but that aren't
        included in this message, will be left unchanged.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_MsgServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Send': grpc.unary_unary_rpc_method_handler(
                    servicer.Send,
                    request_deserializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSend.FromString,
                    response_serializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSendResponse.SerializeToString,
            ),
            'MultiSend': grpc.unary_unary_rpc_method_handler(
                    servicer.MultiSend,
                    request_deserializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgMultiSend.FromString,
                    response_serializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgMultiSendResponse.SerializeToString,
            ),
            'UpdateParams': grpc.unary_unary_rpc_method_handler(
                    servicer.UpdateParams,
                    request_deserializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgUpdateParams.FromString,
                    response_serializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgUpdateParamsResponse.SerializeToString,
            ),
            'SetSendEnabled': grpc.unary_unary_rpc_method_handler(
                    servicer.SetSendEnabled,
                    request_deserializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSetSendEnabled.FromString,
                    response_serializer=cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSetSendEnabledResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'cosmos.bank.v1beta1.Msg', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Msg(object):
    """Msg defines the bank Msg service.
    """

    @staticmethod
    def Send(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cosmos.bank.v1beta1.Msg/Send',
            cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSend.SerializeToString,
            cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSendResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def MultiSend(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cosmos.bank.v1beta1.Msg/MultiSend',
            cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgMultiSend.SerializeToString,
            cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgMultiSendResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

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
        return grpc.experimental.unary_unary(request, target, '/cosmos.bank.v1beta1.Msg/UpdateParams',
            cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgUpdateParams.SerializeToString,
            cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgUpdateParamsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SetSendEnabled(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cosmos.bank.v1beta1.Msg/SetSendEnabled',
            cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSetSendEnabled.SerializeToString,
            cosmos_dot_bank_dot_v1beta1_dot_tx__pb2.MsgSetSendEnabledResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
