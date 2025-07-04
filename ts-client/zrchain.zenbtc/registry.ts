import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgVerifyDepositBlockInclusion } from "./types/zrchain/zenbtc/tx";
import { QueryRedemptionsRequest } from "./types/zrchain/zenbtc/query";
import { QuerySupplyRequest } from "./types/zrchain/zenbtc/query";
import { Redemption } from "./types/zrchain/zenbtc/redemptions";
import { MsgUpdateParamsResponse } from "./types/zrchain/zenbtc/tx";
import { InputHashes } from "./types/zrchain/zenbtc/tx";
import { QueryLockTransactionsResponse } from "./types/zrchain/zenbtc/query";
import { QueryPendingMintTransactionsResponse } from "./types/zrchain/zenbtc/query";
import { QueryPendingMintTransactionResponse } from "./types/zrchain/zenbtc/query";
import { QueryBurnEventsResponse } from "./types/zrchain/zenbtc/query";
import { GenesisState } from "./types/zrchain/zenbtc/genesis";
import { RedemptionData } from "./types/zrchain/zenbtc/redemptions";
import { MsgVerifyDepositBlockInclusionResponse } from "./types/zrchain/zenbtc/tx";
import { QueryParamsResponse } from "./types/zrchain/zenbtc/query";
import { QueryBurnEventsRequest } from "./types/zrchain/zenbtc/query";
import { NonceData } from "./types/zrchain/zenbtc/mint";
import { RequestedBitcoinHeaders } from "./types/zrchain/zenbtc/mint";
import { LockTransaction } from "./types/zrchain/zenbtc/mint";
import { MsgSubmitUnsignedRedemptionTx } from "./types/zrchain/zenbtc/tx";
import { Solana } from "./types/zrchain/zenbtc/params";
import { MsgSubmitUnsignedRedemptionTxResponse } from "./types/zrchain/zenbtc/tx";
import { QueryPendingMintTransactionsRequest } from "./types/zrchain/zenbtc/query";
import { PendingMintTransaction } from "./types/zrchain/zenbtc/mint";
import { PendingMintTransactions } from "./types/zrchain/zenbtc/mint";
import { Params } from "./types/zrchain/zenbtc/params";
import { QueryLockTransactionsRequest } from "./types/zrchain/zenbtc/query";
import { BurnEvent } from "./types/zrchain/zenbtc/redemptions";
import { MsgUpdateParams } from "./types/zrchain/zenbtc/tx";
import { QueryRedemptionsResponse } from "./types/zrchain/zenbtc/query";
import { QuerySupplyResponse } from "./types/zrchain/zenbtc/query";
import { Supply } from "./types/zrchain/zenbtc/supply";
import { QueryParamsRequest } from "./types/zrchain/zenbtc/query";
import { QueryPendingMintTransactionRequest } from "./types/zrchain/zenbtc/query";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.zenbtc.MsgVerifyDepositBlockInclusion", MsgVerifyDepositBlockInclusion],
    ["/zrchain.zenbtc.QueryRedemptionsRequest", QueryRedemptionsRequest],
    ["/zrchain.zenbtc.QuerySupplyRequest", QuerySupplyRequest],
    ["/zrchain.zenbtc.Redemption", Redemption],
    ["/zrchain.zenbtc.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.zenbtc.InputHashes", InputHashes],
    ["/zrchain.zenbtc.QueryLockTransactionsResponse", QueryLockTransactionsResponse],
    ["/zrchain.zenbtc.QueryPendingMintTransactionsResponse", QueryPendingMintTransactionsResponse],
    ["/zrchain.zenbtc.QueryPendingMintTransactionResponse", QueryPendingMintTransactionResponse],
    ["/zrchain.zenbtc.QueryBurnEventsResponse", QueryBurnEventsResponse],
    ["/zrchain.zenbtc.GenesisState", GenesisState],
    ["/zrchain.zenbtc.RedemptionData", RedemptionData],
    ["/zrchain.zenbtc.MsgVerifyDepositBlockInclusionResponse", MsgVerifyDepositBlockInclusionResponse],
    ["/zrchain.zenbtc.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.zenbtc.QueryBurnEventsRequest", QueryBurnEventsRequest],
    ["/zrchain.zenbtc.NonceData", NonceData],
    ["/zrchain.zenbtc.RequestedBitcoinHeaders", RequestedBitcoinHeaders],
    ["/zrchain.zenbtc.LockTransaction", LockTransaction],
    ["/zrchain.zenbtc.MsgSubmitUnsignedRedemptionTx", MsgSubmitUnsignedRedemptionTx],
    ["/zrchain.zenbtc.Solana", Solana],
    ["/zrchain.zenbtc.MsgSubmitUnsignedRedemptionTxResponse", MsgSubmitUnsignedRedemptionTxResponse],
    ["/zrchain.zenbtc.QueryPendingMintTransactionsRequest", QueryPendingMintTransactionsRequest],
    ["/zrchain.zenbtc.PendingMintTransaction", PendingMintTransaction],
    ["/zrchain.zenbtc.PendingMintTransactions", PendingMintTransactions],
    ["/zrchain.zenbtc.Params", Params],
    ["/zrchain.zenbtc.QueryLockTransactionsRequest", QueryLockTransactionsRequest],
    ["/zrchain.zenbtc.BurnEvent", BurnEvent],
    ["/zrchain.zenbtc.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.zenbtc.QueryRedemptionsResponse", QueryRedemptionsResponse],
    ["/zrchain.zenbtc.QuerySupplyResponse", QuerySupplyResponse],
    ["/zrchain.zenbtc.Supply", Supply],
    ["/zrchain.zenbtc.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.zenbtc.QueryPendingMintTransactionRequest", QueryPendingMintTransactionRequest],
    
];

export { msgTypes }