from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ModuleOptions(_message.Message):
    __slots__ = ("tx", "query")
    TX_FIELD_NUMBER: _ClassVar[int]
    QUERY_FIELD_NUMBER: _ClassVar[int]
    tx: ServiceCommandDescriptor
    query: ServiceCommandDescriptor
    def __init__(self, tx: _Optional[_Union[ServiceCommandDescriptor, _Mapping]] = ..., query: _Optional[_Union[ServiceCommandDescriptor, _Mapping]] = ...) -> None: ...

class ServiceCommandDescriptor(_message.Message):
    __slots__ = ("service", "rpc_command_options", "sub_commands", "enhance_custom_command", "short")
    class SubCommandsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: ServiceCommandDescriptor
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[ServiceCommandDescriptor, _Mapping]] = ...) -> None: ...
    SERVICE_FIELD_NUMBER: _ClassVar[int]
    RPC_COMMAND_OPTIONS_FIELD_NUMBER: _ClassVar[int]
    SUB_COMMANDS_FIELD_NUMBER: _ClassVar[int]
    ENHANCE_CUSTOM_COMMAND_FIELD_NUMBER: _ClassVar[int]
    SHORT_FIELD_NUMBER: _ClassVar[int]
    service: str
    rpc_command_options: _containers.RepeatedCompositeFieldContainer[RpcCommandOptions]
    sub_commands: _containers.MessageMap[str, ServiceCommandDescriptor]
    enhance_custom_command: bool
    short: str
    def __init__(self, service: _Optional[str] = ..., rpc_command_options: _Optional[_Iterable[_Union[RpcCommandOptions, _Mapping]]] = ..., sub_commands: _Optional[_Mapping[str, ServiceCommandDescriptor]] = ..., enhance_custom_command: bool = ..., short: _Optional[str] = ...) -> None: ...

class RpcCommandOptions(_message.Message):
    __slots__ = ("rpc_method", "use", "long", "short", "example", "alias", "suggest_for", "deprecated", "version", "flag_options", "positional_args", "skip", "gov_proposal")
    class FlagOptionsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: FlagOptions
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[FlagOptions, _Mapping]] = ...) -> None: ...
    RPC_METHOD_FIELD_NUMBER: _ClassVar[int]
    USE_FIELD_NUMBER: _ClassVar[int]
    LONG_FIELD_NUMBER: _ClassVar[int]
    SHORT_FIELD_NUMBER: _ClassVar[int]
    EXAMPLE_FIELD_NUMBER: _ClassVar[int]
    ALIAS_FIELD_NUMBER: _ClassVar[int]
    SUGGEST_FOR_FIELD_NUMBER: _ClassVar[int]
    DEPRECATED_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    FLAG_OPTIONS_FIELD_NUMBER: _ClassVar[int]
    POSITIONAL_ARGS_FIELD_NUMBER: _ClassVar[int]
    SKIP_FIELD_NUMBER: _ClassVar[int]
    GOV_PROPOSAL_FIELD_NUMBER: _ClassVar[int]
    rpc_method: str
    use: str
    long: str
    short: str
    example: str
    alias: _containers.RepeatedScalarFieldContainer[str]
    suggest_for: _containers.RepeatedScalarFieldContainer[str]
    deprecated: str
    version: str
    flag_options: _containers.MessageMap[str, FlagOptions]
    positional_args: _containers.RepeatedCompositeFieldContainer[PositionalArgDescriptor]
    skip: bool
    gov_proposal: bool
    def __init__(self, rpc_method: _Optional[str] = ..., use: _Optional[str] = ..., long: _Optional[str] = ..., short: _Optional[str] = ..., example: _Optional[str] = ..., alias: _Optional[_Iterable[str]] = ..., suggest_for: _Optional[_Iterable[str]] = ..., deprecated: _Optional[str] = ..., version: _Optional[str] = ..., flag_options: _Optional[_Mapping[str, FlagOptions]] = ..., positional_args: _Optional[_Iterable[_Union[PositionalArgDescriptor, _Mapping]]] = ..., skip: bool = ..., gov_proposal: bool = ...) -> None: ...

class FlagOptions(_message.Message):
    __slots__ = ("name", "shorthand", "usage", "default_value", "deprecated", "shorthand_deprecated", "hidden")
    NAME_FIELD_NUMBER: _ClassVar[int]
    SHORTHAND_FIELD_NUMBER: _ClassVar[int]
    USAGE_FIELD_NUMBER: _ClassVar[int]
    DEFAULT_VALUE_FIELD_NUMBER: _ClassVar[int]
    DEPRECATED_FIELD_NUMBER: _ClassVar[int]
    SHORTHAND_DEPRECATED_FIELD_NUMBER: _ClassVar[int]
    HIDDEN_FIELD_NUMBER: _ClassVar[int]
    name: str
    shorthand: str
    usage: str
    default_value: str
    deprecated: str
    shorthand_deprecated: str
    hidden: bool
    def __init__(self, name: _Optional[str] = ..., shorthand: _Optional[str] = ..., usage: _Optional[str] = ..., default_value: _Optional[str] = ..., deprecated: _Optional[str] = ..., shorthand_deprecated: _Optional[str] = ..., hidden: bool = ...) -> None: ...

class PositionalArgDescriptor(_message.Message):
    __slots__ = ("proto_field", "varargs", "optional")
    PROTO_FIELD_FIELD_NUMBER: _ClassVar[int]
    VARARGS_FIELD_NUMBER: _ClassVar[int]
    OPTIONAL_FIELD_NUMBER: _ClassVar[int]
    proto_field: str
    varargs: bool
    optional: bool
    def __init__(self, proto_field: _Optional[str] = ..., varargs: bool = ..., optional: bool = ...) -> None: ...
