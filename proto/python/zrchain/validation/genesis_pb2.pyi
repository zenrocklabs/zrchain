from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.validation import staking_pb2 as _staking_pb2
from zrchain.validation import hybrid_validation_pb2 as _hybrid_validation_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("params", "last_total_power", "last_validator_powers", "validators", "delegations", "unbonding_delegations", "redelegations", "exported", "HVParams")
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    LAST_TOTAL_POWER_FIELD_NUMBER: _ClassVar[int]
    LAST_VALIDATOR_POWERS_FIELD_NUMBER: _ClassVar[int]
    VALIDATORS_FIELD_NUMBER: _ClassVar[int]
    DELEGATIONS_FIELD_NUMBER: _ClassVar[int]
    UNBONDING_DELEGATIONS_FIELD_NUMBER: _ClassVar[int]
    REDELEGATIONS_FIELD_NUMBER: _ClassVar[int]
    EXPORTED_FIELD_NUMBER: _ClassVar[int]
    HVPARAMS_FIELD_NUMBER: _ClassVar[int]
    params: _staking_pb2.Params
    last_total_power: bytes
    last_validator_powers: _containers.RepeatedCompositeFieldContainer[LastValidatorPower]
    validators: _containers.RepeatedCompositeFieldContainer[_hybrid_validation_pb2.ValidatorHV]
    delegations: _containers.RepeatedCompositeFieldContainer[_staking_pb2.Delegation]
    unbonding_delegations: _containers.RepeatedCompositeFieldContainer[_staking_pb2.UnbondingDelegation]
    redelegations: _containers.RepeatedCompositeFieldContainer[_staking_pb2.Redelegation]
    exported: bool
    HVParams: _hybrid_validation_pb2.HVParams
    def __init__(self, params: _Optional[_Union[_staking_pb2.Params, _Mapping]] = ..., last_total_power: _Optional[bytes] = ..., last_validator_powers: _Optional[_Iterable[_Union[LastValidatorPower, _Mapping]]] = ..., validators: _Optional[_Iterable[_Union[_hybrid_validation_pb2.ValidatorHV, _Mapping]]] = ..., delegations: _Optional[_Iterable[_Union[_staking_pb2.Delegation, _Mapping]]] = ..., unbonding_delegations: _Optional[_Iterable[_Union[_staking_pb2.UnbondingDelegation, _Mapping]]] = ..., redelegations: _Optional[_Iterable[_Union[_staking_pb2.Redelegation, _Mapping]]] = ..., exported: bool = ..., HVParams: _Optional[_Union[_hybrid_validation_pb2.HVParams, _Mapping]] = ...) -> None: ...

class LastValidatorPower(_message.Message):
    __slots__ = ("address", "power")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    POWER_FIELD_NUMBER: _ClassVar[int]
    address: str
    power: int
    def __init__(self, address: _Optional[str] = ..., power: _Optional[int] = ...) -> None: ...
