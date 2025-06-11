import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgModuleQuerySafeResponse } from "./types/ibc/applications/interchain_accounts/host/v1/tx";
import { MsgUpdateParams } from "./types/ibc/applications/interchain_accounts/host/v1/tx";
import { QueryRequest } from "./types/ibc/applications/interchain_accounts/host/v1/host";
import { QueryParamsRequest } from "./types/ibc/applications/interchain_accounts/host/v1/query";
import { QueryParamsResponse } from "./types/ibc/applications/interchain_accounts/host/v1/query";
import { Params } from "./types/ibc/applications/interchain_accounts/host/v1/host";
import { MsgUpdateParamsResponse } from "./types/ibc/applications/interchain_accounts/host/v1/tx";
import { MsgModuleQuerySafe } from "./types/ibc/applications/interchain_accounts/host/v1/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/ibc.applications.interchain_accounts.host.v1.MsgModuleQuerySafeResponse", MsgModuleQuerySafeResponse],
    ["/ibc.applications.interchain_accounts.host.v1.MsgUpdateParams", MsgUpdateParams],
    ["/ibc.applications.interchain_accounts.host.v1.QueryRequest", QueryRequest],
    ["/ibc.applications.interchain_accounts.host.v1.QueryParamsRequest", QueryParamsRequest],
    ["/ibc.applications.interchain_accounts.host.v1.QueryParamsResponse", QueryParamsResponse],
    ["/ibc.applications.interchain_accounts.host.v1.Params", Params],
    ["/ibc.applications.interchain_accounts.host.v1.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/ibc.applications.interchain_accounts.host.v1.MsgModuleQuerySafe", MsgModuleQuerySafe],
    
];

export { msgTypes }