from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import any_pb2 as _any_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class AccessType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    ACCESS_TYPE_UNSPECIFIED: _ClassVar[AccessType]
    ACCESS_TYPE_NOBODY: _ClassVar[AccessType]
    ACCESS_TYPE_EVERYBODY: _ClassVar[AccessType]
    ACCESS_TYPE_ANY_OF_ADDRESSES: _ClassVar[AccessType]

class ContractCodeHistoryOperationType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    CONTRACT_CODE_HISTORY_OPERATION_TYPE_UNSPECIFIED: _ClassVar[ContractCodeHistoryOperationType]
    CONTRACT_CODE_HISTORY_OPERATION_TYPE_INIT: _ClassVar[ContractCodeHistoryOperationType]
    CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE: _ClassVar[ContractCodeHistoryOperationType]
    CONTRACT_CODE_HISTORY_OPERATION_TYPE_GENESIS: _ClassVar[ContractCodeHistoryOperationType]
ACCESS_TYPE_UNSPECIFIED: AccessType
ACCESS_TYPE_NOBODY: AccessType
ACCESS_TYPE_EVERYBODY: AccessType
ACCESS_TYPE_ANY_OF_ADDRESSES: AccessType
CONTRACT_CODE_HISTORY_OPERATION_TYPE_UNSPECIFIED: ContractCodeHistoryOperationType
CONTRACT_CODE_HISTORY_OPERATION_TYPE_INIT: ContractCodeHistoryOperationType
CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE: ContractCodeHistoryOperationType
CONTRACT_CODE_HISTORY_OPERATION_TYPE_GENESIS: ContractCodeHistoryOperationType

class AccessTypeParam(_message.Message):
    __slots__ = ("value",)
    VALUE_FIELD_NUMBER: _ClassVar[int]
    value: AccessType
    def __init__(self, value: _Optional[_Union[AccessType, str]] = ...) -> None: ...

class AccessConfig(_message.Message):
    __slots__ = ("permission", "addresses")
    PERMISSION_FIELD_NUMBER: _ClassVar[int]
    ADDRESSES_FIELD_NUMBER: _ClassVar[int]
    permission: AccessType
    addresses: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, permission: _Optional[_Union[AccessType, str]] = ..., addresses: _Optional[_Iterable[str]] = ...) -> None: ...

class Params(_message.Message):
    __slots__ = ("code_upload_access", "instantiate_default_permission")
    CODE_UPLOAD_ACCESS_FIELD_NUMBER: _ClassVar[int]
    INSTANTIATE_DEFAULT_PERMISSION_FIELD_NUMBER: _ClassVar[int]
    code_upload_access: AccessConfig
    instantiate_default_permission: AccessType
    def __init__(self, code_upload_access: _Optional[_Union[AccessConfig, _Mapping]] = ..., instantiate_default_permission: _Optional[_Union[AccessType, str]] = ...) -> None: ...

class CodeInfo(_message.Message):
    __slots__ = ("code_hash", "creator", "instantiate_config")
    CODE_HASH_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    INSTANTIATE_CONFIG_FIELD_NUMBER: _ClassVar[int]
    code_hash: bytes
    creator: str
    instantiate_config: AccessConfig
    def __init__(self, code_hash: _Optional[bytes] = ..., creator: _Optional[str] = ..., instantiate_config: _Optional[_Union[AccessConfig, _Mapping]] = ...) -> None: ...

class ContractInfo(_message.Message):
    __slots__ = ("code_id", "creator", "admin", "label", "created", "ibc_port_id", "extension", "ibc2_port_id")
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    ADMIN_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    CREATED_FIELD_NUMBER: _ClassVar[int]
    IBC_PORT_ID_FIELD_NUMBER: _ClassVar[int]
    EXTENSION_FIELD_NUMBER: _ClassVar[int]
    IBC2_PORT_ID_FIELD_NUMBER: _ClassVar[int]
    code_id: int
    creator: str
    admin: str
    label: str
    created: AbsoluteTxPosition
    ibc_port_id: str
    extension: _any_pb2.Any
    ibc2_port_id: str
    def __init__(self, code_id: _Optional[int] = ..., creator: _Optional[str] = ..., admin: _Optional[str] = ..., label: _Optional[str] = ..., created: _Optional[_Union[AbsoluteTxPosition, _Mapping]] = ..., ibc_port_id: _Optional[str] = ..., extension: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., ibc2_port_id: _Optional[str] = ...) -> None: ...

class ContractCodeHistoryEntry(_message.Message):
    __slots__ = ("operation", "code_id", "updated", "msg")
    OPERATION_FIELD_NUMBER: _ClassVar[int]
    CODE_ID_FIELD_NUMBER: _ClassVar[int]
    UPDATED_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    operation: ContractCodeHistoryOperationType
    code_id: int
    updated: AbsoluteTxPosition
    msg: bytes
    def __init__(self, operation: _Optional[_Union[ContractCodeHistoryOperationType, str]] = ..., code_id: _Optional[int] = ..., updated: _Optional[_Union[AbsoluteTxPosition, _Mapping]] = ..., msg: _Optional[bytes] = ...) -> None: ...

class AbsoluteTxPosition(_message.Message):
    __slots__ = ("block_height", "tx_index")
    BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    TX_INDEX_FIELD_NUMBER: _ClassVar[int]
    block_height: int
    tx_index: int
    def __init__(self, block_height: _Optional[int] = ..., tx_index: _Optional[int] = ...) -> None: ...

class Model(_message.Message):
    __slots__ = ("key", "value")
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    key: bytes
    value: bytes
    def __init__(self, key: _Optional[bytes] = ..., value: _Optional[bytes] = ...) -> None: ...
