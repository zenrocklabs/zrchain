import { GeneratedType } from "@cosmjs/proto-signing";
import { Policy } from "./types/zrchain/policy/policy";
import { PolicyParticipant } from "./types/zrchain/policy/policy";
import { QueryParamsRequest } from "./types/zrchain/policy/query";
import { QuerySignMethodsByAddressRequest } from "./types/zrchain/policy/query";
import { MsgRemoveMultiGrant } from "./types/zrchain/policy/tx";
import { KeyValue } from "./types/zrchain/policy/action";
import { MsgApproveAction } from "./types/zrchain/policy/tx";
import { MsgAddSignMethod } from "./types/zrchain/policy/tx";
import { Params } from "./types/zrchain/policy/params";
import { ActionResponse } from "./types/zrchain/policy/action";
import { QueryParamsResponse } from "./types/zrchain/policy/query";
import { QueryActionsResponse } from "./types/zrchain/policy/query";
import { QueryPoliciesRequest } from "./types/zrchain/policy/query";
import { MsgUpdateParams } from "./types/zrchain/policy/tx";
import { MsgApproveActionResponse } from "./types/zrchain/policy/tx";
import { MsgRemoveSignMethod } from "./types/zrchain/policy/tx";
import { MsgAddMultiGrantResponse } from "./types/zrchain/policy/tx";
import { QueryPoliciesByCreatorRequest } from "./types/zrchain/policy/query";
import { QueryActionDetailsByIdRequest } from "./types/zrchain/policy/query";
import { MsgNewPolicy } from "./types/zrchain/policy/tx";
import { AdditionalSignaturePasskey } from "./types/zrchain/policy/additional_signature_passkey";
import { GenesisSignMethod } from "./types/zrchain/policy/genesis";
import { NoData } from "./types/zrchain/policy/packet";
import { MsgAddSignMethodResponse } from "./types/zrchain/policy/tx";
import { MsgAddMultiGrant } from "./types/zrchain/policy/tx";
import { GenesisState } from "./types/zrchain/policy/genesis";
import { QueryPolicyByIdResponse } from "./types/zrchain/policy/query";
import { QuerySignMethodsByAddressResponse } from "./types/zrchain/policy/query";
import { MsgRevokeAction } from "./types/zrchain/policy/tx";
import { MsgRemoveSignMethodResponse } from "./types/zrchain/policy/tx";
import { MsgRemoveMultiGrantResponse } from "./types/zrchain/policy/tx";
import { BoolparserPolicy } from "./types/zrchain/policy/policy";
import { QueryActionsRequest } from "./types/zrchain/policy/query";
import { PolicyResponse } from "./types/zrchain/policy/query";
import { QueryPoliciesResponse } from "./types/zrchain/policy/query";
import { QueryPolicyByIdRequest } from "./types/zrchain/policy/query";
import { QueryPoliciesByCreatorResponse } from "./types/zrchain/policy/query";
import { MsgUpdateParamsResponse } from "./types/zrchain/policy/tx";
import { SignMethodPasskey } from "./types/zrchain/policy/sign_method_passkey";
import { PolicyPacketData } from "./types/zrchain/policy/packet";
import { QueryActionDetailsByIdResponse } from "./types/zrchain/policy/query";
import { MsgNewPolicyResponse } from "./types/zrchain/policy/tx";
import { MsgRevokeActionResponse } from "./types/zrchain/policy/tx";
import { Action } from "./types/zrchain/policy/action";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.policy.Policy", Policy],
    ["/zrchain.policy.PolicyParticipant", PolicyParticipant],
    ["/zrchain.policy.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.policy.QuerySignMethodsByAddressRequest", QuerySignMethodsByAddressRequest],
    ["/zrchain.policy.MsgRemoveMultiGrant", MsgRemoveMultiGrant],
    ["/zrchain.policy.KeyValue", KeyValue],
    ["/zrchain.policy.MsgApproveAction", MsgApproveAction],
    ["/zrchain.policy.MsgAddSignMethod", MsgAddSignMethod],
    ["/zrchain.policy.Params", Params],
    ["/zrchain.policy.ActionResponse", ActionResponse],
    ["/zrchain.policy.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.policy.QueryActionsResponse", QueryActionsResponse],
    ["/zrchain.policy.QueryPoliciesRequest", QueryPoliciesRequest],
    ["/zrchain.policy.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.policy.MsgApproveActionResponse", MsgApproveActionResponse],
    ["/zrchain.policy.MsgRemoveSignMethod", MsgRemoveSignMethod],
    ["/zrchain.policy.MsgAddMultiGrantResponse", MsgAddMultiGrantResponse],
    ["/zrchain.policy.QueryPoliciesByCreatorRequest", QueryPoliciesByCreatorRequest],
    ["/zrchain.policy.QueryActionDetailsByIdRequest", QueryActionDetailsByIdRequest],
    ["/zrchain.policy.MsgNewPolicy", MsgNewPolicy],
    ["/zrchain.policy.AdditionalSignaturePasskey", AdditionalSignaturePasskey],
    ["/zrchain.policy.GenesisSignMethod", GenesisSignMethod],
    ["/zrchain.policy.NoData", NoData],
    ["/zrchain.policy.MsgAddSignMethodResponse", MsgAddSignMethodResponse],
    ["/zrchain.policy.MsgAddMultiGrant", MsgAddMultiGrant],
    ["/zrchain.policy.GenesisState", GenesisState],
    ["/zrchain.policy.QueryPolicyByIdResponse", QueryPolicyByIdResponse],
    ["/zrchain.policy.QuerySignMethodsByAddressResponse", QuerySignMethodsByAddressResponse],
    ["/zrchain.policy.MsgRevokeAction", MsgRevokeAction],
    ["/zrchain.policy.MsgRemoveSignMethodResponse", MsgRemoveSignMethodResponse],
    ["/zrchain.policy.MsgRemoveMultiGrantResponse", MsgRemoveMultiGrantResponse],
    ["/zrchain.policy.BoolparserPolicy", BoolparserPolicy],
    ["/zrchain.policy.QueryActionsRequest", QueryActionsRequest],
    ["/zrchain.policy.PolicyResponse", PolicyResponse],
    ["/zrchain.policy.QueryPoliciesResponse", QueryPoliciesResponse],
    ["/zrchain.policy.QueryPolicyByIdRequest", QueryPolicyByIdRequest],
    ["/zrchain.policy.QueryPoliciesByCreatorResponse", QueryPoliciesByCreatorResponse],
    ["/zrchain.policy.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.policy.SignMethodPasskey", SignMethodPasskey],
    ["/zrchain.policy.PolicyPacketData", PolicyPacketData],
    ["/zrchain.policy.QueryActionDetailsByIdResponse", QueryActionDetailsByIdResponse],
    ["/zrchain.policy.MsgNewPolicyResponse", MsgNewPolicyResponse],
    ["/zrchain.policy.MsgRevokeActionResponse", MsgRevokeActionResponse],
    ["/zrchain.policy.Action", Action],
    
];

export { msgTypes }