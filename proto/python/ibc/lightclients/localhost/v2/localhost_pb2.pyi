from ibc.core.client.v1 import client_pb2 as _client_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ClientState(_message.Message):
    __slots__ = ("latest_height",)
    LATEST_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    latest_height: _client_pb2.Height
    def __init__(self, latest_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...
