from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgInit(_message.Message):
    __slots__ = ("initial_value",)
    INITIAL_VALUE_FIELD_NUMBER: _ClassVar[int]
    initial_value: int
    def __init__(self, initial_value: _Optional[int] = ...) -> None: ...

class MsgInitResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgIncreaseCounter(_message.Message):
    __slots__ = ("amount",)
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    amount: int
    def __init__(self, amount: _Optional[int] = ...) -> None: ...

class MsgIncreaseCounterResponse(_message.Message):
    __slots__ = ("new_amount",)
    NEW_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    new_amount: int
    def __init__(self, new_amount: _Optional[int] = ...) -> None: ...

class MsgTestDependencies(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgTestDependenciesResponse(_message.Message):
    __slots__ = ("chain_id", "address", "before_gas", "after_gas", "funds")
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    BEFORE_GAS_FIELD_NUMBER: _ClassVar[int]
    AFTER_GAS_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    chain_id: str
    address: str
    before_gas: int
    after_gas: int
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, chain_id: _Optional[str] = ..., address: _Optional[str] = ..., before_gas: _Optional[int] = ..., after_gas: _Optional[int] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class QueryCounterRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryCounterResponse(_message.Message):
    __slots__ = ("value",)
    VALUE_FIELD_NUMBER: _ClassVar[int]
    value: int
    def __init__(self, value: _Optional[int] = ...) -> None: ...
