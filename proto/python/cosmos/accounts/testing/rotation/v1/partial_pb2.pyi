from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MsgInit(_message.Message):
    __slots__ = ("pub_key_bytes",)
    PUB_KEY_BYTES_FIELD_NUMBER: _ClassVar[int]
    pub_key_bytes: bytes
    def __init__(self, pub_key_bytes: _Optional[bytes] = ...) -> None: ...

class MsgInitResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgRotatePubKey(_message.Message):
    __slots__ = ("new_pub_key_bytes",)
    NEW_PUB_KEY_BYTES_FIELD_NUMBER: _ClassVar[int]
    new_pub_key_bytes: bytes
    def __init__(self, new_pub_key_bytes: _Optional[bytes] = ...) -> None: ...

class MsgRotatePubKeyResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
