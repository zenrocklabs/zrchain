from google.protobuf import any_pb2 as _any_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Policy(_message.Message):
    __slots__ = ("creator", "id", "name", "policy", "btl")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    POLICY_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    creator: str
    id: int
    name: str
    policy: _any_pb2.Any
    btl: int
    def __init__(self, creator: _Optional[str] = ..., id: _Optional[int] = ..., name: _Optional[str] = ..., policy: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., btl: _Optional[int] = ...) -> None: ...

class BoolparserPolicy(_message.Message):
    __slots__ = ("definition", "participants")
    DEFINITION_FIELD_NUMBER: _ClassVar[int]
    PARTICIPANTS_FIELD_NUMBER: _ClassVar[int]
    definition: str
    participants: _containers.RepeatedCompositeFieldContainer[PolicyParticipant]
    def __init__(self, definition: _Optional[str] = ..., participants: _Optional[_Iterable[_Union[PolicyParticipant, _Mapping]]] = ...) -> None: ...

class PolicyParticipant(_message.Message):
    __slots__ = ("abbreviation", "address")
    ABBREVIATION_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    abbreviation: str
    address: str
    def __init__(self, abbreviation: _Optional[str] = ..., address: _Optional[str] = ...) -> None: ...
