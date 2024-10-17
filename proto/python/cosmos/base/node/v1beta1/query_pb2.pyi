from google.api import annotations_pb2 as _annotations_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ConfigRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ConfigResponse(_message.Message):
    __slots__ = ("minimum_gas_price", "pruning_keep_recent", "pruning_interval", "halt_height")
    MINIMUM_GAS_PRICE_FIELD_NUMBER: _ClassVar[int]
    PRUNING_KEEP_RECENT_FIELD_NUMBER: _ClassVar[int]
    PRUNING_INTERVAL_FIELD_NUMBER: _ClassVar[int]
    HALT_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    minimum_gas_price: str
    pruning_keep_recent: str
    pruning_interval: str
    halt_height: int
    def __init__(self, minimum_gas_price: _Optional[str] = ..., pruning_keep_recent: _Optional[str] = ..., pruning_interval: _Optional[str] = ..., halt_height: _Optional[int] = ...) -> None: ...

class StatusRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class StatusResponse(_message.Message):
    __slots__ = ("earliest_store_height", "height", "timestamp", "app_hash", "validator_hash")
    EARLIEST_STORE_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    APP_HASH_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_HASH_FIELD_NUMBER: _ClassVar[int]
    earliest_store_height: int
    height: int
    timestamp: _timestamp_pb2.Timestamp
    app_hash: bytes
    validator_hash: bytes
    def __init__(self, earliest_store_height: _Optional[int] = ..., height: _Optional[int] = ..., timestamp: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., app_hash: _Optional[bytes] = ..., validator_hash: _Optional[bytes] = ...) -> None: ...
