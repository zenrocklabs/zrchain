from tendermint.crypto import proof_pb2 as _proof_pb2
from tendermint.crypto import keys_pb2 as _keys_pb2
from tendermint.types import params_pb2 as _params_pb2
from tendermint.types import validator_pb2 as _validator_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class CheckTxType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    NEW: _ClassVar[CheckTxType]
    RECHECK: _ClassVar[CheckTxType]

class MisbehaviorType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNKNOWN: _ClassVar[MisbehaviorType]
    DUPLICATE_VOTE: _ClassVar[MisbehaviorType]
    LIGHT_CLIENT_ATTACK: _ClassVar[MisbehaviorType]
NEW: CheckTxType
RECHECK: CheckTxType
UNKNOWN: MisbehaviorType
DUPLICATE_VOTE: MisbehaviorType
LIGHT_CLIENT_ATTACK: MisbehaviorType

class Request(_message.Message):
    __slots__ = ("echo", "flush", "info", "init_chain", "query", "check_tx", "commit", "list_snapshots", "offer_snapshot", "load_snapshot_chunk", "apply_snapshot_chunk", "prepare_proposal", "process_proposal", "extend_vote", "verify_vote_extension", "finalize_block")
    ECHO_FIELD_NUMBER: _ClassVar[int]
    FLUSH_FIELD_NUMBER: _ClassVar[int]
    INFO_FIELD_NUMBER: _ClassVar[int]
    INIT_CHAIN_FIELD_NUMBER: _ClassVar[int]
    QUERY_FIELD_NUMBER: _ClassVar[int]
    CHECK_TX_FIELD_NUMBER: _ClassVar[int]
    COMMIT_FIELD_NUMBER: _ClassVar[int]
    LIST_SNAPSHOTS_FIELD_NUMBER: _ClassVar[int]
    OFFER_SNAPSHOT_FIELD_NUMBER: _ClassVar[int]
    LOAD_SNAPSHOT_CHUNK_FIELD_NUMBER: _ClassVar[int]
    APPLY_SNAPSHOT_CHUNK_FIELD_NUMBER: _ClassVar[int]
    PREPARE_PROPOSAL_FIELD_NUMBER: _ClassVar[int]
    PROCESS_PROPOSAL_FIELD_NUMBER: _ClassVar[int]
    EXTEND_VOTE_FIELD_NUMBER: _ClassVar[int]
    VERIFY_VOTE_EXTENSION_FIELD_NUMBER: _ClassVar[int]
    FINALIZE_BLOCK_FIELD_NUMBER: _ClassVar[int]
    echo: RequestEcho
    flush: RequestFlush
    info: RequestInfo
    init_chain: RequestInitChain
    query: RequestQuery
    check_tx: RequestCheckTx
    commit: RequestCommit
    list_snapshots: RequestListSnapshots
    offer_snapshot: RequestOfferSnapshot
    load_snapshot_chunk: RequestLoadSnapshotChunk
    apply_snapshot_chunk: RequestApplySnapshotChunk
    prepare_proposal: RequestPrepareProposal
    process_proposal: RequestProcessProposal
    extend_vote: RequestExtendVote
    verify_vote_extension: RequestVerifyVoteExtension
    finalize_block: RequestFinalizeBlock
    def __init__(self, echo: _Optional[_Union[RequestEcho, _Mapping]] = ..., flush: _Optional[_Union[RequestFlush, _Mapping]] = ..., info: _Optional[_Union[RequestInfo, _Mapping]] = ..., init_chain: _Optional[_Union[RequestInitChain, _Mapping]] = ..., query: _Optional[_Union[RequestQuery, _Mapping]] = ..., check_tx: _Optional[_Union[RequestCheckTx, _Mapping]] = ..., commit: _Optional[_Union[RequestCommit, _Mapping]] = ..., list_snapshots: _Optional[_Union[RequestListSnapshots, _Mapping]] = ..., offer_snapshot: _Optional[_Union[RequestOfferSnapshot, _Mapping]] = ..., load_snapshot_chunk: _Optional[_Union[RequestLoadSnapshotChunk, _Mapping]] = ..., apply_snapshot_chunk: _Optional[_Union[RequestApplySnapshotChunk, _Mapping]] = ..., prepare_proposal: _Optional[_Union[RequestPrepareProposal, _Mapping]] = ..., process_proposal: _Optional[_Union[RequestProcessProposal, _Mapping]] = ..., extend_vote: _Optional[_Union[RequestExtendVote, _Mapping]] = ..., verify_vote_extension: _Optional[_Union[RequestVerifyVoteExtension, _Mapping]] = ..., finalize_block: _Optional[_Union[RequestFinalizeBlock, _Mapping]] = ...) -> None: ...

class RequestEcho(_message.Message):
    __slots__ = ("message",)
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    message: str
    def __init__(self, message: _Optional[str] = ...) -> None: ...

class RequestFlush(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class RequestInfo(_message.Message):
    __slots__ = ("version", "block_version", "p2p_version", "abci_version")
    VERSION_FIELD_NUMBER: _ClassVar[int]
    BLOCK_VERSION_FIELD_NUMBER: _ClassVar[int]
    P2P_VERSION_FIELD_NUMBER: _ClassVar[int]
    ABCI_VERSION_FIELD_NUMBER: _ClassVar[int]
    version: str
    block_version: int
    p2p_version: int
    abci_version: str
    def __init__(self, version: _Optional[str] = ..., block_version: _Optional[int] = ..., p2p_version: _Optional[int] = ..., abci_version: _Optional[str] = ...) -> None: ...

class RequestInitChain(_message.Message):
    __slots__ = ("time", "chain_id", "consensus_params", "validators", "app_state_bytes", "initial_height")
    TIME_FIELD_NUMBER: _ClassVar[int]
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    CONSENSUS_PARAMS_FIELD_NUMBER: _ClassVar[int]
    VALIDATORS_FIELD_NUMBER: _ClassVar[int]
    APP_STATE_BYTES_FIELD_NUMBER: _ClassVar[int]
    INITIAL_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    time: _timestamp_pb2.Timestamp
    chain_id: str
    consensus_params: _params_pb2.ConsensusParams
    validators: _containers.RepeatedCompositeFieldContainer[ValidatorUpdate]
    app_state_bytes: bytes
    initial_height: int
    def __init__(self, time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., chain_id: _Optional[str] = ..., consensus_params: _Optional[_Union[_params_pb2.ConsensusParams, _Mapping]] = ..., validators: _Optional[_Iterable[_Union[ValidatorUpdate, _Mapping]]] = ..., app_state_bytes: _Optional[bytes] = ..., initial_height: _Optional[int] = ...) -> None: ...

class RequestQuery(_message.Message):
    __slots__ = ("data", "path", "height", "prove")
    DATA_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    PROVE_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    path: str
    height: int
    prove: bool
    def __init__(self, data: _Optional[bytes] = ..., path: _Optional[str] = ..., height: _Optional[int] = ..., prove: bool = ...) -> None: ...

class RequestCheckTx(_message.Message):
    __slots__ = ("tx", "type")
    TX_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    tx: bytes
    type: CheckTxType
    def __init__(self, tx: _Optional[bytes] = ..., type: _Optional[_Union[CheckTxType, str]] = ...) -> None: ...

class RequestCommit(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class RequestListSnapshots(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class RequestOfferSnapshot(_message.Message):
    __slots__ = ("snapshot", "app_hash")
    SNAPSHOT_FIELD_NUMBER: _ClassVar[int]
    APP_HASH_FIELD_NUMBER: _ClassVar[int]
    snapshot: Snapshot
    app_hash: bytes
    def __init__(self, snapshot: _Optional[_Union[Snapshot, _Mapping]] = ..., app_hash: _Optional[bytes] = ...) -> None: ...

class RequestLoadSnapshotChunk(_message.Message):
    __slots__ = ("height", "format", "chunk")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    FORMAT_FIELD_NUMBER: _ClassVar[int]
    CHUNK_FIELD_NUMBER: _ClassVar[int]
    height: int
    format: int
    chunk: int
    def __init__(self, height: _Optional[int] = ..., format: _Optional[int] = ..., chunk: _Optional[int] = ...) -> None: ...

class RequestApplySnapshotChunk(_message.Message):
    __slots__ = ("index", "chunk", "sender")
    INDEX_FIELD_NUMBER: _ClassVar[int]
    CHUNK_FIELD_NUMBER: _ClassVar[int]
    SENDER_FIELD_NUMBER: _ClassVar[int]
    index: int
    chunk: bytes
    sender: str
    def __init__(self, index: _Optional[int] = ..., chunk: _Optional[bytes] = ..., sender: _Optional[str] = ...) -> None: ...

class RequestPrepareProposal(_message.Message):
    __slots__ = ("max_tx_bytes", "txs", "local_last_commit", "misbehavior", "height", "time", "next_validators_hash", "proposer_address")
    MAX_TX_BYTES_FIELD_NUMBER: _ClassVar[int]
    TXS_FIELD_NUMBER: _ClassVar[int]
    LOCAL_LAST_COMMIT_FIELD_NUMBER: _ClassVar[int]
    MISBEHAVIOR_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    NEXT_VALIDATORS_HASH_FIELD_NUMBER: _ClassVar[int]
    PROPOSER_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    max_tx_bytes: int
    txs: _containers.RepeatedScalarFieldContainer[bytes]
    local_last_commit: ExtendedCommitInfo
    misbehavior: _containers.RepeatedCompositeFieldContainer[Misbehavior]
    height: int
    time: _timestamp_pb2.Timestamp
    next_validators_hash: bytes
    proposer_address: bytes
    def __init__(self, max_tx_bytes: _Optional[int] = ..., txs: _Optional[_Iterable[bytes]] = ..., local_last_commit: _Optional[_Union[ExtendedCommitInfo, _Mapping]] = ..., misbehavior: _Optional[_Iterable[_Union[Misbehavior, _Mapping]]] = ..., height: _Optional[int] = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., next_validators_hash: _Optional[bytes] = ..., proposer_address: _Optional[bytes] = ...) -> None: ...

class RequestProcessProposal(_message.Message):
    __slots__ = ("txs", "proposed_last_commit", "misbehavior", "hash", "height", "time", "next_validators_hash", "proposer_address")
    TXS_FIELD_NUMBER: _ClassVar[int]
    PROPOSED_LAST_COMMIT_FIELD_NUMBER: _ClassVar[int]
    MISBEHAVIOR_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    NEXT_VALIDATORS_HASH_FIELD_NUMBER: _ClassVar[int]
    PROPOSER_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    txs: _containers.RepeatedScalarFieldContainer[bytes]
    proposed_last_commit: CommitInfo
    misbehavior: _containers.RepeatedCompositeFieldContainer[Misbehavior]
    hash: bytes
    height: int
    time: _timestamp_pb2.Timestamp
    next_validators_hash: bytes
    proposer_address: bytes
    def __init__(self, txs: _Optional[_Iterable[bytes]] = ..., proposed_last_commit: _Optional[_Union[CommitInfo, _Mapping]] = ..., misbehavior: _Optional[_Iterable[_Union[Misbehavior, _Mapping]]] = ..., hash: _Optional[bytes] = ..., height: _Optional[int] = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., next_validators_hash: _Optional[bytes] = ..., proposer_address: _Optional[bytes] = ...) -> None: ...

class RequestExtendVote(_message.Message):
    __slots__ = ("hash", "height", "time", "txs", "proposed_last_commit", "misbehavior", "next_validators_hash", "proposer_address")
    HASH_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    TXS_FIELD_NUMBER: _ClassVar[int]
    PROPOSED_LAST_COMMIT_FIELD_NUMBER: _ClassVar[int]
    MISBEHAVIOR_FIELD_NUMBER: _ClassVar[int]
    NEXT_VALIDATORS_HASH_FIELD_NUMBER: _ClassVar[int]
    PROPOSER_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    hash: bytes
    height: int
    time: _timestamp_pb2.Timestamp
    txs: _containers.RepeatedScalarFieldContainer[bytes]
    proposed_last_commit: CommitInfo
    misbehavior: _containers.RepeatedCompositeFieldContainer[Misbehavior]
    next_validators_hash: bytes
    proposer_address: bytes
    def __init__(self, hash: _Optional[bytes] = ..., height: _Optional[int] = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., txs: _Optional[_Iterable[bytes]] = ..., proposed_last_commit: _Optional[_Union[CommitInfo, _Mapping]] = ..., misbehavior: _Optional[_Iterable[_Union[Misbehavior, _Mapping]]] = ..., next_validators_hash: _Optional[bytes] = ..., proposer_address: _Optional[bytes] = ...) -> None: ...

class RequestVerifyVoteExtension(_message.Message):
    __slots__ = ("hash", "validator_address", "height", "vote_extension")
    HASH_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    VOTE_EXTENSION_FIELD_NUMBER: _ClassVar[int]
    hash: bytes
    validator_address: bytes
    height: int
    vote_extension: bytes
    def __init__(self, hash: _Optional[bytes] = ..., validator_address: _Optional[bytes] = ..., height: _Optional[int] = ..., vote_extension: _Optional[bytes] = ...) -> None: ...

class RequestFinalizeBlock(_message.Message):
    __slots__ = ("txs", "decided_last_commit", "misbehavior", "hash", "height", "time", "next_validators_hash", "proposer_address")
    TXS_FIELD_NUMBER: _ClassVar[int]
    DECIDED_LAST_COMMIT_FIELD_NUMBER: _ClassVar[int]
    MISBEHAVIOR_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    NEXT_VALIDATORS_HASH_FIELD_NUMBER: _ClassVar[int]
    PROPOSER_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    txs: _containers.RepeatedScalarFieldContainer[bytes]
    decided_last_commit: CommitInfo
    misbehavior: _containers.RepeatedCompositeFieldContainer[Misbehavior]
    hash: bytes
    height: int
    time: _timestamp_pb2.Timestamp
    next_validators_hash: bytes
    proposer_address: bytes
    def __init__(self, txs: _Optional[_Iterable[bytes]] = ..., decided_last_commit: _Optional[_Union[CommitInfo, _Mapping]] = ..., misbehavior: _Optional[_Iterable[_Union[Misbehavior, _Mapping]]] = ..., hash: _Optional[bytes] = ..., height: _Optional[int] = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., next_validators_hash: _Optional[bytes] = ..., proposer_address: _Optional[bytes] = ...) -> None: ...

class Response(_message.Message):
    __slots__ = ("exception", "echo", "flush", "info", "init_chain", "query", "check_tx", "commit", "list_snapshots", "offer_snapshot", "load_snapshot_chunk", "apply_snapshot_chunk", "prepare_proposal", "process_proposal", "extend_vote", "verify_vote_extension", "finalize_block")
    EXCEPTION_FIELD_NUMBER: _ClassVar[int]
    ECHO_FIELD_NUMBER: _ClassVar[int]
    FLUSH_FIELD_NUMBER: _ClassVar[int]
    INFO_FIELD_NUMBER: _ClassVar[int]
    INIT_CHAIN_FIELD_NUMBER: _ClassVar[int]
    QUERY_FIELD_NUMBER: _ClassVar[int]
    CHECK_TX_FIELD_NUMBER: _ClassVar[int]
    COMMIT_FIELD_NUMBER: _ClassVar[int]
    LIST_SNAPSHOTS_FIELD_NUMBER: _ClassVar[int]
    OFFER_SNAPSHOT_FIELD_NUMBER: _ClassVar[int]
    LOAD_SNAPSHOT_CHUNK_FIELD_NUMBER: _ClassVar[int]
    APPLY_SNAPSHOT_CHUNK_FIELD_NUMBER: _ClassVar[int]
    PREPARE_PROPOSAL_FIELD_NUMBER: _ClassVar[int]
    PROCESS_PROPOSAL_FIELD_NUMBER: _ClassVar[int]
    EXTEND_VOTE_FIELD_NUMBER: _ClassVar[int]
    VERIFY_VOTE_EXTENSION_FIELD_NUMBER: _ClassVar[int]
    FINALIZE_BLOCK_FIELD_NUMBER: _ClassVar[int]
    exception: ResponseException
    echo: ResponseEcho
    flush: ResponseFlush
    info: ResponseInfo
    init_chain: ResponseInitChain
    query: ResponseQuery
    check_tx: ResponseCheckTx
    commit: ResponseCommit
    list_snapshots: ResponseListSnapshots
    offer_snapshot: ResponseOfferSnapshot
    load_snapshot_chunk: ResponseLoadSnapshotChunk
    apply_snapshot_chunk: ResponseApplySnapshotChunk
    prepare_proposal: ResponsePrepareProposal
    process_proposal: ResponseProcessProposal
    extend_vote: ResponseExtendVote
    verify_vote_extension: ResponseVerifyVoteExtension
    finalize_block: ResponseFinalizeBlock
    def __init__(self, exception: _Optional[_Union[ResponseException, _Mapping]] = ..., echo: _Optional[_Union[ResponseEcho, _Mapping]] = ..., flush: _Optional[_Union[ResponseFlush, _Mapping]] = ..., info: _Optional[_Union[ResponseInfo, _Mapping]] = ..., init_chain: _Optional[_Union[ResponseInitChain, _Mapping]] = ..., query: _Optional[_Union[ResponseQuery, _Mapping]] = ..., check_tx: _Optional[_Union[ResponseCheckTx, _Mapping]] = ..., commit: _Optional[_Union[ResponseCommit, _Mapping]] = ..., list_snapshots: _Optional[_Union[ResponseListSnapshots, _Mapping]] = ..., offer_snapshot: _Optional[_Union[ResponseOfferSnapshot, _Mapping]] = ..., load_snapshot_chunk: _Optional[_Union[ResponseLoadSnapshotChunk, _Mapping]] = ..., apply_snapshot_chunk: _Optional[_Union[ResponseApplySnapshotChunk, _Mapping]] = ..., prepare_proposal: _Optional[_Union[ResponsePrepareProposal, _Mapping]] = ..., process_proposal: _Optional[_Union[ResponseProcessProposal, _Mapping]] = ..., extend_vote: _Optional[_Union[ResponseExtendVote, _Mapping]] = ..., verify_vote_extension: _Optional[_Union[ResponseVerifyVoteExtension, _Mapping]] = ..., finalize_block: _Optional[_Union[ResponseFinalizeBlock, _Mapping]] = ...) -> None: ...

class ResponseException(_message.Message):
    __slots__ = ("error",)
    ERROR_FIELD_NUMBER: _ClassVar[int]
    error: str
    def __init__(self, error: _Optional[str] = ...) -> None: ...

class ResponseEcho(_message.Message):
    __slots__ = ("message",)
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    message: str
    def __init__(self, message: _Optional[str] = ...) -> None: ...

class ResponseFlush(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ResponseInfo(_message.Message):
    __slots__ = ("data", "version", "app_version", "last_block_height", "last_block_app_hash")
    DATA_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    APP_VERSION_FIELD_NUMBER: _ClassVar[int]
    LAST_BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    LAST_BLOCK_APP_HASH_FIELD_NUMBER: _ClassVar[int]
    data: str
    version: str
    app_version: int
    last_block_height: int
    last_block_app_hash: bytes
    def __init__(self, data: _Optional[str] = ..., version: _Optional[str] = ..., app_version: _Optional[int] = ..., last_block_height: _Optional[int] = ..., last_block_app_hash: _Optional[bytes] = ...) -> None: ...

class ResponseInitChain(_message.Message):
    __slots__ = ("consensus_params", "validators", "app_hash")
    CONSENSUS_PARAMS_FIELD_NUMBER: _ClassVar[int]
    VALIDATORS_FIELD_NUMBER: _ClassVar[int]
    APP_HASH_FIELD_NUMBER: _ClassVar[int]
    consensus_params: _params_pb2.ConsensusParams
    validators: _containers.RepeatedCompositeFieldContainer[ValidatorUpdate]
    app_hash: bytes
    def __init__(self, consensus_params: _Optional[_Union[_params_pb2.ConsensusParams, _Mapping]] = ..., validators: _Optional[_Iterable[_Union[ValidatorUpdate, _Mapping]]] = ..., app_hash: _Optional[bytes] = ...) -> None: ...

class ResponseQuery(_message.Message):
    __slots__ = ("code", "log", "info", "index", "key", "value", "proof_ops", "height", "codespace")
    CODE_FIELD_NUMBER: _ClassVar[int]
    LOG_FIELD_NUMBER: _ClassVar[int]
    INFO_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    PROOF_OPS_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    CODESPACE_FIELD_NUMBER: _ClassVar[int]
    code: int
    log: str
    info: str
    index: int
    key: bytes
    value: bytes
    proof_ops: _proof_pb2.ProofOps
    height: int
    codespace: str
    def __init__(self, code: _Optional[int] = ..., log: _Optional[str] = ..., info: _Optional[str] = ..., index: _Optional[int] = ..., key: _Optional[bytes] = ..., value: _Optional[bytes] = ..., proof_ops: _Optional[_Union[_proof_pb2.ProofOps, _Mapping]] = ..., height: _Optional[int] = ..., codespace: _Optional[str] = ...) -> None: ...

class ResponseCheckTx(_message.Message):
    __slots__ = ("code", "data", "log", "info", "gas_wanted", "gas_used", "events", "codespace")
    CODE_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    LOG_FIELD_NUMBER: _ClassVar[int]
    INFO_FIELD_NUMBER: _ClassVar[int]
    GAS_WANTED_FIELD_NUMBER: _ClassVar[int]
    GAS_USED_FIELD_NUMBER: _ClassVar[int]
    EVENTS_FIELD_NUMBER: _ClassVar[int]
    CODESPACE_FIELD_NUMBER: _ClassVar[int]
    code: int
    data: bytes
    log: str
    info: str
    gas_wanted: int
    gas_used: int
    events: _containers.RepeatedCompositeFieldContainer[Event]
    codespace: str
    def __init__(self, code: _Optional[int] = ..., data: _Optional[bytes] = ..., log: _Optional[str] = ..., info: _Optional[str] = ..., gas_wanted: _Optional[int] = ..., gas_used: _Optional[int] = ..., events: _Optional[_Iterable[_Union[Event, _Mapping]]] = ..., codespace: _Optional[str] = ...) -> None: ...

class ResponseCommit(_message.Message):
    __slots__ = ("retain_height",)
    RETAIN_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    retain_height: int
    def __init__(self, retain_height: _Optional[int] = ...) -> None: ...

class ResponseListSnapshots(_message.Message):
    __slots__ = ("snapshots",)
    SNAPSHOTS_FIELD_NUMBER: _ClassVar[int]
    snapshots: _containers.RepeatedCompositeFieldContainer[Snapshot]
    def __init__(self, snapshots: _Optional[_Iterable[_Union[Snapshot, _Mapping]]] = ...) -> None: ...

class ResponseOfferSnapshot(_message.Message):
    __slots__ = ("result",)
    class Result(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        UNKNOWN: _ClassVar[ResponseOfferSnapshot.Result]
        ACCEPT: _ClassVar[ResponseOfferSnapshot.Result]
        ABORT: _ClassVar[ResponseOfferSnapshot.Result]
        REJECT: _ClassVar[ResponseOfferSnapshot.Result]
        REJECT_FORMAT: _ClassVar[ResponseOfferSnapshot.Result]
        REJECT_SENDER: _ClassVar[ResponseOfferSnapshot.Result]
    UNKNOWN: ResponseOfferSnapshot.Result
    ACCEPT: ResponseOfferSnapshot.Result
    ABORT: ResponseOfferSnapshot.Result
    REJECT: ResponseOfferSnapshot.Result
    REJECT_FORMAT: ResponseOfferSnapshot.Result
    REJECT_SENDER: ResponseOfferSnapshot.Result
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ResponseOfferSnapshot.Result
    def __init__(self, result: _Optional[_Union[ResponseOfferSnapshot.Result, str]] = ...) -> None: ...

class ResponseLoadSnapshotChunk(_message.Message):
    __slots__ = ("chunk",)
    CHUNK_FIELD_NUMBER: _ClassVar[int]
    chunk: bytes
    def __init__(self, chunk: _Optional[bytes] = ...) -> None: ...

class ResponseApplySnapshotChunk(_message.Message):
    __slots__ = ("result", "refetch_chunks", "reject_senders")
    class Result(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        UNKNOWN: _ClassVar[ResponseApplySnapshotChunk.Result]
        ACCEPT: _ClassVar[ResponseApplySnapshotChunk.Result]
        ABORT: _ClassVar[ResponseApplySnapshotChunk.Result]
        RETRY: _ClassVar[ResponseApplySnapshotChunk.Result]
        RETRY_SNAPSHOT: _ClassVar[ResponseApplySnapshotChunk.Result]
        REJECT_SNAPSHOT: _ClassVar[ResponseApplySnapshotChunk.Result]
    UNKNOWN: ResponseApplySnapshotChunk.Result
    ACCEPT: ResponseApplySnapshotChunk.Result
    ABORT: ResponseApplySnapshotChunk.Result
    RETRY: ResponseApplySnapshotChunk.Result
    RETRY_SNAPSHOT: ResponseApplySnapshotChunk.Result
    REJECT_SNAPSHOT: ResponseApplySnapshotChunk.Result
    RESULT_FIELD_NUMBER: _ClassVar[int]
    REFETCH_CHUNKS_FIELD_NUMBER: _ClassVar[int]
    REJECT_SENDERS_FIELD_NUMBER: _ClassVar[int]
    result: ResponseApplySnapshotChunk.Result
    refetch_chunks: _containers.RepeatedScalarFieldContainer[int]
    reject_senders: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, result: _Optional[_Union[ResponseApplySnapshotChunk.Result, str]] = ..., refetch_chunks: _Optional[_Iterable[int]] = ..., reject_senders: _Optional[_Iterable[str]] = ...) -> None: ...

class ResponsePrepareProposal(_message.Message):
    __slots__ = ("txs",)
    TXS_FIELD_NUMBER: _ClassVar[int]
    txs: _containers.RepeatedScalarFieldContainer[bytes]
    def __init__(self, txs: _Optional[_Iterable[bytes]] = ...) -> None: ...

class ResponseProcessProposal(_message.Message):
    __slots__ = ("status",)
    class ProposalStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        UNKNOWN: _ClassVar[ResponseProcessProposal.ProposalStatus]
        ACCEPT: _ClassVar[ResponseProcessProposal.ProposalStatus]
        REJECT: _ClassVar[ResponseProcessProposal.ProposalStatus]
    UNKNOWN: ResponseProcessProposal.ProposalStatus
    ACCEPT: ResponseProcessProposal.ProposalStatus
    REJECT: ResponseProcessProposal.ProposalStatus
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: ResponseProcessProposal.ProposalStatus
    def __init__(self, status: _Optional[_Union[ResponseProcessProposal.ProposalStatus, str]] = ...) -> None: ...

class ResponseExtendVote(_message.Message):
    __slots__ = ("vote_extension",)
    VOTE_EXTENSION_FIELD_NUMBER: _ClassVar[int]
    vote_extension: bytes
    def __init__(self, vote_extension: _Optional[bytes] = ...) -> None: ...

class ResponseVerifyVoteExtension(_message.Message):
    __slots__ = ("status",)
    class VerifyStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        UNKNOWN: _ClassVar[ResponseVerifyVoteExtension.VerifyStatus]
        ACCEPT: _ClassVar[ResponseVerifyVoteExtension.VerifyStatus]
        REJECT: _ClassVar[ResponseVerifyVoteExtension.VerifyStatus]
    UNKNOWN: ResponseVerifyVoteExtension.VerifyStatus
    ACCEPT: ResponseVerifyVoteExtension.VerifyStatus
    REJECT: ResponseVerifyVoteExtension.VerifyStatus
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: ResponseVerifyVoteExtension.VerifyStatus
    def __init__(self, status: _Optional[_Union[ResponseVerifyVoteExtension.VerifyStatus, str]] = ...) -> None: ...

class ResponseFinalizeBlock(_message.Message):
    __slots__ = ("events", "tx_results", "validator_updates", "consensus_param_updates", "app_hash")
    EVENTS_FIELD_NUMBER: _ClassVar[int]
    TX_RESULTS_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_UPDATES_FIELD_NUMBER: _ClassVar[int]
    CONSENSUS_PARAM_UPDATES_FIELD_NUMBER: _ClassVar[int]
    APP_HASH_FIELD_NUMBER: _ClassVar[int]
    events: _containers.RepeatedCompositeFieldContainer[Event]
    tx_results: _containers.RepeatedCompositeFieldContainer[ExecTxResult]
    validator_updates: _containers.RepeatedCompositeFieldContainer[ValidatorUpdate]
    consensus_param_updates: _params_pb2.ConsensusParams
    app_hash: bytes
    def __init__(self, events: _Optional[_Iterable[_Union[Event, _Mapping]]] = ..., tx_results: _Optional[_Iterable[_Union[ExecTxResult, _Mapping]]] = ..., validator_updates: _Optional[_Iterable[_Union[ValidatorUpdate, _Mapping]]] = ..., consensus_param_updates: _Optional[_Union[_params_pb2.ConsensusParams, _Mapping]] = ..., app_hash: _Optional[bytes] = ...) -> None: ...

class CommitInfo(_message.Message):
    __slots__ = ("round", "votes")
    ROUND_FIELD_NUMBER: _ClassVar[int]
    VOTES_FIELD_NUMBER: _ClassVar[int]
    round: int
    votes: _containers.RepeatedCompositeFieldContainer[VoteInfo]
    def __init__(self, round: _Optional[int] = ..., votes: _Optional[_Iterable[_Union[VoteInfo, _Mapping]]] = ...) -> None: ...

class ExtendedCommitInfo(_message.Message):
    __slots__ = ("round", "votes")
    ROUND_FIELD_NUMBER: _ClassVar[int]
    VOTES_FIELD_NUMBER: _ClassVar[int]
    round: int
    votes: _containers.RepeatedCompositeFieldContainer[ExtendedVoteInfo]
    def __init__(self, round: _Optional[int] = ..., votes: _Optional[_Iterable[_Union[ExtendedVoteInfo, _Mapping]]] = ...) -> None: ...

class Event(_message.Message):
    __slots__ = ("type", "attributes")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    ATTRIBUTES_FIELD_NUMBER: _ClassVar[int]
    type: str
    attributes: _containers.RepeatedCompositeFieldContainer[EventAttribute]
    def __init__(self, type: _Optional[str] = ..., attributes: _Optional[_Iterable[_Union[EventAttribute, _Mapping]]] = ...) -> None: ...

class EventAttribute(_message.Message):
    __slots__ = ("key", "value", "index")
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    key: str
    value: str
    index: bool
    def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ..., index: bool = ...) -> None: ...

class ExecTxResult(_message.Message):
    __slots__ = ("code", "data", "log", "info", "gas_wanted", "gas_used", "events", "codespace")
    CODE_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    LOG_FIELD_NUMBER: _ClassVar[int]
    INFO_FIELD_NUMBER: _ClassVar[int]
    GAS_WANTED_FIELD_NUMBER: _ClassVar[int]
    GAS_USED_FIELD_NUMBER: _ClassVar[int]
    EVENTS_FIELD_NUMBER: _ClassVar[int]
    CODESPACE_FIELD_NUMBER: _ClassVar[int]
    code: int
    data: bytes
    log: str
    info: str
    gas_wanted: int
    gas_used: int
    events: _containers.RepeatedCompositeFieldContainer[Event]
    codespace: str
    def __init__(self, code: _Optional[int] = ..., data: _Optional[bytes] = ..., log: _Optional[str] = ..., info: _Optional[str] = ..., gas_wanted: _Optional[int] = ..., gas_used: _Optional[int] = ..., events: _Optional[_Iterable[_Union[Event, _Mapping]]] = ..., codespace: _Optional[str] = ...) -> None: ...

class TxResult(_message.Message):
    __slots__ = ("height", "index", "tx", "result")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    TX_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    height: int
    index: int
    tx: bytes
    result: ExecTxResult
    def __init__(self, height: _Optional[int] = ..., index: _Optional[int] = ..., tx: _Optional[bytes] = ..., result: _Optional[_Union[ExecTxResult, _Mapping]] = ...) -> None: ...

class Validator(_message.Message):
    __slots__ = ("address", "power")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    POWER_FIELD_NUMBER: _ClassVar[int]
    address: bytes
    power: int
    def __init__(self, address: _Optional[bytes] = ..., power: _Optional[int] = ...) -> None: ...

class ValidatorUpdate(_message.Message):
    __slots__ = ("pub_key", "power")
    PUB_KEY_FIELD_NUMBER: _ClassVar[int]
    POWER_FIELD_NUMBER: _ClassVar[int]
    pub_key: _keys_pb2.PublicKey
    power: int
    def __init__(self, pub_key: _Optional[_Union[_keys_pb2.PublicKey, _Mapping]] = ..., power: _Optional[int] = ...) -> None: ...

class VoteInfo(_message.Message):
    __slots__ = ("validator", "block_id_flag")
    VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    BLOCK_ID_FLAG_FIELD_NUMBER: _ClassVar[int]
    validator: Validator
    block_id_flag: _validator_pb2.BlockIDFlag
    def __init__(self, validator: _Optional[_Union[Validator, _Mapping]] = ..., block_id_flag: _Optional[_Union[_validator_pb2.BlockIDFlag, str]] = ...) -> None: ...

class ExtendedVoteInfo(_message.Message):
    __slots__ = ("validator", "vote_extension", "extension_signature", "block_id_flag")
    VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    VOTE_EXTENSION_FIELD_NUMBER: _ClassVar[int]
    EXTENSION_SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    BLOCK_ID_FLAG_FIELD_NUMBER: _ClassVar[int]
    validator: Validator
    vote_extension: bytes
    extension_signature: bytes
    block_id_flag: _validator_pb2.BlockIDFlag
    def __init__(self, validator: _Optional[_Union[Validator, _Mapping]] = ..., vote_extension: _Optional[bytes] = ..., extension_signature: _Optional[bytes] = ..., block_id_flag: _Optional[_Union[_validator_pb2.BlockIDFlag, str]] = ...) -> None: ...

class Misbehavior(_message.Message):
    __slots__ = ("type", "validator", "height", "time", "total_voting_power")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_VOTING_POWER_FIELD_NUMBER: _ClassVar[int]
    type: MisbehaviorType
    validator: Validator
    height: int
    time: _timestamp_pb2.Timestamp
    total_voting_power: int
    def __init__(self, type: _Optional[_Union[MisbehaviorType, str]] = ..., validator: _Optional[_Union[Validator, _Mapping]] = ..., height: _Optional[int] = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., total_voting_power: _Optional[int] = ...) -> None: ...

class Snapshot(_message.Message):
    __slots__ = ("height", "format", "chunks", "hash", "metadata")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    FORMAT_FIELD_NUMBER: _ClassVar[int]
    CHUNKS_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    height: int
    format: int
    chunks: int
    hash: bytes
    metadata: bytes
    def __init__(self, height: _Optional[int] = ..., format: _Optional[int] = ..., chunks: _Optional[int] = ..., hash: _Optional[bytes] = ..., metadata: _Optional[bytes] = ...) -> None: ...
