from cosmos.app.v1alpha1 import module_pb2 as _module_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Module(_message.Message):
    __slots__ = ("hooks_order", "authority", "bech32_prefix_validator", "bech32_prefix_consensus")
    HOOKS_ORDER_FIELD_NUMBER: _ClassVar[int]
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    BECH32_PREFIX_VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    BECH32_PREFIX_CONSENSUS_FIELD_NUMBER: _ClassVar[int]
    hooks_order: _containers.RepeatedScalarFieldContainer[str]
    authority: str
    bech32_prefix_validator: str
    bech32_prefix_consensus: str
    def __init__(self, hooks_order: _Optional[_Iterable[str]] = ..., authority: _Optional[str] = ..., bech32_prefix_validator: _Optional[str] = ..., bech32_prefix_consensus: _Optional[str] = ...) -> None: ...
