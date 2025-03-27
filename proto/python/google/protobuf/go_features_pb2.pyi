from google.protobuf import descriptor_pb2 as _descriptor_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor
GO_FIELD_NUMBER: _ClassVar[int]
go: _descriptor.FieldDescriptor

class GoFeatures(_message.Message):
    __slots__ = ("legacy_unmarshal_json_enum", "api_level", "strip_enum_prefix")
    class APILevel(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        API_LEVEL_UNSPECIFIED: _ClassVar[GoFeatures.APILevel]
        API_OPEN: _ClassVar[GoFeatures.APILevel]
        API_HYBRID: _ClassVar[GoFeatures.APILevel]
        API_OPAQUE: _ClassVar[GoFeatures.APILevel]
    API_LEVEL_UNSPECIFIED: GoFeatures.APILevel
    API_OPEN: GoFeatures.APILevel
    API_HYBRID: GoFeatures.APILevel
    API_OPAQUE: GoFeatures.APILevel
    class StripEnumPrefix(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        STRIP_ENUM_PREFIX_UNSPECIFIED: _ClassVar[GoFeatures.StripEnumPrefix]
        STRIP_ENUM_PREFIX_KEEP: _ClassVar[GoFeatures.StripEnumPrefix]
        STRIP_ENUM_PREFIX_GENERATE_BOTH: _ClassVar[GoFeatures.StripEnumPrefix]
        STRIP_ENUM_PREFIX_STRIP: _ClassVar[GoFeatures.StripEnumPrefix]
    STRIP_ENUM_PREFIX_UNSPECIFIED: GoFeatures.StripEnumPrefix
    STRIP_ENUM_PREFIX_KEEP: GoFeatures.StripEnumPrefix
    STRIP_ENUM_PREFIX_GENERATE_BOTH: GoFeatures.StripEnumPrefix
    STRIP_ENUM_PREFIX_STRIP: GoFeatures.StripEnumPrefix
    LEGACY_UNMARSHAL_JSON_ENUM_FIELD_NUMBER: _ClassVar[int]
    API_LEVEL_FIELD_NUMBER: _ClassVar[int]
    STRIP_ENUM_PREFIX_FIELD_NUMBER: _ClassVar[int]
    legacy_unmarshal_json_enum: bool
    api_level: GoFeatures.APILevel
    strip_enum_prefix: GoFeatures.StripEnumPrefix
    def __init__(self, legacy_unmarshal_json_enum: bool = ..., api_level: _Optional[_Union[GoFeatures.APILevel, str]] = ..., strip_enum_prefix: _Optional[_Union[GoFeatures.StripEnumPrefix, str]] = ...) -> None: ...
