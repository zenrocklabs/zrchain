from zrchain.dct import params_pb2 as _params_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Supply(_message.Message):
    __slots__ = ("asset", "custodied_amount", "minted_amount", "pending_amount")
    ASSET_FIELD_NUMBER: _ClassVar[int]
    CUSTODIED_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    MINTED_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    PENDING_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    asset: _params_pb2.Asset
    custodied_amount: int
    minted_amount: int
    pending_amount: int
    def __init__(self, asset: _Optional[_Union[_params_pb2.Asset, str]] = ..., custodied_amount: _Optional[int] = ..., minted_amount: _Optional[int] = ..., pending_amount: _Optional[int] = ...) -> None: ...
