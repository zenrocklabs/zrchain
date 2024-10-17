from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from cosmwasm.wasm.v1 import types_pb2 as _types_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgStoreCode(_message.Message):
    __slots__ = ("sender", "wasm_byte_code", "instantiate_permission")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    WASM_BYTE_CODE_FIELD_NUMBER: _ClassVar[int]
    INSTANTIATE_PERMISSION_FIELD_NUMBER: _ClassVar[int]
    sender: str
    wasm_byte_code: bytes
    instantiate_permission: _types_pb2.AccessConfig
    def __init__(self, sender: _Optional[str] = ..., wasm_byte_code: _Optional[bytes] = ..., instantiate_permission: _Optional[_Union[_types_pb2.AccessConfig, _Mapping]] = ...) -> None: ...

class MsgStoreCodeResponse(_message.Message):
    __slots__ = ("code_id", "checksum")
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    CHECKSUM_FIELD_NUMBER: _ClassVar[int]
    code_id: int
    checksum: bytes
    def __init__(self, code_id: _Optional[int] = ..., checksum: _Optional[bytes] = ...) -> None: ...

class MsgInstantiateContract(_message.Message):
    __slots__ = ("sender", "admin", "code_id", "label", "msg", "funds")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    ADMIN_FIELD_NUMBER: _ClassVar[int]
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    sender: str
    admin: str
    code_id: int
    label: str
    msg: bytes
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, sender: _Optional[str] = ..., admin: _Optional[str] = ..., code_id: _Optional[int] = ..., label: _Optional[str] = ..., msg: _Optional[bytes] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class MsgInstantiateContractResponse(_message.Message):
    __slots__ = ("address", "data")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    address: str
    data: bytes
    def __init__(self, address: _Optional[str] = ..., data: _Optional[bytes] = ...) -> None: ...

class MsgInstantiateContract2(_message.Message):
    __slots__ = ("sender", "admin", "code_id", "label", "msg", "funds", "salt", "fix_msg")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    ADMIN_FIELD_NUMBER: _ClassVar[int]
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    SALT_FIELD_NUMBER: _ClassVar[int]
    FIX_MSG_FIELD_NUMBER: _ClassVar[int]
    sender: str
    admin: str
    code_id: int
    label: str
    msg: bytes
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    salt: bytes
    fix_msg: bool
    def __init__(self, sender: _Optional[str] = ..., admin: _Optional[str] = ..., code_id: _Optional[int] = ..., label: _Optional[str] = ..., msg: _Optional[bytes] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., salt: _Optional[bytes] = ..., fix_msg: bool = ...) -> None: ...

class MsgInstantiateContract2Response(_message.Message):
    __slots__ = ("address", "data")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    address: str
    data: bytes
    def __init__(self, address: _Optional[str] = ..., data: _Optional[bytes] = ...) -> None: ...

class MsgExecuteContract(_message.Message):
    __slots__ = ("sender", "contract", "msg", "funds")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    sender: str
    contract: str
    msg: bytes
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, sender: _Optional[str] = ..., contract: _Optional[str] = ..., msg: _Optional[bytes] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class MsgExecuteContractResponse(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...

class MsgMigrateContract(_message.Message):
    __slots__ = ("sender", "contract", "code_id", "msg")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    sender: str
    contract: str
    code_id: int
    msg: bytes
    def __init__(self, sender: _Optional[str] = ..., contract: _Optional[str] = ..., code_id: _Optional[int] = ..., msg: _Optional[bytes] = ...) -> None: ...

class MsgMigrateContractResponse(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...

class MsgUpdateAdmin(_message.Message):
    __slots__ = ("sender", "new_admin", "contract")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    NEW_ADMIN_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    sender: str
    new_admin: str
    contract: str
    def __init__(self, sender: _Optional[str] = ..., new_admin: _Optional[str] = ..., contract: _Optional[str] = ...) -> None: ...

class MsgUpdateAdminResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgClearAdmin(_message.Message):
    __slots__ = ("sender", "contract")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    sender: str
    contract: str
    def __init__(self, sender: _Optional[str] = ..., contract: _Optional[str] = ...) -> None: ...

class MsgClearAdminResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateInstantiateConfig(_message.Message):
    __slots__ = ("sender", "code_id", "new_instantiate_permission")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    NEW_INSTANTIATE_PERMISSION_FIELD_NUMBER: _ClassVar[int]
    sender: str
    code_id: int
    new_instantiate_permission: _types_pb2.AccessConfig
    def __init__(self, sender: _Optional[str] = ..., code_id: _Optional[int] = ..., new_instantiate_permission: _Optional[_Union[_types_pb2.AccessConfig, _Mapping]] = ...) -> None: ...

class MsgUpdateInstantiateConfigResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    params: _types_pb2.Params
    def __init__(self, authority: _Optional[str] = ..., params: _Optional[_Union[_types_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgSudoContract(_message.Message):
    __slots__ = ("authority", "contract", "msg")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    authority: str
    contract: str
    msg: bytes
    def __init__(self, authority: _Optional[str] = ..., contract: _Optional[str] = ..., msg: _Optional[bytes] = ...) -> None: ...

class MsgSudoContractResponse(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...

class MsgPinCodes(_message.Message):
    __slots__ = ("authority", "code_ids")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    CODE_IDS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    code_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, authority: _Optional[str] = ..., code_ids: _Optional[_Iterable[int]] = ...) -> None: ...

class MsgPinCodesResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUnpinCodes(_message.Message):
    __slots__ = ("authority", "code_ids")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    CODE_IDS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    code_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, authority: _Optional[str] = ..., code_ids: _Optional[_Iterable[int]] = ...) -> None: ...

class MsgUnpinCodesResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgStoreAndInstantiateContract(_message.Message):
    __slots__ = ("authority", "wasm_byte_code", "instantiate_permission", "unpin_code", "admin", "label", "msg", "funds", "source", "builder", "code_hash")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    WASM_BYTE_CODE_FIELD_NUMBER: _ClassVar[int]
    INSTANTIATE_PERMISSION_FIELD_NUMBER: _ClassVar[int]
    UNPIN_CODE_FIELD_NUMBER: _ClassVar[int]
    ADMIN_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    BUILDER_FIELD_NUMBER: _ClassVar[int]
    CODE_HASH_FIELD_NUMBER: _ClassVar[int]
    authority: str
    wasm_byte_code: bytes
    instantiate_permission: _types_pb2.AccessConfig
    unpin_code: bool
    admin: str
    label: str
    msg: bytes
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    source: str
    builder: str
    code_hash: bytes
    def __init__(self, authority: _Optional[str] = ..., wasm_byte_code: _Optional[bytes] = ..., instantiate_permission: _Optional[_Union[_types_pb2.AccessConfig, _Mapping]] = ..., unpin_code: bool = ..., admin: _Optional[str] = ..., label: _Optional[str] = ..., msg: _Optional[bytes] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., source: _Optional[str] = ..., builder: _Optional[str] = ..., code_hash: _Optional[bytes] = ...) -> None: ...

class MsgStoreAndInstantiateContractResponse(_message.Message):
    __slots__ = ("address", "data")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    address: str
    data: bytes
    def __init__(self, address: _Optional[str] = ..., data: _Optional[bytes] = ...) -> None: ...

class MsgAddCodeUploadParamsAddresses(_message.Message):
    __slots__ = ("authority", "addresses")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    ADDRESSES_FIELD_NUMBER: _ClassVar[int]
    authority: str
    addresses: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, authority: _Optional[str] = ..., addresses: _Optional[_Iterable[str]] = ...) -> None: ...

class MsgAddCodeUploadParamsAddressesResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgRemoveCodeUploadParamsAddresses(_message.Message):
    __slots__ = ("authority", "addresses")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    ADDRESSES_FIELD_NUMBER: _ClassVar[int]
    authority: str
    addresses: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, authority: _Optional[str] = ..., addresses: _Optional[_Iterable[str]] = ...) -> None: ...

class MsgRemoveCodeUploadParamsAddressesResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgStoreAndMigrateContract(_message.Message):
    __slots__ = ("authority", "wasm_byte_code", "instantiate_permission", "contract", "msg")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    WASM_BYTE_CODE_FIELD_NUMBER: _ClassVar[int]
    INSTANTIATE_PERMISSION_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    authority: str
    wasm_byte_code: bytes
    instantiate_permission: _types_pb2.AccessConfig
    contract: str
    msg: bytes
    def __init__(self, authority: _Optional[str] = ..., wasm_byte_code: _Optional[bytes] = ..., instantiate_permission: _Optional[_Union[_types_pb2.AccessConfig, _Mapping]] = ..., contract: _Optional[str] = ..., msg: _Optional[bytes] = ...) -> None: ...

class MsgStoreAndMigrateContractResponse(_message.Message):
    __slots__ = ("code_id", "checksum", "data")
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    CHECKSUM_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    code_id: int
    checksum: bytes
    data: bytes
    def __init__(self, code_id: _Optional[int] = ..., checksum: _Optional[bytes] = ..., data: _Optional[bytes] = ...) -> None: ...

class MsgUpdateContractLabel(_message.Message):
    __slots__ = ("sender", "new_label", "contract")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    NEW_LABEL_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    sender: str
    new_label: str
    contract: str
    def __init__(self, sender: _Optional[str] = ..., new_label: _Optional[str] = ..., contract: _Optional[str] = ...) -> None: ...

class MsgUpdateContractLabelResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
