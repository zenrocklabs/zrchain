# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from ibc.core.client.v2 import query_pb2 as ibc_dot_core_dot_client_dot_v2_dot_query__pb2


class QueryStub(object):
    """Query provides defines the gRPC querier service
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.CounterpartyInfo = channel.unary_unary(
                '/ibc.core.client.v2.Query/CounterpartyInfo',
                request_serializer=ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryCounterpartyInfoRequest.SerializeToString,
                response_deserializer=ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryCounterpartyInfoResponse.FromString,
                )
        self.Config = channel.unary_unary(
                '/ibc.core.client.v2.Query/Config',
                request_serializer=ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryConfigRequest.SerializeToString,
                response_deserializer=ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryConfigResponse.FromString,
                )


class QueryServicer(object):
    """Query provides defines the gRPC querier service
    """

    def CounterpartyInfo(self, request, context):
        """CounterpartyInfo queries an IBC light counter party info.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Config(self, request, context):
        """Config queries the IBC client v2 configuration for a given client.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_QueryServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'CounterpartyInfo': grpc.unary_unary_rpc_method_handler(
                    servicer.CounterpartyInfo,
                    request_deserializer=ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryCounterpartyInfoRequest.FromString,
                    response_serializer=ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryCounterpartyInfoResponse.SerializeToString,
            ),
            'Config': grpc.unary_unary_rpc_method_handler(
                    servicer.Config,
                    request_deserializer=ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryConfigRequest.FromString,
                    response_serializer=ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryConfigResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'ibc.core.client.v2.Query', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Query(object):
    """Query provides defines the gRPC querier service
    """

    @staticmethod
    def CounterpartyInfo(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/ibc.core.client.v2.Query/CounterpartyInfo',
            ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryCounterpartyInfoRequest.SerializeToString,
            ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryCounterpartyInfoResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Config(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/ibc.core.client.v2.Query/Config',
            ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryConfigRequest.SerializeToString,
            ibc_dot_core_dot_client_dot_v2_dot_query__pb2.QueryConfigResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
