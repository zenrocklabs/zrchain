from ibc.core.client.v2 import counterparty_pb2 as _counterparty_pb2
from ibc.core.client.v2 import config_pb2 as _config_pb2
from google.api import annotations_pb2 as _annotations_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryCounterpartyInfoRequest(_message.Message):
    __slots__ = ("client_id",)
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    def __init__(self, client_id: _Optional[str] = ...) -> None: ...

class QueryCounterpartyInfoResponse(_message.Message):
    __slots__ = ("counterparty_info",)
    COUNTERPARTY_INFO_FIELD_NUMBER: _ClassVar[int]
    counterparty_info: _counterparty_pb2.CounterpartyInfo
    def __init__(self, counterparty_info: _Optional[_Union[_counterparty_pb2.CounterpartyInfo, _Mapping]] = ...) -> None: ...

class QueryConfigRequest(_message.Message):
    __slots__ = ("client_id",)
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    def __init__(self, client_id: _Optional[str] = ...) -> None: ...

class QueryConfigResponse(_message.Message):
    __slots__ = ("config",)
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    config: _config_pb2.Config
    def __init__(self, config: _Optional[_Union[_config_pb2.Config, _Mapping]] = ...) -> None: ...
