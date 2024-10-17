from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmwasm.wasm.v1 import types_pb2 as _types_pb2
from google.protobuf import any_pb2 as _any_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class StoreCodeAuthorization(_message.Message):
    __slots__ = ("grants",)
    GRANTS_FIELD_NUMBER: _ClassVar[int]
    grants: _containers.RepeatedCompositeFieldContainer[CodeGrant]
    def __init__(self, grants: _Optional[_Iterable[_Union[CodeGrant, _Mapping]]] = ...) -> None: ...

class ContractExecutionAuthorization(_message.Message):
    __slots__ = ("grants",)
    GRANTS_FIELD_NUMBER: _ClassVar[int]
    grants: _containers.RepeatedCompositeFieldContainer[ContractGrant]
    def __init__(self, grants: _Optional[_Iterable[_Union[ContractGrant, _Mapping]]] = ...) -> None: ...

class ContractMigrationAuthorization(_message.Message):
    __slots__ = ("grants",)
    GRANTS_FIELD_NUMBER: _ClassVar[int]
    grants: _containers.RepeatedCompositeFieldContainer[ContractGrant]
    def __init__(self, grants: _Optional[_Iterable[_Union[ContractGrant, _Mapping]]] = ...) -> None: ...

class CodeGrant(_message.Message):
    __slots__ = ("code_hash", "instantiate_permission")
    CODE_HASH_FIELD_NUMBER: _ClassVar[int]
    INSTANTIATE_PERMISSION_FIELD_NUMBER: _ClassVar[int]
    code_hash: bytes
    instantiate_permission: _types_pb2.AccessConfig
    def __init__(self, code_hash: _Optional[bytes] = ..., instantiate_permission: _Optional[_Union[_types_pb2.AccessConfig, _Mapping]] = ...) -> None: ...

class ContractGrant(_message.Message):
    __slots__ = ("contract", "limit", "filter")
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    FILTER_FIELD_NUMBER: _ClassVar[int]
    contract: str
    limit: _any_pb2.Any
    filter: _any_pb2.Any
    def __init__(self, contract: _Optional[str] = ..., limit: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., filter: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class MaxCallsLimit(_message.Message):
    __slots__ = ("remaining",)
    REMAINING_FIELD_NUMBER: _ClassVar[int]
    remaining: int
    def __init__(self, remaining: _Optional[int] = ...) -> None: ...

class MaxFundsLimit(_message.Message):
    __slots__ = ("amounts",)
    AMOUNTS_FIELD_NUMBER: _ClassVar[int]
    amounts: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, amounts: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class CombinedLimit(_message.Message):
    __slots__ = ("calls_remaining", "amounts")
    CALLS_REMAINING_FIELD_NUMBER: _ClassVar[int]
    AMOUNTS_FIELD_NUMBER: _ClassVar[int]
    calls_remaining: int
    amounts: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, calls_remaining: _Optional[int] = ..., amounts: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class AllowAllMessagesFilter(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class AcceptedMessageKeysFilter(_message.Message):
    __slots__ = ("keys",)
    KEYS_FIELD_NUMBER: _ClassVar[int]
    keys: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, keys: _Optional[_Iterable[str]] = ...) -> None: ...

class AcceptedMessagesFilter(_message.Message):
    __slots__ = ("messages",)
    MESSAGES_FIELD_NUMBER: _ClassVar[int]
    messages: _containers.RepeatedScalarFieldContainer[bytes]
    def __init__(self, messages: _Optional[_Iterable[bytes]] = ...) -> None: ...
