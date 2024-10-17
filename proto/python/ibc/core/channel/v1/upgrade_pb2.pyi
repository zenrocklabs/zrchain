from gogoproto import gogo_pb2 as _gogo_pb2
from ibc.core.channel.v1 import channel_pb2 as _channel_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Upgrade(_message.Message):
    __slots__ = ("fields", "timeout", "next_sequence_send")
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_FIELD_NUMBER: _ClassVar[int]
    NEXT_SEQUENCE_SEND_FIELD_NUMBER: _ClassVar[int]
    fields: UpgradeFields
    timeout: _channel_pb2.Timeout
    next_sequence_send: int
    def __init__(self, fields: _Optional[_Union[UpgradeFields, _Mapping]] = ..., timeout: _Optional[_Union[_channel_pb2.Timeout, _Mapping]] = ..., next_sequence_send: _Optional[int] = ...) -> None: ...

class UpgradeFields(_message.Message):
    __slots__ = ("ordering", "connection_hops", "version")
    ORDERING_FIELD_NUMBER: _ClassVar[int]
    CONNECTION_HOPS_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    ordering: _channel_pb2.Order
    connection_hops: _containers.RepeatedScalarFieldContainer[str]
    version: str
    def __init__(self, ordering: _Optional[_Union[_channel_pb2.Order, str]] = ..., connection_hops: _Optional[_Iterable[str]] = ..., version: _Optional[str] = ...) -> None: ...

class ErrorReceipt(_message.Message):
    __slots__ = ("sequence", "message")
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    sequence: int
    message: str
    def __init__(self, sequence: _Optional[int] = ..., message: _Optional[str] = ...) -> None: ...
