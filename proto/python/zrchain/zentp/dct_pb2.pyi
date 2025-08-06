from zrchain.zentp import params_pb2 as _params_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class DctStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    DCT_STATUS_UNSPECIFIED: _ClassVar[DctStatus]
    DCT_STATUS_KEYS_REQUESTED: _ClassVar[DctStatus]
    DCT_STATUS_KEYS_CREATED: _ClassVar[DctStatus]
    DCT_STATUS_SPL_REQUESTED: _ClassVar[DctStatus]
    DCT_STATUS_COMPLETED: _ClassVar[DctStatus]
    DCT_STATUS_FAILED: _ClassVar[DctStatus]
DCT_STATUS_UNSPECIFIED: DctStatus
DCT_STATUS_KEYS_REQUESTED: DctStatus
DCT_STATUS_KEYS_CREATED: DctStatus
DCT_STATUS_SPL_REQUESTED: DctStatus
DCT_STATUS_COMPLETED: DctStatus
DCT_STATUS_FAILED: DctStatus

class Dct(_message.Message):
    __slots__ = ("denom", "solana", "status")
    DENOM_FIELD_NUMBER: _ClassVar[int]
    SOLANA_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    denom: str
    solana: _params_pb2.Solana
    status: DctStatus
    def __init__(self, denom: _Optional[str] = ..., solana: _Optional[_Union[_params_pb2.Solana, _Mapping]] = ..., status: _Optional[_Union[DctStatus, str]] = ...) -> None: ...
