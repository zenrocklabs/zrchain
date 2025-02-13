from amino import amino_pb2 as _amino_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.zentp import params_pb2 as _params_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    params: _params_pb2.Params
    def __init__(self, authority: _Optional[str] = ..., params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgMintRock(_message.Message):
    __slots__ = ("creator", "amount", "source_key_id", "destination_chain", "recipient_key_id")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    SOURCE_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_CHAIN_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    creator: str
    amount: int
    source_key_id: int
    destination_chain: str
    recipient_key_id: int
    def __init__(self, creator: _Optional[str] = ..., amount: _Optional[int] = ..., source_key_id: _Optional[int] = ..., destination_chain: _Optional[str] = ..., recipient_key_id: _Optional[int] = ...) -> None: ...

class MsgMintRockResponse(_message.Message):
    __slots__ = ("tx_hash",)
    TX_HASH_FIELD_NUMBER: _ClassVar[int]
    tx_hash: str
    def __init__(self, tx_hash: _Optional[str] = ...) -> None: ...

class MsgBurn(_message.Message):
    __slots__ = ("authority", "module_account", "denom", "amount")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    MODULE_ACCOUNT_FIELD_NUMBER: _ClassVar[int]
    DENOM_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    authority: str
    module_account: str
    denom: str
    amount: int
    def __init__(self, authority: _Optional[str] = ..., module_account: _Optional[str] = ..., denom: _Optional[str] = ..., amount: _Optional[int] = ...) -> None: ...

class MsgBurnResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
