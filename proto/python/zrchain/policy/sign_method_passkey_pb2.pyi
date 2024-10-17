from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class SignMethodPasskey(_message.Message):
    __slots__ = ("raw_id", "attestation_object", "client_data_json", "active")
    RAW_ID_FIELD_NUMBER: _ClassVar[int]
    ATTESTATION_OBJECT_FIELD_NUMBER: _ClassVar[int]
    CLIENT_DATA_JSON_FIELD_NUMBER: _ClassVar[int]
    ACTIVE_FIELD_NUMBER: _ClassVar[int]
    raw_id: bytes
    attestation_object: bytes
    client_data_json: bytes
    active: bool
    def __init__(self, raw_id: _Optional[bytes] = ..., attestation_object: _Optional[bytes] = ..., client_data_json: _Optional[bytes] = ..., active: bool = ...) -> None: ...
