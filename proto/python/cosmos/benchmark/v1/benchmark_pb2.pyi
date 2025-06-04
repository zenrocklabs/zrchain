from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Op(_message.Message):
    __slots__ = ("seed", "actor", "key_length", "value_length", "iterations", "delete", "exists")
    SEED_FIELD_NUMBER: _ClassVar[int]
    ACTOR_FIELD_NUMBER: _ClassVar[int]
    KEY_LENGTH_FIELD_NUMBER: _ClassVar[int]
    VALUE_LENGTH_FIELD_NUMBER: _ClassVar[int]
    ITERATIONS_FIELD_NUMBER: _ClassVar[int]
    DELETE_FIELD_NUMBER: _ClassVar[int]
    EXISTS_FIELD_NUMBER: _ClassVar[int]
    seed: int
    actor: str
    key_length: int
    value_length: int
    iterations: int
    delete: bool
    exists: bool
    def __init__(self, seed: _Optional[int] = ..., actor: _Optional[str] = ..., key_length: _Optional[int] = ..., value_length: _Optional[int] = ..., iterations: _Optional[int] = ..., delete: bool = ..., exists: bool = ...) -> None: ...
