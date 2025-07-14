from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.validation import staking_pb2 as _staking_pb2
from zrchain.validation import solana_pb2 as _solana_pb2
from zrchain.validation import tx_pb2 as _tx_pb2
from zrchain.validation import asset_data_pb2 as _asset_data_pb2
from zrchain.validation import hybrid_validation_pb2 as _hybrid_validation_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from amino import amino_pb2 as _amino_pb2
from api import sidecar_service_pb2 as _sidecar_service_pb2
from zrchain.zenbtc import mint_pb2 as _mint_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("params", "last_total_power", "last_validator_powers", "validators", "delegations", "unbonding_delegations", "redelegations", "exported", "HVParams", "asset_prices", "last_valid_ve_height", "slash_events", "slash_event_count", "validation_infos", "btc_block_headers", "last_used_solana_nonce", "backfill_request", "last_used_ethereum_nonce", "requested_historical_bitcoin_headers", "avs_rewards_pool", "ethereum_nonce_requested", "solana_nonce_requested", "solana_zentp_accounts_requested", "solana_accounts_requested", "validator_mismatch_counts", "last_completed_zentp_mint_id")
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    LAST_TOTAL_POWER_FIELD_NUMBER: _ClassVar[int]
    LAST_VALIDATOR_POWERS_FIELD_NUMBER: _ClassVar[int]
    VALIDATORS_FIELD_NUMBER: _ClassVar[int]
    DELEGATIONS_FIELD_NUMBER: _ClassVar[int]
    UNBONDING_DELEGATIONS_FIELD_NUMBER: _ClassVar[int]
    REDELEGATIONS_FIELD_NUMBER: _ClassVar[int]
    EXPORTED_FIELD_NUMBER: _ClassVar[int]
    HVPARAMS_FIELD_NUMBER: _ClassVar[int]
    ASSET_PRICES_FIELD_NUMBER: _ClassVar[int]
    LAST_VALID_VE_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SLASH_EVENTS_FIELD_NUMBER: _ClassVar[int]
    SLASH_EVENT_COUNT_FIELD_NUMBER: _ClassVar[int]
    VALIDATION_INFOS_FIELD_NUMBER: _ClassVar[int]
    BTC_BLOCK_HEADERS_FIELD_NUMBER: _ClassVar[int]
    LAST_USED_SOLANA_NONCE_FIELD_NUMBER: _ClassVar[int]
    BACKFILL_REQUEST_FIELD_NUMBER: _ClassVar[int]
    LAST_USED_ETHEREUM_NONCE_FIELD_NUMBER: _ClassVar[int]
    REQUESTED_HISTORICAL_BITCOIN_HEADERS_FIELD_NUMBER: _ClassVar[int]
    AVS_REWARDS_POOL_FIELD_NUMBER: _ClassVar[int]
    ETHEREUM_NONCE_REQUESTED_FIELD_NUMBER: _ClassVar[int]
    SOLANA_NONCE_REQUESTED_FIELD_NUMBER: _ClassVar[int]
    SOLANA_ZENTP_ACCOUNTS_REQUESTED_FIELD_NUMBER: _ClassVar[int]
    SOLANA_ACCOUNTS_REQUESTED_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_MISMATCH_COUNTS_FIELD_NUMBER: _ClassVar[int]
    LAST_COMPLETED_ZENTP_MINT_ID_FIELD_NUMBER: _ClassVar[int]
    params: _staking_pb2.Params
    last_total_power: bytes
    last_validator_powers: _containers.RepeatedCompositeFieldContainer[LastValidatorPower]
    validators: _containers.RepeatedCompositeFieldContainer[_hybrid_validation_pb2.ValidatorHV]
    delegations: _containers.RepeatedCompositeFieldContainer[_staking_pb2.Delegation]
    unbonding_delegations: _containers.RepeatedCompositeFieldContainer[_staking_pb2.UnbondingDelegation]
    redelegations: _containers.RepeatedCompositeFieldContainer[_staking_pb2.Redelegation]
    exported: bool
    HVParams: _hybrid_validation_pb2.HVParams
    asset_prices: _containers.RepeatedCompositeFieldContainer[_asset_data_pb2.AssetData]
    last_valid_ve_height: int
    slash_events: _containers.RepeatedCompositeFieldContainer[_hybrid_validation_pb2.SlashEvent]
    slash_event_count: int
    validation_infos: _containers.RepeatedCompositeFieldContainer[_hybrid_validation_pb2.ValidationInfo]
    btc_block_headers: _containers.RepeatedCompositeFieldContainer[_sidecar_service_pb2.BTCBlockHeader]
    last_used_solana_nonce: _containers.RepeatedCompositeFieldContainer[_solana_pb2.SolanaNonce]
    backfill_request: _tx_pb2.BackfillRequests
    last_used_ethereum_nonce: _containers.RepeatedCompositeFieldContainer[_mint_pb2.NonceData]
    requested_historical_bitcoin_headers: _containers.RepeatedCompositeFieldContainer[_mint_pb2.RequestedBitcoinHeaders]
    avs_rewards_pool: _containers.RepeatedScalarFieldContainer[str]
    ethereum_nonce_requested: _containers.RepeatedScalarFieldContainer[int]
    solana_nonce_requested: _containers.RepeatedScalarFieldContainer[int]
    solana_zentp_accounts_requested: _containers.RepeatedScalarFieldContainer[str]
    solana_accounts_requested: _containers.RepeatedScalarFieldContainer[str]
    validator_mismatch_counts: _containers.RepeatedCompositeFieldContainer[_hybrid_validation_pb2.ValidatorMismatchCount]
    last_completed_zentp_mint_id: int
    def __init__(self, params: _Optional[_Union[_staking_pb2.Params, _Mapping]] = ..., last_total_power: _Optional[bytes] = ..., last_validator_powers: _Optional[_Iterable[_Union[LastValidatorPower, _Mapping]]] = ..., validators: _Optional[_Iterable[_Union[_hybrid_validation_pb2.ValidatorHV, _Mapping]]] = ..., delegations: _Optional[_Iterable[_Union[_staking_pb2.Delegation, _Mapping]]] = ..., unbonding_delegations: _Optional[_Iterable[_Union[_staking_pb2.UnbondingDelegation, _Mapping]]] = ..., redelegations: _Optional[_Iterable[_Union[_staking_pb2.Redelegation, _Mapping]]] = ..., exported: bool = ..., HVParams: _Optional[_Union[_hybrid_validation_pb2.HVParams, _Mapping]] = ..., asset_prices: _Optional[_Iterable[_Union[_asset_data_pb2.AssetData, _Mapping]]] = ..., last_valid_ve_height: _Optional[int] = ..., slash_events: _Optional[_Iterable[_Union[_hybrid_validation_pb2.SlashEvent, _Mapping]]] = ..., slash_event_count: _Optional[int] = ..., validation_infos: _Optional[_Iterable[_Union[_hybrid_validation_pb2.ValidationInfo, _Mapping]]] = ..., btc_block_headers: _Optional[_Iterable[_Union[_sidecar_service_pb2.BTCBlockHeader, _Mapping]]] = ..., last_used_solana_nonce: _Optional[_Iterable[_Union[_solana_pb2.SolanaNonce, _Mapping]]] = ..., backfill_request: _Optional[_Union[_tx_pb2.BackfillRequests, _Mapping]] = ..., last_used_ethereum_nonce: _Optional[_Iterable[_Union[_mint_pb2.NonceData, _Mapping]]] = ..., requested_historical_bitcoin_headers: _Optional[_Iterable[_Union[_mint_pb2.RequestedBitcoinHeaders, _Mapping]]] = ..., avs_rewards_pool: _Optional[_Iterable[str]] = ..., ethereum_nonce_requested: _Optional[_Iterable[int]] = ..., solana_nonce_requested: _Optional[_Iterable[int]] = ..., solana_zentp_accounts_requested: _Optional[_Iterable[str]] = ..., solana_accounts_requested: _Optional[_Iterable[str]] = ..., validator_mismatch_counts: _Optional[_Iterable[_Union[_hybrid_validation_pb2.ValidatorMismatchCount, _Mapping]]] = ..., last_completed_zentp_mint_id: _Optional[int] = ...) -> None: ...

class LastValidatorPower(_message.Message):
    __slots__ = ("address", "power")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    POWER_FIELD_NUMBER: _ClassVar[int]
    address: str
    power: int
    def __init__(self, address: _Optional[str] = ..., power: _Optional[int] = ...) -> None: ...
