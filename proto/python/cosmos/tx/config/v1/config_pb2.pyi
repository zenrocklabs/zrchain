from cosmos.app.v1alpha1 import module_pb2 as _module_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Config(_message.Message):
    __slots__ = ("skip_ante_handler", "skip_post_handler")
    SKIP_ANTE_HANDLER_FIELD_NUMBER: _ClassVar[int]
    SKIP_POST_HANDLER_FIELD_NUMBER: _ClassVar[int]
    skip_ante_handler: bool
    skip_post_handler: bool
    def __init__(self, skip_ante_handler: bool = ..., skip_post_handler: bool = ...) -> None: ...
