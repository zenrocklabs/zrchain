from google.protobuf import descriptor_pb2 as _descriptor_pb2
from cosmos.query.v1 import query_pb2 as _query_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class FileDescriptorsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class FileDescriptorsResponse(_message.Message):
    __slots__ = ("files",)
    FILES_FIELD_NUMBER: _ClassVar[int]
    files: _containers.RepeatedCompositeFieldContainer[_descriptor_pb2.FileDescriptorProto]
    def __init__(self, files: _Optional[_Iterable[_Union[_descriptor_pb2.FileDescriptorProto, _Mapping]]] = ...) -> None: ...
