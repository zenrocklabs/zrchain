from google.protobuf import any_pb2 as _any_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Config(_message.Message):
    __slots__ = ("modules", "golang_bindings")
    MODULES_FIELD_NUMBER: _ClassVar[int]
    GOLANG_BINDINGS_FIELD_NUMBER: _ClassVar[int]
    modules: _containers.RepeatedCompositeFieldContainer[ModuleConfig]
    golang_bindings: _containers.RepeatedCompositeFieldContainer[GolangBinding]
    def __init__(self, modules: _Optional[_Iterable[_Union[ModuleConfig, _Mapping]]] = ..., golang_bindings: _Optional[_Iterable[_Union[GolangBinding, _Mapping]]] = ...) -> None: ...

class ModuleConfig(_message.Message):
    __slots__ = ("name", "config", "golang_bindings")
    NAME_FIELD_NUMBER: _ClassVar[int]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    GOLANG_BINDINGS_FIELD_NUMBER: _ClassVar[int]
    name: str
    config: _any_pb2.Any
    golang_bindings: _containers.RepeatedCompositeFieldContainer[GolangBinding]
    def __init__(self, name: _Optional[str] = ..., config: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., golang_bindings: _Optional[_Iterable[_Union[GolangBinding, _Mapping]]] = ...) -> None: ...

class GolangBinding(_message.Message):
    __slots__ = ("interface_type", "implementation")
    INTERFACE_TYPE_FIELD_NUMBER: _ClassVar[int]
    IMPLEMENTATION_FIELD_NUMBER: _ClassVar[int]
    interface_type: str
    implementation: str
    def __init__(self, interface_type: _Optional[str] = ..., implementation: _Optional[str] = ...) -> None: ...
