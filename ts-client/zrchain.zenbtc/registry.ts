import { GeneratedType } from "@cosmjs/proto-signing";
import { QueryParamsRequest } from "./types/zrchain/zenbtc/query";
import { QuerySupplyResponse } from "./types/zrchain/zenbtc/query";
import { NonceData } from "./types/zrchain/zenbtc/mint";
import { PendingMintTransactions } from "./types/zrchain/zenbtc/mint";
import { QueryParamsResponse } from "./types/zrchain/zenbtc/query";
import { QueryLockTransactionsRequest } from "./types/zrchain/zenbtc/query";
import { QueryRedemptionsResponse } from "./types/zrchain/zenbtc/query";
import { BurnEvent } from "./types/zrchain/zenbtc/redemptions";
import { MsgSubmitUnsignedRedemptionTxResponse } from "./types/zrchain/zenbtc/tx";
import { QueryRedemptionsRequest } from "./types/zrchain/zenbtc/query";
import { QueryPendingMintTransactionsRequest } from "./types/zrchain/zenbtc/query";
import { RedemptionData } from "./types/zrchain/zenbtc/redemptions";
import { GenesisState } from "./types/zrchain/zenbtc/genesis";
import { RequestedBitcoinHeaders } from "./types/zrchain/zenbtc/mint";
import { QuerySupplyRequest } from "./types/zrchain/zenbtc/query";
import { Supply } from "./types/zrchain/zenbtc/supply";
import { Solana } from "./types/zrchain/zenbtc/params";
import { MsgUpdateParamsResponse } from "./types/zrchain/zenbtc/tx";
import { MsgVerifyDepositBlockInclusionResponse } from "./types/zrchain/zenbtc/tx";
import { PendingMintTransaction } from "./types/zrchain/zenbtc/mint";
import { QueryLockTransactionsResponse } from "./types/zrchain/zenbtc/query";
import { QueryPendingMintTransactionsResponse } from "./types/zrchain/zenbtc/query";
import { Params } from "./types/zrchain/zenbtc/params";
import { Redemption } from "./types/zrchain/zenbtc/redemptions";
import { MsgSubmitUnsignedRedemptionTx } from "./types/zrchain/zenbtc/tx";
import { MsgUpdateParams } from "./types/zrchain/zenbtc/tx";
import { MsgVerifyDepositBlockInclusion } from "./types/zrchain/zenbtc/tx";
import { InputHashes } from "./types/zrchain/zenbtc/tx";
import { LockTransaction } from "./types/zrchain/zenbtc/mint";
import { QueryBurnEventsRequest } from "./types/zrchain/zenbtc/query";
import { QueryBurnEventsResponse } from "./types/zrchain/zenbtc/query";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.zenbtc.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.zenbtc.QuerySupplyResponse", QuerySupplyResponse],
    ["/zrchain.zenbtc.NonceData", NonceData],
    ["/zrchain.zenbtc.PendingMintTransactions", PendingMintTransactions],
    ["/zrchain.zenbtc.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.zenbtc.QueryLockTransactionsRequest", QueryLockTransactionsRequest],
    ["/zrchain.zenbtc.QueryRedemptionsResponse", QueryRedemptionsResponse],
    ["/zrchain.zenbtc.BurnEvent", BurnEvent],
    ["/zrchain.zenbtc.MsgSubmitUnsignedRedemptionTxResponse", MsgSubmitUnsignedRedemptionTxResponse],
    ["/zrchain.zenbtc.QueryRedemptionsRequest", QueryRedemptionsRequest],
    ["/zrchain.zenbtc.QueryPendingMintTransactionsRequest", QueryPendingMintTransactionsRequest],
    ["/zrchain.zenbtc.RedemptionData", RedemptionData],
    ["/zrchain.zenbtc.GenesisState", GenesisState],
    ["/zrchain.zenbtc.RequestedBitcoinHeaders", RequestedBitcoinHeaders],
    ["/zrchain.zenbtc.QuerySupplyRequest", QuerySupplyRequest],
    ["/zrchain.zenbtc.Supply", Supply],
    ["/zrchain.zenbtc.Solana", Solana],
    ["/zrchain.zenbtc.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.zenbtc.MsgVerifyDepositBlockInclusionResponse", MsgVerifyDepositBlockInclusionResponse],
    ["/zrchain.zenbtc.PendingMintTransaction", PendingMintTransaction],
    ["/zrchain.zenbtc.QueryLockTransactionsResponse", QueryLockTransactionsResponse],
    ["/zrchain.zenbtc.QueryPendingMintTransactionsResponse", QueryPendingMintTransactionsResponse],
    ["/zrchain.zenbtc.Params", Params],
    ["/zrchain.zenbtc.Redemption", Redemption],
    ["/zrchain.zenbtc.MsgSubmitUnsignedRedemptionTx", MsgSubmitUnsignedRedemptionTx],
    ["/zrchain.zenbtc.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.zenbtc.MsgVerifyDepositBlockInclusion", MsgVerifyDepositBlockInclusion],
    ["/zrchain.zenbtc.InputHashes", InputHashes],
    ["/zrchain.zenbtc.LockTransaction", LockTransaction],
    ["/zrchain.zenbtc.QueryBurnEventsRequest", QueryBurnEventsRequest],
    ["/zrchain.zenbtc.QueryBurnEventsResponse", QueryBurnEventsResponse],
    
];

export { msgTypes }