from google.protobuf import descriptor_pb2 as _descriptor_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor
JAVA_FIELD_NUMBER: _ClassVar[int]
java: _descriptor.FieldDescriptor

class JavaFeatures(_message.Message):
    __slots__ = ("legacy_closed_enum", "utf8_validation", "large_enum", "use_old_outer_classname_default", "nest_in_file_class")
    class Utf8Validation(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        UTF8_VALIDATION_UNKNOWN: _ClassVar[JavaFeatures.Utf8Validation]
        DEFAULT: _ClassVar[JavaFeatures.Utf8Validation]
        VERIFY: _ClassVar[JavaFeatures.Utf8Validation]
    UTF8_VALIDATION_UNKNOWN: JavaFeatures.Utf8Validation
    DEFAULT: JavaFeatures.Utf8Validation
    VERIFY: JavaFeatures.Utf8Validation
    class NestInFileClassFeature(_message.Message):
        __slots__ = ()
        class NestInFileClass(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
            __slots__ = ()
            NEST_IN_FILE_CLASS_UNKNOWN: _ClassVar[JavaFeatures.NestInFileClassFeature.NestInFileClass]
            NO: _ClassVar[JavaFeatures.NestInFileClassFeature.NestInFileClass]
            YES: _ClassVar[JavaFeatures.NestInFileClassFeature.NestInFileClass]
            LEGACY: _ClassVar[JavaFeatures.NestInFileClassFeature.NestInFileClass]
        NEST_IN_FILE_CLASS_UNKNOWN: JavaFeatures.NestInFileClassFeature.NestInFileClass
        NO: JavaFeatures.NestInFileClassFeature.NestInFileClass
        YES: JavaFeatures.NestInFileClassFeature.NestInFileClass
        LEGACY: JavaFeatures.NestInFileClassFeature.NestInFileClass
        def __init__(self) -> None: ...
    LEGACY_CLOSED_ENUM_FIELD_NUMBER: _ClassVar[int]
    UTF8_VALIDATION_FIELD_NUMBER: _ClassVar[int]
    LARGE_ENUM_FIELD_NUMBER: _ClassVar[int]
    USE_OLD_OUTER_CLASSNAME_DEFAULT_FIELD_NUMBER: _ClassVar[int]
    NEST_IN_FILE_CLASS_FIELD_NUMBER: _ClassVar[int]
    legacy_closed_enum: bool
    utf8_validation: JavaFeatures.Utf8Validation
    large_enum: bool
    use_old_outer_classname_default: bool
    nest_in_file_class: JavaFeatures.NestInFileClassFeature.NestInFileClass
    def __init__(self, legacy_closed_enum: bool = ..., utf8_validation: _Optional[_Union[JavaFeatures.Utf8Validation, str]] = ..., large_enum: bool = ..., use_old_outer_classname_default: bool = ..., nest_in_file_class: _Optional[_Union[JavaFeatures.NestInFileClassFeature.NestInFileClass, str]] = ...) -> None: ...
