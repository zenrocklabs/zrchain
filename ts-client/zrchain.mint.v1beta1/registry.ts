import { GeneratedType } from "@cosmjs/proto-signing";
import { QueryInflationResponse } from "./types/zrchain/mint/v1beta1/query";
import { QueryAnnualProvisionsRequest } from "./types/zrchain/mint/v1beta1/query";
import { GenesisState } from "./types/zrchain/mint/v1beta1/genesis";
import { Minter } from "./types/zrchain/mint/v1beta1/mint";
import { QueryParamsResponse } from "./types/zrchain/mint/v1beta1/query";
import { QueryInflationRequest } from "./types/zrchain/mint/v1beta1/query";
import { Params } from "./types/zrchain/mint/v1beta1/mint";
import { MsgUpdateParams } from "./types/zrchain/mint/v1beta1/tx";
import { MsgUpdateParamsResponse } from "./types/zrchain/mint/v1beta1/tx";
import { QueryParamsRequest } from "./types/zrchain/mint/v1beta1/query";
import { QueryAnnualProvisionsResponse } from "./types/zrchain/mint/v1beta1/query";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.mint.v1beta1.QueryInflationResponse", QueryInflationResponse],
    ["/zrchain.mint.v1beta1.QueryAnnualProvisionsRequest", QueryAnnualProvisionsRequest],
    ["/zrchain.mint.v1beta1.GenesisState", GenesisState],
    ["/zrchain.mint.v1beta1.Minter", Minter],
    ["/zrchain.mint.v1beta1.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.mint.v1beta1.QueryInflationRequest", QueryInflationRequest],
    ["/zrchain.mint.v1beta1.Params", Params],
    ["/zrchain.mint.v1beta1.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.mint.v1beta1.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.mint.v1beta1.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.mint.v1beta1.QueryAnnualProvisionsResponse", QueryAnnualProvisionsResponse],
    
];

export { msgTypes }