import { GeneratedType } from "@cosmjs/proto-signing";
import { GenesisState } from "./types/zrchain/zentp/genesis";
import { Params } from "./types/zrchain/zentp/params";
import { QueryParamsRequest } from "./types/zrchain/zentp/query";
import { QueryParamsResponse } from "./types/zrchain/zentp/query";
import { QueryMintsResponse } from "./types/zrchain/zentp/query";
import { QueryStatsRequest } from "./types/zrchain/zentp/query";
import { QueryStatsResponse } from "./types/zrchain/zentp/query";
import { Bridge } from "./types/zrchain/zentp/bridge";
import { QueryBurnsRequest } from "./types/zrchain/zentp/query";
import { QuerySolanaROCKSupplyRequest } from "./types/zrchain/zentp/query";
import { QueryBurnsResponse } from "./types/zrchain/zentp/query";
import { QuerySolanaROCKSupplyResponse } from "./types/zrchain/zentp/query";
import { MsgUpdateParamsResponse } from "./types/zrchain/zentp/tx";
import { MsgBridgeResponse } from "./types/zrchain/zentp/tx";
import { MsgBurnResponse } from "./types/zrchain/zentp/tx";
import { MsgSetSolanaROCKSupplyResponse } from "./types/zrchain/zentp/tx";
import { QueryMintsRequest } from "./types/zrchain/zentp/query";
import { MsgUpdateParams } from "./types/zrchain/zentp/tx";
import { MsgBridge } from "./types/zrchain/zentp/tx";
import { MsgBurn } from "./types/zrchain/zentp/tx";
import { MsgSetSolanaROCKSupply } from "./types/zrchain/zentp/tx";
import { Solana } from "./types/zrchain/zentp/params";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.zentp.GenesisState", GenesisState],
    ["/zrchain.zentp.Params", Params],
    ["/zrchain.zentp.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.zentp.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.zentp.QueryMintsResponse", QueryMintsResponse],
    ["/zrchain.zentp.QueryStatsRequest", QueryStatsRequest],
    ["/zrchain.zentp.QueryStatsResponse", QueryStatsResponse],
    ["/zrchain.zentp.Bridge", Bridge],
    ["/zrchain.zentp.QueryBurnsRequest", QueryBurnsRequest],
    ["/zrchain.zentp.QuerySolanaROCKSupplyRequest", QuerySolanaROCKSupplyRequest],
    ["/zrchain.zentp.QueryBurnsResponse", QueryBurnsResponse],
    ["/zrchain.zentp.QuerySolanaROCKSupplyResponse", QuerySolanaROCKSupplyResponse],
    ["/zrchain.zentp.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.zentp.MsgBridgeResponse", MsgBridgeResponse],
    ["/zrchain.zentp.MsgBurnResponse", MsgBurnResponse],
    ["/zrchain.zentp.MsgSetSolanaROCKSupplyResponse", MsgSetSolanaROCKSupplyResponse],
    ["/zrchain.zentp.QueryMintsRequest", QueryMintsRequest],
    ["/zrchain.zentp.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.zentp.MsgBridge", MsgBridge],
    ["/zrchain.zentp.MsgBurn", MsgBurn],
    ["/zrchain.zentp.MsgSetSolanaROCKSupply", MsgSetSolanaROCKSupply],
    ["/zrchain.zentp.Solana", Solana],
    
];

export { msgTypes }