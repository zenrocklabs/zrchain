from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PolicyPacketData(_message.Message):
    __slots__ = ("no_data",)
    NO_DATA_FIELD_NUMBER: _ClassVar[int]
    no_data: NoData
    def __init__(self, no_data: _Optional[_Union[NoData, _Mapping]] = ...) -> None: ...

class NoData(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
