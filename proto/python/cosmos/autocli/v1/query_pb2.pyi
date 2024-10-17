from cosmos.autocli.v1 import options_pb2 as _options_pb2
from cosmos.query.v1 import query_pb2 as _query_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class AppOptionsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class AppOptionsResponse(_message.Message):
    __slots__ = ("module_options",)
    class ModuleOptionsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _options_pb2.ModuleOptions
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_options_pb2.ModuleOptions, _Mapping]] = ...) -> None: ...
    MODULE_OPTIONS_FIELD_NUMBER: _ClassVar[int]
    module_options: _containers.MessageMap[str, _options_pb2.ModuleOptions]
    def __init__(self, module_options: _Optional[_Mapping[str, _options_pb2.ModuleOptions]] = ...) -> None: ...
