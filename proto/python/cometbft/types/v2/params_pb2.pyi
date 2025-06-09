from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf import wrappers_pb2 as _wrappers_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ConsensusParams(_message.Message):
    __slots__ = ("block", "evidence", "validator", "version", "abci", "synchrony", "feature")
    BLOCK_FIELD_NUMBER: _ClassVar[int]
    EVIDENCE_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    ABCI_FIELD_NUMBER: _ClassVar[int]
    SYNCHRONY_FIELD_NUMBER: _ClassVar[int]
    FEATURE_FIELD_NUMBER: _ClassVar[int]
    block: BlockParams
    evidence: EvidenceParams
    validator: ValidatorParams
    version: VersionParams
    abci: ABCIParams
    synchrony: SynchronyParams
    feature: FeatureParams
    def __init__(self, block: _Optional[_Union[BlockParams, _Mapping]] = ..., evidence: _Optional[_Union[EvidenceParams, _Mapping]] = ..., validator: _Optional[_Union[ValidatorParams, _Mapping]] = ..., version: _Optional[_Union[VersionParams, _Mapping]] = ..., abci: _Optional[_Union[ABCIParams, _Mapping]] = ..., synchrony: _Optional[_Union[SynchronyParams, _Mapping]] = ..., feature: _Optional[_Union[FeatureParams, _Mapping]] = ...) -> None: ...

class BlockParams(_message.Message):
    __slots__ = ("max_bytes", "max_gas")
    MAX_BYTES_FIELD_NUMBER: _ClassVar[int]
    MAX_GAS_FIELD_NUMBER: _ClassVar[int]
    max_bytes: int
    max_gas: int
    def __init__(self, max_bytes: _Optional[int] = ..., max_gas: _Optional[int] = ...) -> None: ...

class EvidenceParams(_message.Message):
    __slots__ = ("max_age_num_blocks", "max_age_duration", "max_bytes")
    MAX_AGE_NUM_BLOCKS_FIELD_NUMBER: _ClassVar[int]
    MAX_AGE_DURATION_FIELD_NUMBER: _ClassVar[int]
    MAX_BYTES_FIELD_NUMBER: _ClassVar[int]
    max_age_num_blocks: int
    max_age_duration: _duration_pb2.Duration
    max_bytes: int
    def __init__(self, max_age_num_blocks: _Optional[int] = ..., max_age_duration: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., max_bytes: _Optional[int] = ...) -> None: ...

class ValidatorParams(_message.Message):
    __slots__ = ("pub_key_types",)
    PUB_KEY_TYPES_FIELD_NUMBER: _ClassVar[int]
    pub_key_types: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, pub_key_types: _Optional[_Iterable[str]] = ...) -> None: ...

class VersionParams(_message.Message):
    __slots__ = ("app",)
    APP_FIELD_NUMBER: _ClassVar[int]
    app: int
    def __init__(self, app: _Optional[int] = ...) -> None: ...

class HashedParams(_message.Message):
    __slots__ = ("block_max_bytes", "block_max_gas")
    BLOCK_MAX_BYTES_FIELD_NUMBER: _ClassVar[int]
    BLOCK_MAX_GAS_FIELD_NUMBER: _ClassVar[int]
    block_max_bytes: int
    block_max_gas: int
    def __init__(self, block_max_bytes: _Optional[int] = ..., block_max_gas: _Optional[int] = ...) -> None: ...

class SynchronyParams(_message.Message):
    __slots__ = ("precision", "message_delay")
    PRECISION_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_DELAY_FIELD_NUMBER: _ClassVar[int]
    precision: _duration_pb2.Duration
    message_delay: _duration_pb2.Duration
    def __init__(self, precision: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., message_delay: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ...) -> None: ...

class FeatureParams(_message.Message):
    __slots__ = ("vote_extensions_enable_height", "pbts_enable_height")
    VOTE_EXTENSIONS_ENABLE_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    PBTS_ENABLE_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    vote_extensions_enable_height: _wrappers_pb2.Int64Value
    pbts_enable_height: _wrappers_pb2.Int64Value
    def __init__(self, vote_extensions_enable_height: _Optional[_Union[_wrappers_pb2.Int64Value, _Mapping]] = ..., pbts_enable_height: _Optional[_Union[_wrappers_pb2.Int64Value, _Mapping]] = ...) -> None: ...

class ABCIParams(_message.Message):
    __slots__ = ("vote_extensions_enable_height",)
    VOTE_EXTENSIONS_ENABLE_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    vote_extensions_enable_height: int
    def __init__(self, vote_extensions_enable_height: _Optional[int] = ...) -> None: ...
