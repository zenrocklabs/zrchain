from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class AdditionalSignaturePasskey(_message.Message):
    __slots__ = ("raw_id", "authenticator_data", "client_data_json", "signature")
    RAW_ID_FIELD_NUMBER: _ClassVar[int]
    AUTHENTICATOR_DATA_FIELD_NUMBER: _ClassVar[int]
    CLIENT_DATA_JSON_FIELD_NUMBER: _ClassVar[int]
    SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    raw_id: bytes
    authenticator_data: bytes
    client_data_json: bytes
    signature: bytes
    def __init__(self, raw_id: _Optional[bytes] = ..., authenticator_data: _Optional[bytes] = ..., client_data_json: _Optional[bytes] = ..., signature: _Optional[bytes] = ...) -> None: ...
