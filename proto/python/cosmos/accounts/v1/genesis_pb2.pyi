from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("account_number", "accounts")
    ACCOUNT_NUMBER_FIELD_NUMBER: _ClassVar[int]
    ACCOUNTS_FIELD_NUMBER: _ClassVar[int]
    account_number: int
    accounts: _containers.RepeatedCompositeFieldContainer[GenesisAccount]
    def __init__(self, account_number: _Optional[int] = ..., accounts: _Optional[_Iterable[_Union[GenesisAccount, _Mapping]]] = ...) -> None: ...

class GenesisAccount(_message.Message):
    __slots__ = ("address", "account_type", "account_number", "state")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    ACCOUNT_TYPE_FIELD_NUMBER: _ClassVar[int]
    ACCOUNT_NUMBER_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    address: str
    account_type: str
    account_number: int
    state: _containers.RepeatedCompositeFieldContainer[KVPair]
    def __init__(self, address: _Optional[str] = ..., account_type: _Optional[str] = ..., account_number: _Optional[int] = ..., state: _Optional[_Iterable[_Union[KVPair, _Mapping]]] = ...) -> None: ...

class KVPair(_message.Message):
    __slots__ = ("key", "value")
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    key: bytes
    value: bytes
    def __init__(self, key: _Optional[bytes] = ..., value: _Optional[bytes] = ...) -> None: ...
