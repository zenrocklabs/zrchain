from ibc.core.client.v2 import counterparty_pb2 as _counterparty_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisCounterpartyInfo(_message.Message):
    __slots__ = ("client_id", "counterparty_info")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_INFO_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    counterparty_info: _counterparty_pb2.CounterpartyInfo
    def __init__(self, client_id: _Optional[str] = ..., counterparty_info: _Optional[_Union[_counterparty_pb2.CounterpartyInfo, _Mapping]] = ...) -> None: ...

class GenesisState(_message.Message):
    __slots__ = ("counterparty_infos",)
    COUNTERPARTY_INFOS_FIELD_NUMBER: _ClassVar[int]
    counterparty_infos: _containers.RepeatedCompositeFieldContainer[GenesisCounterpartyInfo]
    def __init__(self, counterparty_infos: _Optional[_Iterable[_Union[GenesisCounterpartyInfo, _Mapping]]] = ...) -> None: ...
