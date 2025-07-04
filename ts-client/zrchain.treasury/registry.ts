import { GeneratedType } from "@cosmjs/proto-signing";
import { MetadataSolana } from "./types/zrchain/treasury/tx";
import { MsgUpdateParams } from "./types/zrchain/treasury/tx";
import { QueryKeyRequestByIDRequest } from "./types/zrchain/treasury/query";
import { QueryKeyRequestByIDResponse } from "./types/zrchain/treasury/query";
import { ICATransactionRequest } from "./types/zrchain/treasury/mpcsign";
import { Params } from "./types/zrchain/treasury/params";
import { MsgFulfilSignatureRequestResponse } from "./types/zrchain/treasury/tx";
import { MetadataEthereum } from "./types/zrchain/treasury/tx";
import { MsgNewICATransactionRequest } from "./types/zrchain/treasury/tx";
import { KeyAndWalletResponse } from "./types/zrchain/treasury/query";
import { QuerySignatureRequestByIDResponse } from "./types/zrchain/treasury/query";
import { SignedDataWithID } from "./types/zrchain/treasury/mpcsign";
import { MsgNewSignatureRequestResponse } from "./types/zrchain/treasury/tx";
import { MsgNewSignTransactionRequestResponse } from "./types/zrchain/treasury/tx";
import { MsgTransferFromKeyring } from "./types/zrchain/treasury/tx";
import { QuerySignTransactionRequestByIDResponse } from "./types/zrchain/treasury/query";
import { KeyRequest } from "./types/zrchain/treasury/key";
import { PartySignature } from "./types/zrchain/treasury/key";
import { SignRequest } from "./types/zrchain/treasury/mpcsign";
import { MsgFulfilKeyRequest } from "./types/zrchain/treasury/tx";
import { MsgNewSignatureRequest } from "./types/zrchain/treasury/tx";
import { QueryKeyRequestsResponse } from "./types/zrchain/treasury/query";
import { QueryKeyByIDRequest } from "./types/zrchain/treasury/query";
import { QueryKeyByIDResponse } from "./types/zrchain/treasury/query";
import { QueryZenbtcWalletsRequest } from "./types/zrchain/treasury/query";
import { MsgUpdateParamsResponse } from "./types/zrchain/treasury/tx";
import { QuerySignatureRequestsResponse } from "./types/zrchain/treasury/query";
import { QueryZrSignKeysResponse } from "./types/zrchain/treasury/query";
import { GenesisState } from "./types/zrchain/treasury/genesis";
import { Key } from "./types/zrchain/treasury/key";
import { KeyReqResponse } from "./types/zrchain/treasury/key";
import { KeyResponse } from "./types/zrchain/treasury/key";
import { NoData } from "./types/zrchain/treasury/packet";
import { SignReqResponse } from "./types/zrchain/treasury/mpcsign";
import { MsgFulfilSignatureRequest } from "./types/zrchain/treasury/tx";
import { MsgNewZrSignSignatureRequest } from "./types/zrchain/treasury/tx";
import { QueryKeyRequestsRequest } from "./types/zrchain/treasury/query";
import { QuerySignatureRequestsRequest } from "./types/zrchain/treasury/query";
import { QuerySignatureRequestByIDRequest } from "./types/zrchain/treasury/query";
import { QuerySignTransactionRequestsResponse } from "./types/zrchain/treasury/query";
import { MsgNewKeyRequest } from "./types/zrchain/treasury/tx";
import { MsgFulfilKeyRequestResponse } from "./types/zrchain/treasury/tx";
import { MsgTransferFromKeyringResponse } from "./types/zrchain/treasury/tx";
import { QueryKeysRequest } from "./types/zrchain/treasury/query";
import { SignTransactionRequestsResponse } from "./types/zrchain/treasury/query";
import { QueryZenbtcWalletsResponse } from "./types/zrchain/treasury/query";
import { MsgNewICATransactionRequestResponse } from "./types/zrchain/treasury/tx";
import { QueryKeysResponse } from "./types/zrchain/treasury/query";
import { QuerySignTransactionRequestsRequest } from "./types/zrchain/treasury/query";
import { QueryKeyByAddressRequest } from "./types/zrchain/treasury/query";
import { QueryKeyByAddressResponse } from "./types/zrchain/treasury/query";
import { MsgNewKey } from "./types/zrchain/treasury/tx";
import { MsgUpdateKeyPolicy } from "./types/zrchain/treasury/tx";
import { QueryParamsResponse } from "./types/zrchain/treasury/query";
import { QuerySignTransactionRequestByIDRequest } from "./types/zrchain/treasury/query";
import { QueryZrSignKeysRequest } from "./types/zrchain/treasury/query";
import { MsgNewKeyRequestResponse } from "./types/zrchain/treasury/tx";
import { MsgFulfilICATransactionRequest } from "./types/zrchain/treasury/tx";
import { SignTransactionRequest } from "./types/zrchain/treasury/mpcsign";
import { SignTxReqResponse } from "./types/zrchain/treasury/mpcsign";
import { WalletResponse } from "./types/zrchain/treasury/query";
import { TreasuryPacketData } from "./types/zrchain/treasury/packet";
import { MsgNewSignTransactionRequest } from "./types/zrchain/treasury/tx";
import { MsgFulfilICATransactionRequestResponse } from "./types/zrchain/treasury/tx";
import { MsgUpdateKeyPolicyResponse } from "./types/zrchain/treasury/tx";
import { QueryParamsRequest } from "./types/zrchain/treasury/query";
import { ZrSignKeyEntry } from "./types/zrchain/treasury/query";
import { ZenBTCMetadata } from "./types/zrchain/treasury/key";
import { MsgNewZrSignSignatureRequestResponse } from "./types/zrchain/treasury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.treasury.MetadataSolana", MetadataSolana],
    ["/zrchain.treasury.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.treasury.QueryKeyRequestByIDRequest", QueryKeyRequestByIDRequest],
    ["/zrchain.treasury.QueryKeyRequestByIDResponse", QueryKeyRequestByIDResponse],
    ["/zrchain.treasury.ICATransactionRequest", ICATransactionRequest],
    ["/zrchain.treasury.Params", Params],
    ["/zrchain.treasury.MsgFulfilSignatureRequestResponse", MsgFulfilSignatureRequestResponse],
    ["/zrchain.treasury.MetadataEthereum", MetadataEthereum],
    ["/zrchain.treasury.MsgNewICATransactionRequest", MsgNewICATransactionRequest],
    ["/zrchain.treasury.KeyAndWalletResponse", KeyAndWalletResponse],
    ["/zrchain.treasury.QuerySignatureRequestByIDResponse", QuerySignatureRequestByIDResponse],
    ["/zrchain.treasury.SignedDataWithID", SignedDataWithID],
    ["/zrchain.treasury.MsgNewSignatureRequestResponse", MsgNewSignatureRequestResponse],
    ["/zrchain.treasury.MsgNewSignTransactionRequestResponse", MsgNewSignTransactionRequestResponse],
    ["/zrchain.treasury.MsgTransferFromKeyring", MsgTransferFromKeyring],
    ["/zrchain.treasury.QuerySignTransactionRequestByIDResponse", QuerySignTransactionRequestByIDResponse],
    ["/zrchain.treasury.KeyRequest", KeyRequest],
    ["/zrchain.treasury.PartySignature", PartySignature],
    ["/zrchain.treasury.SignRequest", SignRequest],
    ["/zrchain.treasury.MsgFulfilKeyRequest", MsgFulfilKeyRequest],
    ["/zrchain.treasury.MsgNewSignatureRequest", MsgNewSignatureRequest],
    ["/zrchain.treasury.QueryKeyRequestsResponse", QueryKeyRequestsResponse],
    ["/zrchain.treasury.QueryKeyByIDRequest", QueryKeyByIDRequest],
    ["/zrchain.treasury.QueryKeyByIDResponse", QueryKeyByIDResponse],
    ["/zrchain.treasury.QueryZenbtcWalletsRequest", QueryZenbtcWalletsRequest],
    ["/zrchain.treasury.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.treasury.QuerySignatureRequestsResponse", QuerySignatureRequestsResponse],
    ["/zrchain.treasury.QueryZrSignKeysResponse", QueryZrSignKeysResponse],
    ["/zrchain.treasury.GenesisState", GenesisState],
    ["/zrchain.treasury.Key", Key],
    ["/zrchain.treasury.KeyReqResponse", KeyReqResponse],
    ["/zrchain.treasury.KeyResponse", KeyResponse],
    ["/zrchain.treasury.NoData", NoData],
    ["/zrchain.treasury.SignReqResponse", SignReqResponse],
    ["/zrchain.treasury.MsgFulfilSignatureRequest", MsgFulfilSignatureRequest],
    ["/zrchain.treasury.MsgNewZrSignSignatureRequest", MsgNewZrSignSignatureRequest],
    ["/zrchain.treasury.QueryKeyRequestsRequest", QueryKeyRequestsRequest],
    ["/zrchain.treasury.QuerySignatureRequestsRequest", QuerySignatureRequestsRequest],
    ["/zrchain.treasury.QuerySignatureRequestByIDRequest", QuerySignatureRequestByIDRequest],
    ["/zrchain.treasury.QuerySignTransactionRequestsResponse", QuerySignTransactionRequestsResponse],
    ["/zrchain.treasury.MsgNewKeyRequest", MsgNewKeyRequest],
    ["/zrchain.treasury.MsgFulfilKeyRequestResponse", MsgFulfilKeyRequestResponse],
    ["/zrchain.treasury.MsgTransferFromKeyringResponse", MsgTransferFromKeyringResponse],
    ["/zrchain.treasury.QueryKeysRequest", QueryKeysRequest],
    ["/zrchain.treasury.SignTransactionRequestsResponse", SignTransactionRequestsResponse],
    ["/zrchain.treasury.QueryZenbtcWalletsResponse", QueryZenbtcWalletsResponse],
    ["/zrchain.treasury.MsgNewICATransactionRequestResponse", MsgNewICATransactionRequestResponse],
    ["/zrchain.treasury.QueryKeysResponse", QueryKeysResponse],
    ["/zrchain.treasury.QuerySignTransactionRequestsRequest", QuerySignTransactionRequestsRequest],
    ["/zrchain.treasury.QueryKeyByAddressRequest", QueryKeyByAddressRequest],
    ["/zrchain.treasury.QueryKeyByAddressResponse", QueryKeyByAddressResponse],
    ["/zrchain.treasury.MsgNewKey", MsgNewKey],
    ["/zrchain.treasury.MsgUpdateKeyPolicy", MsgUpdateKeyPolicy],
    ["/zrchain.treasury.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.treasury.QuerySignTransactionRequestByIDRequest", QuerySignTransactionRequestByIDRequest],
    ["/zrchain.treasury.QueryZrSignKeysRequest", QueryZrSignKeysRequest],
    ["/zrchain.treasury.MsgNewKeyRequestResponse", MsgNewKeyRequestResponse],
    ["/zrchain.treasury.MsgFulfilICATransactionRequest", MsgFulfilICATransactionRequest],
    ["/zrchain.treasury.SignTransactionRequest", SignTransactionRequest],
    ["/zrchain.treasury.SignTxReqResponse", SignTxReqResponse],
    ["/zrchain.treasury.WalletResponse", WalletResponse],
    ["/zrchain.treasury.TreasuryPacketData", TreasuryPacketData],
    ["/zrchain.treasury.MsgNewSignTransactionRequest", MsgNewSignTransactionRequest],
    ["/zrchain.treasury.MsgFulfilICATransactionRequestResponse", MsgFulfilICATransactionRequestResponse],
    ["/zrchain.treasury.MsgUpdateKeyPolicyResponse", MsgUpdateKeyPolicyResponse],
    ["/zrchain.treasury.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.treasury.ZrSignKeyEntry", ZrSignKeyEntry],
    ["/zrchain.treasury.ZenBTCMetadata", ZenBTCMetadata],
    ["/zrchain.treasury.MsgNewZrSignSignatureRequestResponse", MsgNewZrSignSignatureRequestResponse],
    
];

export { msgTypes }