from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmos.gov.v1 import gov_pb2 as _gov_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from google.protobuf import any_pb2 as _any_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgSubmitProposal(_message.Message):
    __slots__ = ("messages", "initial_deposit", "proposer", "metadata", "title", "summary", "expedited", "proposal_type")
    MESSAGES_FIELD_NUMBER: _ClassVar[int]
    INITIAL_DEPOSIT_FIELD_NUMBER: _ClassVar[int]
    PROPOSER_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    EXPEDITED_FIELD_NUMBER: _ClassVar[int]
    PROPOSAL_TYPE_FIELD_NUMBER: _ClassVar[int]
    messages: _containers.RepeatedCompositeFieldContainer[_any_pb2.Any]
    initial_deposit: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    proposer: str
    metadata: str
    title: str
    summary: str
    expedited: bool
    proposal_type: _gov_pb2.ProposalType
    def __init__(self, messages: _Optional[_Iterable[_Union[_any_pb2.Any, _Mapping]]] = ..., initial_deposit: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., proposer: _Optional[str] = ..., metadata: _Optional[str] = ..., title: _Optional[str] = ..., summary: _Optional[str] = ..., expedited: bool = ..., proposal_type: _Optional[_Union[_gov_pb2.ProposalType, str]] = ...) -> None: ...

class MsgSubmitProposalResponse(_message.Message):
    __slots__ = ("proposal_id",)
    PROPOSAL_ID_FIELD_NUMBER: _ClassVar[int]
    proposal_id: int
    def __init__(self, proposal_id: _Optional[int] = ...) -> None: ...

class MsgExecLegacyContent(_message.Message):
    __slots__ = ("content", "authority")
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    content: _any_pb2.Any
    authority: str
    def __init__(self, content: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., authority: _Optional[str] = ...) -> None: ...

class MsgExecLegacyContentResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgVote(_message.Message):
    __slots__ = ("proposal_id", "voter", "option", "metadata")
    PROPOSAL_ID_FIELD_NUMBER: _ClassVar[int]
    VOTER_FIELD_NUMBER: _ClassVar[int]
    OPTION_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    proposal_id: int
    voter: str
    option: _gov_pb2.VoteOption
    metadata: str
    def __init__(self, proposal_id: _Optional[int] = ..., voter: _Optional[str] = ..., option: _Optional[_Union[_gov_pb2.VoteOption, str]] = ..., metadata: _Optional[str] = ...) -> None: ...

class MsgVoteResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgVoteWeighted(_message.Message):
    __slots__ = ("proposal_id", "voter", "options", "metadata")
    PROPOSAL_ID_FIELD_NUMBER: _ClassVar[int]
    VOTER_FIELD_NUMBER: _ClassVar[int]
    OPTIONS_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    proposal_id: int
    voter: str
    options: _containers.RepeatedCompositeFieldContainer[_gov_pb2.WeightedVoteOption]
    metadata: str
    def __init__(self, proposal_id: _Optional[int] = ..., voter: _Optional[str] = ..., options: _Optional[_Iterable[_Union[_gov_pb2.WeightedVoteOption, _Mapping]]] = ..., metadata: _Optional[str] = ...) -> None: ...

class MsgVoteWeightedResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgDeposit(_message.Message):
    __slots__ = ("proposal_id", "depositor", "amount")
    PROPOSAL_ID_FIELD_NUMBER: _ClassVar[int]
    DEPOSITOR_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    proposal_id: int
    depositor: str
    amount: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, proposal_id: _Optional[int] = ..., depositor: _Optional[str] = ..., amount: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class MsgDepositResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    params: _gov_pb2.Params
    def __init__(self, authority: _Optional[str] = ..., params: _Optional[_Union[_gov_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgCancelProposal(_message.Message):
    __slots__ = ("proposal_id", "proposer")
    PROPOSAL_ID_FIELD_NUMBER: _ClassVar[int]
    PROPOSER_FIELD_NUMBER: _ClassVar[int]
    proposal_id: int
    proposer: str
    def __init__(self, proposal_id: _Optional[int] = ..., proposer: _Optional[str] = ...) -> None: ...

class MsgCancelProposalResponse(_message.Message):
    __slots__ = ("proposal_id", "canceled_time", "canceled_height")
    PROPOSAL_ID_FIELD_NUMBER: _ClassVar[int]
    CANCELED_TIME_FIELD_NUMBER: _ClassVar[int]
    CANCELED_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    proposal_id: int
    canceled_time: _timestamp_pb2.Timestamp
    canceled_height: int
    def __init__(self, proposal_id: _Optional[int] = ..., canceled_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., canceled_height: _Optional[int] = ...) -> None: ...

class MsgSubmitMultipleChoiceProposal(_message.Message):
    __slots__ = ("initial_deposit", "proposer", "metadata", "title", "summary", "vote_options")
    INITIAL_DEPOSIT_FIELD_NUMBER: _ClassVar[int]
    PROPOSER_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    VOTE_OPTIONS_FIELD_NUMBER: _ClassVar[int]
    initial_deposit: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    proposer: str
    metadata: str
    title: str
    summary: str
    vote_options: _gov_pb2.ProposalVoteOptions
    def __init__(self, initial_deposit: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., proposer: _Optional[str] = ..., metadata: _Optional[str] = ..., title: _Optional[str] = ..., summary: _Optional[str] = ..., vote_options: _Optional[_Union[_gov_pb2.ProposalVoteOptions, _Mapping]] = ...) -> None: ...

class MsgSubmitMultipleChoiceProposalResponse(_message.Message):
    __slots__ = ("proposal_id",)
    PROPOSAL_ID_FIELD_NUMBER: _ClassVar[int]
    proposal_id: int
    def __init__(self, proposal_id: _Optional[int] = ...) -> None: ...

class MsgUpdateMessageParams(_message.Message):
    __slots__ = ("authority", "msg_url", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    MSG_URL_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    msg_url: str
    params: _gov_pb2.MessageBasedParams
    def __init__(self, authority: _Optional[str] = ..., msg_url: _Optional[str] = ..., params: _Optional[_Union[_gov_pb2.MessageBasedParams, _Mapping]] = ...) -> None: ...

class MsgUpdateMessageParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgSudoExec(_message.Message):
    __slots__ = ("authority", "msg")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    authority: str
    msg: _any_pb2.Any
    def __init__(self, authority: _Optional[str] = ..., msg: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class MsgSudoExecResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: bytes
    def __init__(self, result: _Optional[bytes] = ...) -> None: ...
