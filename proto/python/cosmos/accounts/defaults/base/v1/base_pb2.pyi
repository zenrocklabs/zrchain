from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MsgInit(_message.Message):
    __slots__ = ("pub_key",)
    PUB_KEY_FIELD_NUMBER: _ClassVar[int]
    pub_key: bytes
    def __init__(self, pub_key: _Optional[bytes] = ...) -> None: ...

class MsgInitResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgSwapPubKey(_message.Message):
    __slots__ = ("new_pub_key",)
    NEW_PUB_KEY_FIELD_NUMBER: _ClassVar[int]
    new_pub_key: bytes
    def __init__(self, new_pub_key: _Optional[bytes] = ...) -> None: ...

class MsgSwapPubKeyResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QuerySequence(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QuerySequenceResponse(_message.Message):
    __slots__ = ("sequence",)
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    sequence: int
    def __init__(self, sequence: _Optional[int] = ...) -> None: ...
