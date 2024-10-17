from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmwasm.wasm.v1 import types_pb2 as _types_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class StoreCodeProposal(_message.Message):
    __slots__ = ("title", "description", "run_as", "wasm_byte_code", "instantiate_permission", "unpin_code", "source", "builder", "code_hash")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    RUN_AS_FIELD_NUMBER: _ClassVar[int]
    WASM_BYTE_CODE_FIELD_NUMBER: _ClassVar[int]
    INSTANTIATE_PERMISSION_FIELD_NUMBER: _ClassVar[int]
    UNPIN_CODE_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    BUILDER_FIELD_NUMBER: _ClassVar[int]
    CODE_HASH_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    run_as: str
    wasm_byte_code: bytes
    instantiate_permission: _types_pb2.AccessConfig
    unpin_code: bool
    source: str
    builder: str
    code_hash: bytes
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., run_as: _Optional[str] = ..., wasm_byte_code: _Optional[bytes] = ..., instantiate_permission: _Optional[_Union[_types_pb2.AccessConfig, _Mapping]] = ..., unpin_code: bool = ..., source: _Optional[str] = ..., builder: _Optional[str] = ..., code_hash: _Optional[bytes] = ...) -> None: ...

class InstantiateContractProposal(_message.Message):
    __slots__ = ("title", "description", "run_as", "admin", "code_id", "label", "msg", "funds")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    RUN_AS_FIELD_NUMBER: _ClassVar[int]
    ADMIN_FIELD_NUMBER: _ClassVar[int]
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    run_as: str
    admin: str
    code_id: int
    label: str
    msg: bytes
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., run_as: _Optional[str] = ..., admin: _Optional[str] = ..., code_id: _Optional[int] = ..., label: _Optional[str] = ..., msg: _Optional[bytes] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class InstantiateContract2Proposal(_message.Message):
    __slots__ = ("title", "description", "run_as", "admin", "code_id", "label", "msg", "funds", "salt", "fix_msg")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    RUN_AS_FIELD_NUMBER: _ClassVar[int]
    ADMIN_FIELD_NUMBER: _ClassVar[int]
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    SALT_FIELD_NUMBER: _ClassVar[int]
    FIX_MSG_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    run_as: str
    admin: str
    code_id: int
    label: str
    msg: bytes
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    salt: bytes
    fix_msg: bool
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., run_as: _Optional[str] = ..., admin: _Optional[str] = ..., code_id: _Optional[int] = ..., label: _Optional[str] = ..., msg: _Optional[bytes] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., salt: _Optional[bytes] = ..., fix_msg: bool = ...) -> None: ...

class MigrateContractProposal(_message.Message):
    __slots__ = ("title", "description", "contract", "code_id", "msg")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    contract: str
    code_id: int
    msg: bytes
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., contract: _Optional[str] = ..., code_id: _Optional[int] = ..., msg: _Optional[bytes] = ...) -> None: ...

class SudoContractProposal(_message.Message):
    __slots__ = ("title", "description", "contract", "msg")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    contract: str
    msg: bytes
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., contract: _Optional[str] = ..., msg: _Optional[bytes] = ...) -> None: ...

class ExecuteContractProposal(_message.Message):
    __slots__ = ("title", "description", "run_as", "contract", "msg", "funds")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    RUN_AS_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    run_as: str
    contract: str
    msg: bytes
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., run_as: _Optional[str] = ..., contract: _Optional[str] = ..., msg: _Optional[bytes] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class UpdateAdminProposal(_message.Message):
    __slots__ = ("title", "description", "new_admin", "contract")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    NEW_ADMIN_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    new_admin: str
    contract: str
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., new_admin: _Optional[str] = ..., contract: _Optional[str] = ...) -> None: ...

class ClearAdminProposal(_message.Message):
    __slots__ = ("title", "description", "contract")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    contract: str
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., contract: _Optional[str] = ...) -> None: ...

class PinCodesProposal(_message.Message):
    __slots__ = ("title", "description", "code_ids")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    CODE_IDS_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    code_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., code_ids: _Optional[_Iterable[int]] = ...) -> None: ...

class UnpinCodesProposal(_message.Message):
    __slots__ = ("title", "description", "code_ids")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    CODE_IDS_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    code_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., code_ids: _Optional[_Iterable[int]] = ...) -> None: ...

class AccessConfigUpdate(_message.Message):
    __slots__ = ("code_id", "instantiate_permission")
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    INSTANTIATE_PERMISSION_FIELD_NUMBER: _ClassVar[int]
    code_id: int
    instantiate_permission: _types_pb2.AccessConfig
    def __init__(self, code_id: _Optional[int] = ..., instantiate_permission: _Optional[_Union[_types_pb2.AccessConfig, _Mapping]] = ...) -> None: ...

class UpdateInstantiateConfigProposal(_message.Message):
    __slots__ = ("title", "description", "access_config_updates")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    ACCESS_CONFIG_UPDATES_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    access_config_updates: _containers.RepeatedCompositeFieldContainer[AccessConfigUpdate]
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., access_config_updates: _Optional[_Iterable[_Union[AccessConfigUpdate, _Mapping]]] = ...) -> None: ...

class StoreAndInstantiateContractProposal(_message.Message):
    __slots__ = ("title", "description", "run_as", "wasm_byte_code", "instantiate_permission", "unpin_code", "admin", "label", "msg", "funds", "source", "builder", "code_hash")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    RUN_AS_FIELD_NUMBER: _ClassVar[int]
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
    title: str
    description: str
    run_as: str
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
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., run_as: _Optional[str] = ..., wasm_byte_code: _Optional[bytes] = ..., instantiate_permission: _Optional[_Union[_types_pb2.AccessConfig, _Mapping]] = ..., unpin_code: bool = ..., admin: _Optional[str] = ..., label: _Optional[str] = ..., msg: _Optional[bytes] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., source: _Optional[str] = ..., builder: _Optional[str] = ..., code_hash: _Optional[bytes] = ...) -> None: ...
