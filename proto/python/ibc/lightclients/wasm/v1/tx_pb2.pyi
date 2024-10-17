from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MsgStoreCode(_message.Message):
    __slots__ = ("signer", "wasm_byte_code")
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    WASM_BYTE_CODE_FIELD_NUMBER: _ClassVar[int]
    signer: str
    wasm_byte_code: bytes
    def __init__(self, signer: _Optional[str] = ..., wasm_byte_code: _Optional[bytes] = ...) -> None: ...

class MsgStoreCodeResponse(_message.Message):
    __slots__ = ("checksum",)
    CHECKSUM_FIELD_NUMBER: _ClassVar[int]
    checksum: bytes
    def __init__(self, checksum: _Optional[bytes] = ...) -> None: ...

class MsgRemoveChecksum(_message.Message):
    __slots__ = ("signer", "checksum")
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    CHECKSUM_FIELD_NUMBER: _ClassVar[int]
    signer: str
    checksum: bytes
    def __init__(self, signer: _Optional[str] = ..., checksum: _Optional[bytes] = ...) -> None: ...

class MsgRemoveChecksumResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgMigrateContract(_message.Message):
    __slots__ = ("signer", "client_id", "checksum", "msg")
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    CHECKSUM_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    signer: str
    client_id: str
    checksum: bytes
    msg: bytes
    def __init__(self, signer: _Optional[str] = ..., client_id: _Optional[str] = ..., checksum: _Optional[bytes] = ..., msg: _Optional[bytes] = ...) -> None: ...

class MsgMigrateContractResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
