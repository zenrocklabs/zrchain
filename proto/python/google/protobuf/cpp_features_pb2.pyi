from google.protobuf import descriptor_pb2 as _descriptor_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor
CPP_FIELD_NUMBER: _ClassVar[int]
cpp: _descriptor.FieldDescriptor

class CppFeatures(_message.Message):
    __slots__ = ("legacy_closed_enum", "string_type", "enum_name_uses_string_view")
    class StringType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        STRING_TYPE_UNKNOWN: _ClassVar[CppFeatures.StringType]
        VIEW: _ClassVar[CppFeatures.StringType]
        CORD: _ClassVar[CppFeatures.StringType]
        STRING: _ClassVar[CppFeatures.StringType]
    STRING_TYPE_UNKNOWN: CppFeatures.StringType
    VIEW: CppFeatures.StringType
    CORD: CppFeatures.StringType
    STRING: CppFeatures.StringType
    LEGACY_CLOSED_ENUM_FIELD_NUMBER: _ClassVar[int]
    STRING_TYPE_FIELD_NUMBER: _ClassVar[int]
    ENUM_NAME_USES_STRING_VIEW_FIELD_NUMBER: _ClassVar[int]
    legacy_closed_enum: bool
    string_type: CppFeatures.StringType
    enum_name_uses_string_view: bool
    def __init__(self, legacy_closed_enum: bool = ..., string_type: _Optional[_Union[CppFeatures.StringType, str]] = ..., enum_name_uses_string_view: bool = ...) -> None: ...
