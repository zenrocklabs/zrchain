from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf import any_pb2 as _any_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ProposalType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    PROPOSAL_TYPE_UNSPECIFIED: _ClassVar[ProposalType]
    PROPOSAL_TYPE_STANDARD: _ClassVar[ProposalType]
    PROPOSAL_TYPE_MULTIPLE_CHOICE: _ClassVar[ProposalType]
    PROPOSAL_TYPE_OPTIMISTIC: _ClassVar[ProposalType]
    PROPOSAL_TYPE_EXPEDITED: _ClassVar[ProposalType]

class VoteOption(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    VOTE_OPTION_UNSPECIFIED: _ClassVar[VoteOption]
    VOTE_OPTION_ONE: _ClassVar[VoteOption]
    VOTE_OPTION_YES: _ClassVar[VoteOption]
    VOTE_OPTION_TWO: _ClassVar[VoteOption]
    VOTE_OPTION_ABSTAIN: _ClassVar[VoteOption]
    VOTE_OPTION_THREE: _ClassVar[VoteOption]
    VOTE_OPTION_NO: _ClassVar[VoteOption]
    VOTE_OPTION_FOUR: _ClassVar[VoteOption]
    VOTE_OPTION_NO_WITH_VETO: _ClassVar[VoteOption]
    VOTE_OPTION_SPAM: _ClassVar[VoteOption]

class ProposalStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    PROPOSAL_STATUS_UNSPECIFIED: _ClassVar[ProposalStatus]
    PROPOSAL_STATUS_DEPOSIT_PERIOD: _ClassVar[ProposalStatus]
    PROPOSAL_STATUS_VOTING_PERIOD: _ClassVar[ProposalStatus]
    PROPOSAL_STATUS_PASSED: _ClassVar[ProposalStatus]
    PROPOSAL_STATUS_REJECTED: _ClassVar[ProposalStatus]
    PROPOSAL_STATUS_FAILED: _ClassVar[ProposalStatus]
PROPOSAL_TYPE_UNSPECIFIED: ProposalType
PROPOSAL_TYPE_STANDARD: ProposalType
PROPOSAL_TYPE_MULTIPLE_CHOICE: ProposalType
PROPOSAL_TYPE_OPTIMISTIC: ProposalType
PROPOSAL_TYPE_EXPEDITED: ProposalType
VOTE_OPTION_UNSPECIFIED: VoteOption
VOTE_OPTION_ONE: VoteOption
VOTE_OPTION_YES: VoteOption
VOTE_OPTION_TWO: VoteOption
VOTE_OPTION_ABSTAIN: VoteOption
VOTE_OPTION_THREE: VoteOption
VOTE_OPTION_NO: VoteOption
VOTE_OPTION_FOUR: VoteOption
VOTE_OPTION_NO_WITH_VETO: VoteOption
VOTE_OPTION_SPAM: VoteOption
PROPOSAL_STATUS_UNSPECIFIED: ProposalStatus
PROPOSAL_STATUS_DEPOSIT_PERIOD: ProposalStatus
PROPOSAL_STATUS_VOTING_PERIOD: ProposalStatus
PROPOSAL_STATUS_PASSED: ProposalStatus
PROPOSAL_STATUS_REJECTED: ProposalStatus
PROPOSAL_STATUS_FAILED: ProposalStatus

class WeightedVoteOption(_message.Message):
    __slots__ = ("option", "weight")
    OPTION_FIELD_NUMBER: _ClassVar[int]
    WEIGHT_FIELD_NUMBER: _ClassVar[int]
    option: VoteOption
    weight: str
    def __init__(self, option: _Optional[_Union[VoteOption, str]] = ..., weight: _Optional[str] = ...) -> None: ...

class Deposit(_message.Message):
    __slots__ = ("proposal_id", "depositor", "amount")
    PROPOSAL_ID_FIELD_NUMBER: _ClassVar[int]
    DEPOSITOR_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    proposal_id: int
    depositor: str
    amount: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, proposal_id: _Optional[int] = ..., depositor: _Optional[str] = ..., amount: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class Proposal(_message.Message):
    __slots__ = ("id", "messages", "status", "final_tally_result", "submit_time", "deposit_end_time", "total_deposit", "voting_start_time", "voting_end_time", "metadata", "title", "summary", "proposer", "expedited", "failed_reason", "proposal_type")
    ID_FIELD_NUMBER: _ClassVar[int]
    MESSAGES_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    FINAL_TALLY_RESULT_FIELD_NUMBER: _ClassVar[int]
    SUBMIT_TIME_FIELD_NUMBER: _ClassVar[int]
    DEPOSIT_END_TIME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_DEPOSIT_FIELD_NUMBER: _ClassVar[int]
    VOTING_START_TIME_FIELD_NUMBER: _ClassVar[int]
    VOTING_END_TIME_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    PROPOSER_FIELD_NUMBER: _ClassVar[int]
    EXPEDITED_FIELD_NUMBER: _ClassVar[int]
    FAILED_REASON_FIELD_NUMBER: _ClassVar[int]
    PROPOSAL_TYPE_FIELD_NUMBER: _ClassVar[int]
    id: int
    messages: _containers.RepeatedCompositeFieldContainer[_any_pb2.Any]
    status: ProposalStatus
    final_tally_result: TallyResult
    submit_time: _timestamp_pb2.Timestamp
    deposit_end_time: _timestamp_pb2.Timestamp
    total_deposit: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    voting_start_time: _timestamp_pb2.Timestamp
    voting_end_time: _timestamp_pb2.Timestamp
    metadata: str
    title: str
    summary: str
    proposer: str
    expedited: bool
    failed_reason: str
    proposal_type: ProposalType
    def __init__(self, id: _Optional[int] = ..., messages: _Optional[_Iterable[_Union[_any_pb2.Any, _Mapping]]] = ..., status: _Optional[_Union[ProposalStatus, str]] = ..., final_tally_result: _Optional[_Union[TallyResult, _Mapping]] = ..., submit_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., deposit_end_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., total_deposit: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., voting_start_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., voting_end_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., metadata: _Optional[str] = ..., title: _Optional[str] = ..., summary: _Optional[str] = ..., proposer: _Optional[str] = ..., expedited: bool = ..., failed_reason: _Optional[str] = ..., proposal_type: _Optional[_Union[ProposalType, str]] = ...) -> None: ...

class ProposalVoteOptions(_message.Message):
    __slots__ = ("option_one", "option_two", "option_three", "option_four", "option_spam")
    OPTION_ONE_FIELD_NUMBER: _ClassVar[int]
    OPTION_TWO_FIELD_NUMBER: _ClassVar[int]
    OPTION_THREE_FIELD_NUMBER: _ClassVar[int]
    OPTION_FOUR_FIELD_NUMBER: _ClassVar[int]
    OPTION_SPAM_FIELD_NUMBER: _ClassVar[int]
    option_one: str
    option_two: str
    option_three: str
    option_four: str
    option_spam: str
    def __init__(self, option_one: _Optional[str] = ..., option_two: _Optional[str] = ..., option_three: _Optional[str] = ..., option_four: _Optional[str] = ..., option_spam: _Optional[str] = ...) -> None: ...

class TallyResult(_message.Message):
    __slots__ = ("yes_count", "abstain_count", "no_count", "no_with_veto_count", "option_one_count", "option_two_count", "option_three_count", "option_four_count", "spam_count")
    YES_COUNT_FIELD_NUMBER: _ClassVar[int]
    ABSTAIN_COUNT_FIELD_NUMBER: _ClassVar[int]
    NO_COUNT_FIELD_NUMBER: _ClassVar[int]
    NO_WITH_VETO_COUNT_FIELD_NUMBER: _ClassVar[int]
    OPTION_ONE_COUNT_FIELD_NUMBER: _ClassVar[int]
    OPTION_TWO_COUNT_FIELD_NUMBER: _ClassVar[int]
    OPTION_THREE_COUNT_FIELD_NUMBER: _ClassVar[int]
    OPTION_FOUR_COUNT_FIELD_NUMBER: _ClassVar[int]
    SPAM_COUNT_FIELD_NUMBER: _ClassVar[int]
    yes_count: str
    abstain_count: str
    no_count: str
    no_with_veto_count: str
    option_one_count: str
    option_two_count: str
    option_three_count: str
    option_four_count: str
    spam_count: str
    def __init__(self, yes_count: _Optional[str] = ..., abstain_count: _Optional[str] = ..., no_count: _Optional[str] = ..., no_with_veto_count: _Optional[str] = ..., option_one_count: _Optional[str] = ..., option_two_count: _Optional[str] = ..., option_three_count: _Optional[str] = ..., option_four_count: _Optional[str] = ..., spam_count: _Optional[str] = ...) -> None: ...

class Vote(_message.Message):
    __slots__ = ("proposal_id", "voter", "options", "metadata")
    PROPOSAL_ID_FIELD_NUMBER: _ClassVar[int]
    VOTER_FIELD_NUMBER: _ClassVar[int]
    OPTIONS_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    proposal_id: int
    voter: str
    options: _containers.RepeatedCompositeFieldContainer[WeightedVoteOption]
    metadata: str
    def __init__(self, proposal_id: _Optional[int] = ..., voter: _Optional[str] = ..., options: _Optional[_Iterable[_Union[WeightedVoteOption, _Mapping]]] = ..., metadata: _Optional[str] = ...) -> None: ...

class DepositParams(_message.Message):
    __slots__ = ("min_deposit", "max_deposit_period")
    MIN_DEPOSIT_FIELD_NUMBER: _ClassVar[int]
    MAX_DEPOSIT_PERIOD_FIELD_NUMBER: _ClassVar[int]
    min_deposit: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    max_deposit_period: _duration_pb2.Duration
    def __init__(self, min_deposit: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., max_deposit_period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ...) -> None: ...

class VotingParams(_message.Message):
    __slots__ = ("voting_period",)
    VOTING_PERIOD_FIELD_NUMBER: _ClassVar[int]
    voting_period: _duration_pb2.Duration
    def __init__(self, voting_period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ...) -> None: ...

class TallyParams(_message.Message):
    __slots__ = ("quorum", "threshold", "veto_threshold")
    QUORUM_FIELD_NUMBER: _ClassVar[int]
    THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    VETO_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    quorum: str
    threshold: str
    veto_threshold: str
    def __init__(self, quorum: _Optional[str] = ..., threshold: _Optional[str] = ..., veto_threshold: _Optional[str] = ...) -> None: ...

class Params(_message.Message):
    __slots__ = ("min_deposit", "max_deposit_period", "voting_period", "quorum", "threshold", "veto_threshold", "min_initial_deposit_ratio", "proposal_cancel_ratio", "proposal_cancel_dest", "expedited_voting_period", "expedited_threshold", "expedited_min_deposit", "burn_vote_quorum", "burn_proposal_deposit_prevote", "burn_vote_veto", "min_deposit_ratio", "proposal_cancel_max_period", "optimistic_authorized_addresses", "optimistic_rejected_threshold", "yes_quorum")
    MIN_DEPOSIT_FIELD_NUMBER: _ClassVar[int]
    MAX_DEPOSIT_PERIOD_FIELD_NUMBER: _ClassVar[int]
    VOTING_PERIOD_FIELD_NUMBER: _ClassVar[int]
    QUORUM_FIELD_NUMBER: _ClassVar[int]
    THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    VETO_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    MIN_INITIAL_DEPOSIT_RATIO_FIELD_NUMBER: _ClassVar[int]
    PROPOSAL_CANCEL_RATIO_FIELD_NUMBER: _ClassVar[int]
    PROPOSAL_CANCEL_DEST_FIELD_NUMBER: _ClassVar[int]
    EXPEDITED_VOTING_PERIOD_FIELD_NUMBER: _ClassVar[int]
    EXPEDITED_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    EXPEDITED_MIN_DEPOSIT_FIELD_NUMBER: _ClassVar[int]
    BURN_VOTE_QUORUM_FIELD_NUMBER: _ClassVar[int]
    BURN_PROPOSAL_DEPOSIT_PREVOTE_FIELD_NUMBER: _ClassVar[int]
    BURN_VOTE_VETO_FIELD_NUMBER: _ClassVar[int]
    MIN_DEPOSIT_RATIO_FIELD_NUMBER: _ClassVar[int]
    PROPOSAL_CANCEL_MAX_PERIOD_FIELD_NUMBER: _ClassVar[int]
    OPTIMISTIC_AUTHORIZED_ADDRESSES_FIELD_NUMBER: _ClassVar[int]
    OPTIMISTIC_REJECTED_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    YES_QUORUM_FIELD_NUMBER: _ClassVar[int]
    min_deposit: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    max_deposit_period: _duration_pb2.Duration
    voting_period: _duration_pb2.Duration
    quorum: str
    threshold: str
    veto_threshold: str
    min_initial_deposit_ratio: str
    proposal_cancel_ratio: str
    proposal_cancel_dest: str
    expedited_voting_period: _duration_pb2.Duration
    expedited_threshold: str
    expedited_min_deposit: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    burn_vote_quorum: bool
    burn_proposal_deposit_prevote: bool
    burn_vote_veto: bool
    min_deposit_ratio: str
    proposal_cancel_max_period: str
    optimistic_authorized_addresses: _containers.RepeatedScalarFieldContainer[str]
    optimistic_rejected_threshold: str
    yes_quorum: str
    def __init__(self, min_deposit: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., max_deposit_period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., voting_period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., quorum: _Optional[str] = ..., threshold: _Optional[str] = ..., veto_threshold: _Optional[str] = ..., min_initial_deposit_ratio: _Optional[str] = ..., proposal_cancel_ratio: _Optional[str] = ..., proposal_cancel_dest: _Optional[str] = ..., expedited_voting_period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., expedited_threshold: _Optional[str] = ..., expedited_min_deposit: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., burn_vote_quorum: bool = ..., burn_proposal_deposit_prevote: bool = ..., burn_vote_veto: bool = ..., min_deposit_ratio: _Optional[str] = ..., proposal_cancel_max_period: _Optional[str] = ..., optimistic_authorized_addresses: _Optional[_Iterable[str]] = ..., optimistic_rejected_threshold: _Optional[str] = ..., yes_quorum: _Optional[str] = ...) -> None: ...

class MessageBasedParams(_message.Message):
    __slots__ = ("voting_period", "quorum", "yes_quorum", "threshold", "veto_threshold")
    VOTING_PERIOD_FIELD_NUMBER: _ClassVar[int]
    QUORUM_FIELD_NUMBER: _ClassVar[int]
    YES_QUORUM_FIELD_NUMBER: _ClassVar[int]
    THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    VETO_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    voting_period: _duration_pb2.Duration
    quorum: str
    yes_quorum: str
    threshold: str
    veto_threshold: str
    def __init__(self, voting_period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., quorum: _Optional[str] = ..., yes_quorum: _Optional[str] = ..., threshold: _Optional[str] = ..., veto_threshold: _Optional[str] = ...) -> None: ...
