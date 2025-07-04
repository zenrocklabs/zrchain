import { GeneratedType } from "@cosmjs/proto-signing";
import { ValidatorHV } from "./types/zrchain/validation/hybrid_validation";
import { Pool } from "./types/zrchain/validation/staking";
import { HVParams } from "./types/zrchain/validation/hybrid_validation";
import { QueryValidatorsRequest } from "./types/zrchain/validation/query";
import { QueryValidatorUnbondingDelegationsResponse } from "./types/zrchain/validation/query";
import { QueryPowerResponse } from "./types/zrchain/validation/query";
import { RedelegationEntryResponse } from "./types/zrchain/validation/staking";
import { GenesisState } from "./types/zrchain/validation/genesis";
import { QueryUnbondingDelegationResponse } from "./types/zrchain/validation/query";
import { QueryDelegatorUnbondingDelegationsResponse } from "./types/zrchain/validation/query";
import { QueryDelegatorValidatorsResponse } from "./types/zrchain/validation/query";
import { QueryPoolResponse } from "./types/zrchain/validation/query";
import { QueryHistoricalInfoResponse } from "./types/zrchain/validation/query";
import { CommissionRates } from "./types/zrchain/validation/staking";
import { QueryValidatorResponse } from "./types/zrchain/validation/query";
import { QueryHistoricalInfoRequest } from "./types/zrchain/validation/query";
import { QueryPoolRequest } from "./types/zrchain/validation/query";
import { ValAddresses } from "./types/zrchain/validation/staking";
import { LastValidatorPower } from "./types/zrchain/validation/genesis";
import { QueryValidatorUnbondingDelegationsRequest } from "./types/zrchain/validation/query";
import { QueryDelegatorDelegationsResponse } from "./types/zrchain/validation/query";
import { UnbondingDelegationEntry } from "./types/zrchain/validation/staking";
import { AssetData } from "./types/zrchain/validation/asset_data";
import { MsgDelegate } from "./types/zrchain/validation/tx";
import { QueryDelegatorValidatorsRequest } from "./types/zrchain/validation/query";
import { QueryDelegatorValidatorResponse } from "./types/zrchain/validation/query";
import { DVPair } from "./types/zrchain/validation/staking";
import { MsgCancelUnbondingDelegationResponse } from "./types/zrchain/validation/tx";
import { MsgUpdateParamsResponse } from "./types/zrchain/validation/tx";
import { QueryDelegationResponse } from "./types/zrchain/validation/query";
import { QueryUnbondingDelegationRequest } from "./types/zrchain/validation/query";
import { QueryParamsRequest } from "./types/zrchain/validation/query";
import { StakeAuthorization_Validators } from "./types/zrchain/validation/authz";
import { MsgCreateValidator } from "./types/zrchain/validation/tx";
import { SlashEvent } from "./types/zrchain/validation/hybrid_validation";
import { MsgEditValidator } from "./types/zrchain/validation/tx";
import { MsgCancelUnbondingDelegation } from "./types/zrchain/validation/tx";
import { QueryValidatorsResponse } from "./types/zrchain/validation/query";
import { QueryValidatorRequest } from "./types/zrchain/validation/query";
import { QueryRedelegationsResponse } from "./types/zrchain/validation/query";
import { QueryPowerRequest } from "./types/zrchain/validation/query";
import { ValidationInfo } from "./types/zrchain/validation/hybrid_validation";
import { MsgCreateValidatorResponse } from "./types/zrchain/validation/tx";
import { MsgTriggerEventBackfillResponse } from "./types/zrchain/validation/tx";
import { QueryDelegationRequest } from "./types/zrchain/validation/query";
import { QueryDelegatorUnbondingDelegationsRequest } from "./types/zrchain/validation/query";
import { QueryBackfillRequestsRequest } from "./types/zrchain/validation/query";
import { Description } from "./types/zrchain/validation/staking";
import { Params } from "./types/zrchain/validation/staking";
import { MsgUpdateParams } from "./types/zrchain/validation/tx";
import { BackfillRequests } from "./types/zrchain/validation/tx";
import { QueryDelegatorValidatorRequest } from "./types/zrchain/validation/query";
import { DVPairs } from "./types/zrchain/validation/staking";
import { DVVTriplet } from "./types/zrchain/validation/staking";
import { RedelegationEntry } from "./types/zrchain/validation/staking";
import { MsgBeginRedelegateResponse } from "./types/zrchain/validation/tx";
import { MsgUpdateHVParams } from "./types/zrchain/validation/tx";
import { Commission } from "./types/zrchain/validation/staking";
import { QueryDelegatorDelegationsRequest } from "./types/zrchain/validation/query";
import { DelegationResponse } from "./types/zrchain/validation/staking";
import { RedelegationResponse } from "./types/zrchain/validation/staking";
import { MsgDelegateResponse } from "./types/zrchain/validation/tx";
import { MsgUndelegate } from "./types/zrchain/validation/tx";
import { QueryValidatorDelegationsRequest } from "./types/zrchain/validation/query";
import { ValidatorPower } from "./types/zrchain/validation/query";
import { SolanaNonce } from "./types/zrchain/validation/solana";
import { StakeAuthorization } from "./types/zrchain/validation/authz";
import { QueryParamsResponse } from "./types/zrchain/validation/query";
import { QueryBackfillRequestsResponse } from "./types/zrchain/validation/query";
import { Delegation } from "./types/zrchain/validation/staking";
import { UnbondingDelegation } from "./types/zrchain/validation/staking";
import { Redelegation } from "./types/zrchain/validation/staking";
import { ValidatorUpdates } from "./types/zrchain/validation/staking";
import { HistoricalInfoHV } from "./types/zrchain/validation/hybrid_validation";
import { MsgUndelegateResponse } from "./types/zrchain/validation/tx";
import { QueryValidatorDelegationsResponse } from "./types/zrchain/validation/query";
import { QueryRedelegationsRequest } from "./types/zrchain/validation/query";
import { MsgEditValidatorResponse } from "./types/zrchain/validation/tx";
import { MsgBeginRedelegate } from "./types/zrchain/validation/tx";
import { MsgUpdateHVParamsResponse } from "./types/zrchain/validation/tx";
import { MsgTriggerEventBackfill } from "./types/zrchain/validation/tx";
import { DVVTriplets } from "./types/zrchain/validation/staking";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.validation.ValidatorHV", ValidatorHV],
    ["/zrchain.validation.Pool", Pool],
    ["/zrchain.validation.HVParams", HVParams],
    ["/zrchain.validation.QueryValidatorsRequest", QueryValidatorsRequest],
    ["/zrchain.validation.QueryValidatorUnbondingDelegationsResponse", QueryValidatorUnbondingDelegationsResponse],
    ["/zrchain.validation.QueryPowerResponse", QueryPowerResponse],
    ["/zrchain.validation.RedelegationEntryResponse", RedelegationEntryResponse],
    ["/zrchain.validation.GenesisState", GenesisState],
    ["/zrchain.validation.QueryUnbondingDelegationResponse", QueryUnbondingDelegationResponse],
    ["/zrchain.validation.QueryDelegatorUnbondingDelegationsResponse", QueryDelegatorUnbondingDelegationsResponse],
    ["/zrchain.validation.QueryDelegatorValidatorsResponse", QueryDelegatorValidatorsResponse],
    ["/zrchain.validation.QueryPoolResponse", QueryPoolResponse],
    ["/zrchain.validation.QueryHistoricalInfoResponse", QueryHistoricalInfoResponse],
    ["/zrchain.validation.CommissionRates", CommissionRates],
    ["/zrchain.validation.QueryValidatorResponse", QueryValidatorResponse],
    ["/zrchain.validation.QueryHistoricalInfoRequest", QueryHistoricalInfoRequest],
    ["/zrchain.validation.QueryPoolRequest", QueryPoolRequest],
    ["/zrchain.validation.ValAddresses", ValAddresses],
    ["/zrchain.validation.LastValidatorPower", LastValidatorPower],
    ["/zrchain.validation.QueryValidatorUnbondingDelegationsRequest", QueryValidatorUnbondingDelegationsRequest],
    ["/zrchain.validation.QueryDelegatorDelegationsResponse", QueryDelegatorDelegationsResponse],
    ["/zrchain.validation.UnbondingDelegationEntry", UnbondingDelegationEntry],
    ["/zrchain.validation.AssetData", AssetData],
    ["/zrchain.validation.MsgDelegate", MsgDelegate],
    ["/zrchain.validation.QueryDelegatorValidatorsRequest", QueryDelegatorValidatorsRequest],
    ["/zrchain.validation.QueryDelegatorValidatorResponse", QueryDelegatorValidatorResponse],
    ["/zrchain.validation.DVPair", DVPair],
    ["/zrchain.validation.MsgCancelUnbondingDelegationResponse", MsgCancelUnbondingDelegationResponse],
    ["/zrchain.validation.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.validation.QueryDelegationResponse", QueryDelegationResponse],
    ["/zrchain.validation.QueryUnbondingDelegationRequest", QueryUnbondingDelegationRequest],
    ["/zrchain.validation.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.validation.StakeAuthorization_Validators", StakeAuthorization_Validators],
    ["/zrchain.validation.MsgCreateValidator", MsgCreateValidator],
    ["/zrchain.validation.SlashEvent", SlashEvent],
    ["/zrchain.validation.MsgEditValidator", MsgEditValidator],
    ["/zrchain.validation.MsgCancelUnbondingDelegation", MsgCancelUnbondingDelegation],
    ["/zrchain.validation.QueryValidatorsResponse", QueryValidatorsResponse],
    ["/zrchain.validation.QueryValidatorRequest", QueryValidatorRequest],
    ["/zrchain.validation.QueryRedelegationsResponse", QueryRedelegationsResponse],
    ["/zrchain.validation.QueryPowerRequest", QueryPowerRequest],
    ["/zrchain.validation.ValidationInfo", ValidationInfo],
    ["/zrchain.validation.MsgCreateValidatorResponse", MsgCreateValidatorResponse],
    ["/zrchain.validation.MsgTriggerEventBackfillResponse", MsgTriggerEventBackfillResponse],
    ["/zrchain.validation.QueryDelegationRequest", QueryDelegationRequest],
    ["/zrchain.validation.QueryDelegatorUnbondingDelegationsRequest", QueryDelegatorUnbondingDelegationsRequest],
    ["/zrchain.validation.QueryBackfillRequestsRequest", QueryBackfillRequestsRequest],
    ["/zrchain.validation.Description", Description],
    ["/zrchain.validation.Params", Params],
    ["/zrchain.validation.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.validation.BackfillRequests", BackfillRequests],
    ["/zrchain.validation.QueryDelegatorValidatorRequest", QueryDelegatorValidatorRequest],
    ["/zrchain.validation.DVPairs", DVPairs],
    ["/zrchain.validation.DVVTriplet", DVVTriplet],
    ["/zrchain.validation.RedelegationEntry", RedelegationEntry],
    ["/zrchain.validation.MsgBeginRedelegateResponse", MsgBeginRedelegateResponse],
    ["/zrchain.validation.MsgUpdateHVParams", MsgUpdateHVParams],
    ["/zrchain.validation.Commission", Commission],
    ["/zrchain.validation.QueryDelegatorDelegationsRequest", QueryDelegatorDelegationsRequest],
    ["/zrchain.validation.DelegationResponse", DelegationResponse],
    ["/zrchain.validation.RedelegationResponse", RedelegationResponse],
    ["/zrchain.validation.MsgDelegateResponse", MsgDelegateResponse],
    ["/zrchain.validation.MsgUndelegate", MsgUndelegate],
    ["/zrchain.validation.QueryValidatorDelegationsRequest", QueryValidatorDelegationsRequest],
    ["/zrchain.validation.ValidatorPower", ValidatorPower],
    ["/zrchain.validation.SolanaNonce", SolanaNonce],
    ["/zrchain.validation.StakeAuthorization", StakeAuthorization],
    ["/zrchain.validation.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.validation.QueryBackfillRequestsResponse", QueryBackfillRequestsResponse],
    ["/zrchain.validation.Delegation", Delegation],
    ["/zrchain.validation.UnbondingDelegation", UnbondingDelegation],
    ["/zrchain.validation.Redelegation", Redelegation],
    ["/zrchain.validation.ValidatorUpdates", ValidatorUpdates],
    ["/zrchain.validation.HistoricalInfoHV", HistoricalInfoHV],
    ["/zrchain.validation.MsgUndelegateResponse", MsgUndelegateResponse],
    ["/zrchain.validation.QueryValidatorDelegationsResponse", QueryValidatorDelegationsResponse],
    ["/zrchain.validation.QueryRedelegationsRequest", QueryRedelegationsRequest],
    ["/zrchain.validation.MsgEditValidatorResponse", MsgEditValidatorResponse],
    ["/zrchain.validation.MsgBeginRedelegate", MsgBeginRedelegate],
    ["/zrchain.validation.MsgUpdateHVParamsResponse", MsgUpdateHVParamsResponse],
    ["/zrchain.validation.MsgTriggerEventBackfill", MsgTriggerEventBackfill],
    ["/zrchain.validation.DVVTriplets", DVVTriplets],
    
];

export { msgTypes }