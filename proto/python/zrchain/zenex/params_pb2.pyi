from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Params(_message.Message):
    __slots__ = ("btc_proxy_address", "minimum_satoshis")
    BTC_PROXY_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    MINIMUM_SATOSHIS_FIELD_NUMBER: _ClassVar[int]
    btc_proxy_address: str
    minimum_satoshis: int
    def __init__(self, btc_proxy_address: _Optional[str] = ..., minimum_satoshis: _Optional[int] = ...) -> None: ...
