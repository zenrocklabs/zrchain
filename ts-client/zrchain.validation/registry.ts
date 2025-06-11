import { GeneratedType } from "@cosmjs/proto-signing";
import { HistoricalInfoHV } from "./types/zrchain/validation/hybrid_validation";
import { ValAddresses } from "./types/zrchain/validation/staking";
import { Redelegation } from "./types/zrchain/validation/staking";
import { MsgCancelUnbondingDelegation } from "./types/zrchain/validation/tx";
import { QueryDelegationRequest } from "./types/zrchain/validation/query";
import { SlashEvent } from "./types/zrchain/validation/hybrid_validation";
import { MsgBeginRedelegate } from "./types/zrchain/validation/tx";
import { MsgBeginRedelegateResponse } from "./types/zrchain/validation/tx";
import { QueryUnbondingDelegationRequest } from "./types/zrchain/validation/query";
import { QueryUnbondingDelegationResponse } from "./types/zrchain/validation/query";
import { QueryDelegatorDelegationsRequest } from "./types/zrchain/validation/query";
import { MsgCreateValidator } from "./types/zrchain/validation/tx";
import { SolanaNonce } from "./types/zrchain/validation/solana";
import { UnbondingDelegationEntry } from "./types/zrchain/validation/staking";
import { MsgEditValidator } from "./types/zrchain/validation/tx";
import { QueryPowerResponse } from "./types/zrchain/validation/query";
import { ValidatorHV } from "./types/zrchain/validation/hybrid_validation";
import { RedelegationEntry } from "./types/zrchain/validation/staking";
import { RedelegationEntryResponse } from "./types/zrchain/validation/staking";
import { QueryRedelegationsResponse } from "./types/zrchain/validation/query";
import { QueryPoolRequest } from "./types/zrchain/validation/query";
import { LastValidatorPower } from "./types/zrchain/validation/genesis";
import { DelegationResponse } from "./types/zrchain/validation/staking";
import { MsgUndelegate } from "./types/zrchain/validation/tx";
import { QueryDelegatorValidatorsRequest } from "./types/zrchain/validation/query";
import { CommissionRates } from "./types/zrchain/validation/staking";
import { StakeAuthorization } from "./types/zrchain/validation/authz";
import { MsgEditValidatorResponse } from "./types/zrchain/validation/tx";
import { QueryValidatorDelegationsRequest } from "./types/zrchain/validation/query";
import { QueryValidatorUnbondingDelegationsResponse } from "./types/zrchain/validation/query";
import { DVPair } from "./types/zrchain/validation/staking";
import { QueryValidatorResponse } from "./types/zrchain/validation/query";
import { MsgCreateValidatorResponse } from "./types/zrchain/validation/tx";
import { QueryDelegatorValidatorResponse } from "./types/zrchain/validation/query";
import { ValidatorPower } from "./types/zrchain/validation/query";
import { MsgCancelUnbondingDelegationResponse } from "./types/zrchain/validation/tx";
import { UnbondingDelegation } from "./types/zrchain/validation/staking";
import { Params } from "./types/zrchain/validation/staking";
import { HVParams } from "./types/zrchain/validation/hybrid_validation";
import { DVVTriplet } from "./types/zrchain/validation/staking";
import { QueryParamsResponse } from "./types/zrchain/validation/query";
import { DVPairs } from "./types/zrchain/validation/staking";
import { StakeAuthorization_Validators } from "./types/zrchain/validation/authz";
import { QueryValidatorDelegationsResponse } from "./types/zrchain/validation/query";
import { QueryPowerRequest } from "./types/zrchain/validation/query";
import { Delegation } from "./types/zrchain/validation/staking";
import { MsgUpdateHVParams } from "./types/zrchain/validation/tx";
import { QueryDelegatorDelegationsResponse } from "./types/zrchain/validation/query";
import { QueryDelegatorUnbondingDelegationsRequest } from "./types/zrchain/validation/query";
import { QueryDelegatorUnbondingDelegationsResponse } from "./types/zrchain/validation/query";
import { QueryDelegatorValidatorsResponse } from "./types/zrchain/validation/query";
import { Description } from "./types/zrchain/validation/staking";
import { GenesisState } from "./types/zrchain/validation/genesis";
import { QueryDelegationResponse } from "./types/zrchain/validation/query";
import { QueryHistoricalInfoResponse } from "./types/zrchain/validation/query";
import { DVVTriplets } from "./types/zrchain/validation/staking";
import { AssetData } from "./types/zrchain/validation/asset_data";
import { RedelegationResponse } from "./types/zrchain/validation/staking";
import { Pool } from "./types/zrchain/validation/staking";
import { MsgDelegateResponse } from "./types/zrchain/validation/tx";
import { MsgUpdateParams } from "./types/zrchain/validation/tx";
import { QueryRedelegationsRequest } from "./types/zrchain/validation/query";
import { ValidatorUpdates } from "./types/zrchain/validation/staking";
import { MsgDelegate } from "./types/zrchain/validation/tx";
import { MsgUpdateParamsResponse } from "./types/zrchain/validation/tx";
import { QueryValidatorsResponse } from "./types/zrchain/validation/query";
import { QueryDelegatorValidatorRequest } from "./types/zrchain/validation/query";
import { QueryValidatorsRequest } from "./types/zrchain/validation/query";
import { QueryValidatorRequest } from "./types/zrchain/validation/query";
import { QueryValidatorUnbondingDelegationsRequest } from "./types/zrchain/validation/query";
import { QueryPoolResponse } from "./types/zrchain/validation/query";
import { Commission } from "./types/zrchain/validation/staking";
import { ValidationInfo } from "./types/zrchain/validation/hybrid_validation";
import { MsgUndelegateResponse } from "./types/zrchain/validation/tx";
import { QueryHistoricalInfoRequest } from "./types/zrchain/validation/query";
import { QueryParamsRequest } from "./types/zrchain/validation/query";
import { MsgUpdateHVParamsResponse } from "./types/zrchain/validation/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.validation.HistoricalInfoHV", HistoricalInfoHV],
    ["/zrchain.validation.ValAddresses", ValAddresses],
    ["/zrchain.validation.Redelegation", Redelegation],
    ["/zrchain.validation.MsgCancelUnbondingDelegation", MsgCancelUnbondingDelegation],
    ["/zrchain.validation.QueryDelegationRequest", QueryDelegationRequest],
    ["/zrchain.validation.SlashEvent", SlashEvent],
    ["/zrchain.validation.MsgBeginRedelegate", MsgBeginRedelegate],
    ["/zrchain.validation.MsgBeginRedelegateResponse", MsgBeginRedelegateResponse],
    ["/zrchain.validation.QueryUnbondingDelegationRequest", QueryUnbondingDelegationRequest],
    ["/zrchain.validation.QueryUnbondingDelegationResponse", QueryUnbondingDelegationResponse],
    ["/zrchain.validation.QueryDelegatorDelegationsRequest", QueryDelegatorDelegationsRequest],
    ["/zrchain.validation.MsgCreateValidator", MsgCreateValidator],
    ["/zrchain.validation.SolanaNonce", SolanaNonce],
    ["/zrchain.validation.UnbondingDelegationEntry", UnbondingDelegationEntry],
    ["/zrchain.validation.MsgEditValidator", MsgEditValidator],
    ["/zrchain.validation.QueryPowerResponse", QueryPowerResponse],
    ["/zrchain.validation.ValidatorHV", ValidatorHV],
    ["/zrchain.validation.RedelegationEntry", RedelegationEntry],
    ["/zrchain.validation.RedelegationEntryResponse", RedelegationEntryResponse],
    ["/zrchain.validation.QueryRedelegationsResponse", QueryRedelegationsResponse],
    ["/zrchain.validation.QueryPoolRequest", QueryPoolRequest],
    ["/zrchain.validation.LastValidatorPower", LastValidatorPower],
    ["/zrchain.validation.DelegationResponse", DelegationResponse],
    ["/zrchain.validation.MsgUndelegate", MsgUndelegate],
    ["/zrchain.validation.QueryDelegatorValidatorsRequest", QueryDelegatorValidatorsRequest],
    ["/zrchain.validation.CommissionRates", CommissionRates],
    ["/zrchain.validation.StakeAuthorization", StakeAuthorization],
    ["/zrchain.validation.MsgEditValidatorResponse", MsgEditValidatorResponse],
    ["/zrchain.validation.QueryValidatorDelegationsRequest", QueryValidatorDelegationsRequest],
    ["/zrchain.validation.QueryValidatorUnbondingDelegationsResponse", QueryValidatorUnbondingDelegationsResponse],
    ["/zrchain.validation.DVPair", DVPair],
    ["/zrchain.validation.QueryValidatorResponse", QueryValidatorResponse],
    ["/zrchain.validation.MsgCreateValidatorResponse", MsgCreateValidatorResponse],
    ["/zrchain.validation.QueryDelegatorValidatorResponse", QueryDelegatorValidatorResponse],
    ["/zrchain.validation.ValidatorPower", ValidatorPower],
    ["/zrchain.validation.MsgCancelUnbondingDelegationResponse", MsgCancelUnbondingDelegationResponse],
    ["/zrchain.validation.UnbondingDelegation", UnbondingDelegation],
    ["/zrchain.validation.Params", Params],
    ["/zrchain.validation.HVParams", HVParams],
    ["/zrchain.validation.DVVTriplet", DVVTriplet],
    ["/zrchain.validation.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.validation.DVPairs", DVPairs],
    ["/zrchain.validation.StakeAuthorization_Validators", StakeAuthorization_Validators],
    ["/zrchain.validation.QueryValidatorDelegationsResponse", QueryValidatorDelegationsResponse],
    ["/zrchain.validation.QueryPowerRequest", QueryPowerRequest],
    ["/zrchain.validation.Delegation", Delegation],
    ["/zrchain.validation.MsgUpdateHVParams", MsgUpdateHVParams],
    ["/zrchain.validation.QueryDelegatorDelegationsResponse", QueryDelegatorDelegationsResponse],
    ["/zrchain.validation.QueryDelegatorUnbondingDelegationsRequest", QueryDelegatorUnbondingDelegationsRequest],
    ["/zrchain.validation.QueryDelegatorUnbondingDelegationsResponse", QueryDelegatorUnbondingDelegationsResponse],
    ["/zrchain.validation.QueryDelegatorValidatorsResponse", QueryDelegatorValidatorsResponse],
    ["/zrchain.validation.Description", Description],
    ["/zrchain.validation.GenesisState", GenesisState],
    ["/zrchain.validation.QueryDelegationResponse", QueryDelegationResponse],
    ["/zrchain.validation.QueryHistoricalInfoResponse", QueryHistoricalInfoResponse],
    ["/zrchain.validation.DVVTriplets", DVVTriplets],
    ["/zrchain.validation.AssetData", AssetData],
    ["/zrchain.validation.RedelegationResponse", RedelegationResponse],
    ["/zrchain.validation.Pool", Pool],
    ["/zrchain.validation.MsgDelegateResponse", MsgDelegateResponse],
    ["/zrchain.validation.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.validation.QueryRedelegationsRequest", QueryRedelegationsRequest],
    ["/zrchain.validation.ValidatorUpdates", ValidatorUpdates],
    ["/zrchain.validation.MsgDelegate", MsgDelegate],
    ["/zrchain.validation.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.validation.QueryValidatorsResponse", QueryValidatorsResponse],
    ["/zrchain.validation.QueryDelegatorValidatorRequest", QueryDelegatorValidatorRequest],
    ["/zrchain.validation.QueryValidatorsRequest", QueryValidatorsRequest],
    ["/zrchain.validation.QueryValidatorRequest", QueryValidatorRequest],
    ["/zrchain.validation.QueryValidatorUnbondingDelegationsRequest", QueryValidatorUnbondingDelegationsRequest],
    ["/zrchain.validation.QueryPoolResponse", QueryPoolResponse],
    ["/zrchain.validation.Commission", Commission],
    ["/zrchain.validation.ValidationInfo", ValidationInfo],
    ["/zrchain.validation.MsgUndelegateResponse", MsgUndelegateResponse],
    ["/zrchain.validation.QueryHistoricalInfoRequest", QueryHistoricalInfoRequest],
    ["/zrchain.validation.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.validation.MsgUpdateHVParamsResponse", MsgUpdateHVParamsResponse],
    
];

export { msgTypes }