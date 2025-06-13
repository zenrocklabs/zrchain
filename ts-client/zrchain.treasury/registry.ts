import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgTransferFromKeyringResponse } from "./types/zrchain/treasury/tx";
import { QueryParamsRequest } from "./types/zrchain/treasury/query";
import { MsgNewICATransactionRequest } from "./types/zrchain/treasury/tx";
import { MsgNewICATransactionRequestResponse } from "./types/zrchain/treasury/tx";
import { SignRequest } from "./types/zrchain/treasury/mpcsign";
import { QuerySignatureRequestByIDResponse } from "./types/zrchain/treasury/query";
import { MsgNewKeyRequest } from "./types/zrchain/treasury/tx";
import { Params } from "./types/zrchain/treasury/params";
import { SignedDataWithID } from "./types/zrchain/treasury/mpcsign";
import { KeyAndWalletResponse } from "./types/zrchain/treasury/query";
import { WalletResponse } from "./types/zrchain/treasury/query";
import { QuerySignTransactionRequestsResponse } from "./types/zrchain/treasury/query";
import { SignTransactionRequestsResponse } from "./types/zrchain/treasury/query";
import { MsgTransferFromKeyring } from "./types/zrchain/treasury/tx";
import { QueryParamsResponse } from "./types/zrchain/treasury/query";
import { QueryKeyByIDResponse } from "./types/zrchain/treasury/query";
import { QuerySignTransactionRequestByIDResponse } from "./types/zrchain/treasury/query";
import { PartySignature } from "./types/zrchain/treasury/key";
import { MsgFulfilKeyRequest } from "./types/zrchain/treasury/tx";
import { MsgNewSignTransactionRequest } from "./types/zrchain/treasury/tx";
import { QueryKeysResponse } from "./types/zrchain/treasury/query";
import { QuerySignTransactionRequestByIDRequest } from "./types/zrchain/treasury/query";
import { MsgFulfilKeyRequestResponse } from "./types/zrchain/treasury/tx";
import { MsgNewSignatureRequest } from "./types/zrchain/treasury/tx";
import { QuerySignatureRequestsResponse } from "./types/zrchain/treasury/query";
import { QueryZrSignKeysRequest } from "./types/zrchain/treasury/query";
import { MsgUpdateParamsResponse } from "./types/zrchain/treasury/tx";
import { MetadataEthereum } from "./types/zrchain/treasury/tx";
import { MsgNewSignTransactionRequestResponse } from "./types/zrchain/treasury/tx";
import { SignReqResponse } from "./types/zrchain/treasury/mpcsign";
import { QueryKeysRequest } from "./types/zrchain/treasury/query";
import { QuerySignTransactionRequestsRequest } from "./types/zrchain/treasury/query";
import { MsgUpdateKeyPolicy } from "./types/zrchain/treasury/tx";
import { QueryZrSignKeysResponse } from "./types/zrchain/treasury/query";
import { QueryKeyByAddressRequest } from "./types/zrchain/treasury/query";
import { TreasuryPacketData } from "./types/zrchain/treasury/packet";
import { MsgNewKey } from "./types/zrchain/treasury/tx";
import { ZenBTCMetadata } from "./types/zrchain/treasury/key";
import { GenesisState } from "./types/zrchain/treasury/genesis";
import { Key } from "./types/zrchain/treasury/key";
import { MsgFulfilICATransactionRequest } from "./types/zrchain/treasury/tx";
import { MsgUpdateKeyPolicyResponse } from "./types/zrchain/treasury/tx";
import { QueryKeyRequestByIDResponse } from "./types/zrchain/treasury/query";
import { SignTxReqResponse } from "./types/zrchain/treasury/mpcsign";
import { QuerySignatureRequestsRequest } from "./types/zrchain/treasury/query";
import { QuerySignatureRequestByIDRequest } from "./types/zrchain/treasury/query";
import { QueryZenbtcWalletsRequest } from "./types/zrchain/treasury/query";
import { KeyReqResponse } from "./types/zrchain/treasury/key";
import { MsgNewSignatureRequestResponse } from "./types/zrchain/treasury/tx";
import { MetadataSolana } from "./types/zrchain/treasury/tx";
import { MsgNewZrSignSignatureRequestResponse } from "./types/zrchain/treasury/tx";
import { QueryKeyRequestsRequest } from "./types/zrchain/treasury/query";
import { ZrSignKeyEntry } from "./types/zrchain/treasury/query";
import { QueryKeyByAddressResponse } from "./types/zrchain/treasury/query";
import { KeyResponse } from "./types/zrchain/treasury/key";
import { NoData } from "./types/zrchain/treasury/packet";
import { MsgNewZrSignSignatureRequest } from "./types/zrchain/treasury/tx";
import { MsgFulfilICATransactionRequestResponse } from "./types/zrchain/treasury/tx";
import { ICATransactionRequest } from "./types/zrchain/treasury/mpcsign";
import { QueryKeyRequestsResponse } from "./types/zrchain/treasury/query";
import { QueryKeyByIDRequest } from "./types/zrchain/treasury/query";
import { QueryZenbtcWalletsResponse } from "./types/zrchain/treasury/query";
import { QueryKeyRequestByIDRequest } from "./types/zrchain/treasury/query";
import { KeyRequest } from "./types/zrchain/treasury/key";
import { MsgFulfilSignatureRequest } from "./types/zrchain/treasury/tx";
import { SignTransactionRequest } from "./types/zrchain/treasury/mpcsign";
import { MsgUpdateParams } from "./types/zrchain/treasury/tx";
import { MsgNewKeyRequestResponse } from "./types/zrchain/treasury/tx";
import { MsgFulfilSignatureRequestResponse } from "./types/zrchain/treasury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.treasury.MsgTransferFromKeyringResponse", MsgTransferFromKeyringResponse],
    ["/zrchain.treasury.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.treasury.MsgNewICATransactionRequest", MsgNewICATransactionRequest],
    ["/zrchain.treasury.MsgNewICATransactionRequestResponse", MsgNewICATransactionRequestResponse],
    ["/zrchain.treasury.SignRequest", SignRequest],
    ["/zrchain.treasury.QuerySignatureRequestByIDResponse", QuerySignatureRequestByIDResponse],
    ["/zrchain.treasury.MsgNewKeyRequest", MsgNewKeyRequest],
    ["/zrchain.treasury.Params", Params],
    ["/zrchain.treasury.SignedDataWithID", SignedDataWithID],
    ["/zrchain.treasury.KeyAndWalletResponse", KeyAndWalletResponse],
    ["/zrchain.treasury.WalletResponse", WalletResponse],
    ["/zrchain.treasury.QuerySignTransactionRequestsResponse", QuerySignTransactionRequestsResponse],
    ["/zrchain.treasury.SignTransactionRequestsResponse", SignTransactionRequestsResponse],
    ["/zrchain.treasury.MsgTransferFromKeyring", MsgTransferFromKeyring],
    ["/zrchain.treasury.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.treasury.QueryKeyByIDResponse", QueryKeyByIDResponse],
    ["/zrchain.treasury.QuerySignTransactionRequestByIDResponse", QuerySignTransactionRequestByIDResponse],
    ["/zrchain.treasury.PartySignature", PartySignature],
    ["/zrchain.treasury.MsgFulfilKeyRequest", MsgFulfilKeyRequest],
    ["/zrchain.treasury.MsgNewSignTransactionRequest", MsgNewSignTransactionRequest],
    ["/zrchain.treasury.QueryKeysResponse", QueryKeysResponse],
    ["/zrchain.treasury.QuerySignTransactionRequestByIDRequest", QuerySignTransactionRequestByIDRequest],
    ["/zrchain.treasury.MsgFulfilKeyRequestResponse", MsgFulfilKeyRequestResponse],
    ["/zrchain.treasury.MsgNewSignatureRequest", MsgNewSignatureRequest],
    ["/zrchain.treasury.QuerySignatureRequestsResponse", QuerySignatureRequestsResponse],
    ["/zrchain.treasury.QueryZrSignKeysRequest", QueryZrSignKeysRequest],
    ["/zrchain.treasury.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.treasury.MetadataEthereum", MetadataEthereum],
    ["/zrchain.treasury.MsgNewSignTransactionRequestResponse", MsgNewSignTransactionRequestResponse],
    ["/zrchain.treasury.SignReqResponse", SignReqResponse],
    ["/zrchain.treasury.QueryKeysRequest", QueryKeysRequest],
    ["/zrchain.treasury.QuerySignTransactionRequestsRequest", QuerySignTransactionRequestsRequest],
    ["/zrchain.treasury.MsgUpdateKeyPolicy", MsgUpdateKeyPolicy],
    ["/zrchain.treasury.QueryZrSignKeysResponse", QueryZrSignKeysResponse],
    ["/zrchain.treasury.QueryKeyByAddressRequest", QueryKeyByAddressRequest],
    ["/zrchain.treasury.TreasuryPacketData", TreasuryPacketData],
    ["/zrchain.treasury.MsgNewKey", MsgNewKey],
    ["/zrchain.treasury.ZenBTCMetadata", ZenBTCMetadata],
    ["/zrchain.treasury.GenesisState", GenesisState],
    ["/zrchain.treasury.Key", Key],
    ["/zrchain.treasury.MsgFulfilICATransactionRequest", MsgFulfilICATransactionRequest],
    ["/zrchain.treasury.MsgUpdateKeyPolicyResponse", MsgUpdateKeyPolicyResponse],
    ["/zrchain.treasury.QueryKeyRequestByIDResponse", QueryKeyRequestByIDResponse],
    ["/zrchain.treasury.SignTxReqResponse", SignTxReqResponse],
    ["/zrchain.treasury.QuerySignatureRequestsRequest", QuerySignatureRequestsRequest],
    ["/zrchain.treasury.QuerySignatureRequestByIDRequest", QuerySignatureRequestByIDRequest],
    ["/zrchain.treasury.QueryZenbtcWalletsRequest", QueryZenbtcWalletsRequest],
    ["/zrchain.treasury.KeyReqResponse", KeyReqResponse],
    ["/zrchain.treasury.MsgNewSignatureRequestResponse", MsgNewSignatureRequestResponse],
    ["/zrchain.treasury.MetadataSolana", MetadataSolana],
    ["/zrchain.treasury.MsgNewZrSignSignatureRequestResponse", MsgNewZrSignSignatureRequestResponse],
    ["/zrchain.treasury.QueryKeyRequestsRequest", QueryKeyRequestsRequest],
    ["/zrchain.treasury.ZrSignKeyEntry", ZrSignKeyEntry],
    ["/zrchain.treasury.QueryKeyByAddressResponse", QueryKeyByAddressResponse],
    ["/zrchain.treasury.KeyResponse", KeyResponse],
    ["/zrchain.treasury.NoData", NoData],
    ["/zrchain.treasury.MsgNewZrSignSignatureRequest", MsgNewZrSignSignatureRequest],
    ["/zrchain.treasury.MsgFulfilICATransactionRequestResponse", MsgFulfilICATransactionRequestResponse],
    ["/zrchain.treasury.ICATransactionRequest", ICATransactionRequest],
    ["/zrchain.treasury.QueryKeyRequestsResponse", QueryKeyRequestsResponse],
    ["/zrchain.treasury.QueryKeyByIDRequest", QueryKeyByIDRequest],
    ["/zrchain.treasury.QueryZenbtcWalletsResponse", QueryZenbtcWalletsResponse],
    ["/zrchain.treasury.QueryKeyRequestByIDRequest", QueryKeyRequestByIDRequest],
    ["/zrchain.treasury.KeyRequest", KeyRequest],
    ["/zrchain.treasury.MsgFulfilSignatureRequest", MsgFulfilSignatureRequest],
    ["/zrchain.treasury.SignTransactionRequest", SignTransactionRequest],
    ["/zrchain.treasury.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.treasury.MsgNewKeyRequestResponse", MsgNewKeyRequestResponse],
    ["/zrchain.treasury.MsgFulfilSignatureRequestResponse", MsgFulfilSignatureRequestResponse],
    
];

export { msgTypes }