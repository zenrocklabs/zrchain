import { GeneratedType } from "@cosmjs/proto-signing";
import { BoolparserPolicy } from "./types/zrchain/policy/policy";
import { Action } from "./types/zrchain/policy/action";
import { MsgRemoveSignMethod } from "./types/zrchain/policy/tx";
import { NoData } from "./types/zrchain/policy/packet";
import { MsgRemoveMultiGrant } from "./types/zrchain/policy/tx";
import { QueryActionsRequest } from "./types/zrchain/policy/query";
import { QueryPoliciesRequest } from "./types/zrchain/policy/query";
import { QueryPoliciesResponse } from "./types/zrchain/policy/query";
import { MsgApproveAction } from "./types/zrchain/policy/tx";
import { KeyValue } from "./types/zrchain/policy/action";
import { ActionResponse } from "./types/zrchain/policy/action";
import { QueryParamsRequest } from "./types/zrchain/policy/query";
import { QuerySignMethodsByAddressResponse } from "./types/zrchain/policy/query";
import { QueryPoliciesByCreatorRequest } from "./types/zrchain/policy/query";
import { QueryPoliciesByCreatorResponse } from "./types/zrchain/policy/query";
import { GenesisState } from "./types/zrchain/policy/genesis";
import { MsgNewPolicy } from "./types/zrchain/policy/tx";
import { MsgUpdateParams } from "./types/zrchain/policy/tx";
import { QueryPolicyByIdRequest } from "./types/zrchain/policy/query";
import { QueryPolicyByIdResponse } from "./types/zrchain/policy/query";
import { QuerySignMethodsByAddressRequest } from "./types/zrchain/policy/query";
import { MsgRemoveSignMethodResponse } from "./types/zrchain/policy/tx";
import { MsgAddMultiGrantResponse } from "./types/zrchain/policy/tx";
import { MsgUpdateParamsResponse } from "./types/zrchain/policy/tx";
import { MsgNewPolicyResponse } from "./types/zrchain/policy/tx";
import { MsgRevokeActionResponse } from "./types/zrchain/policy/tx";
import { QueryActionDetailsByIdResponse } from "./types/zrchain/policy/query";
import { MsgRevokeAction } from "./types/zrchain/policy/tx";
import { MsgApproveActionResponse } from "./types/zrchain/policy/tx";
import { MsgAddSignMethodResponse } from "./types/zrchain/policy/tx";
import { PolicyResponse } from "./types/zrchain/policy/query";
import { QueryActionDetailsByIdRequest } from "./types/zrchain/policy/query";
import { PolicyPacketData } from "./types/zrchain/policy/packet";
import { MsgRemoveMultiGrantResponse } from "./types/zrchain/policy/tx";
import { SignMethodPasskey } from "./types/zrchain/policy/sign_method_passkey";
import { MsgAddMultiGrant } from "./types/zrchain/policy/tx";
import { QueryParamsResponse } from "./types/zrchain/policy/query";
import { AdditionalSignaturePasskey } from "./types/zrchain/policy/additional_signature_passkey";
import { MsgAddSignMethod } from "./types/zrchain/policy/tx";
import { Policy } from "./types/zrchain/policy/policy";
import { PolicyParticipant } from "./types/zrchain/policy/policy";
import { Params } from "./types/zrchain/policy/params";
import { GenesisSignMethod } from "./types/zrchain/policy/genesis";
import { QueryActionsResponse } from "./types/zrchain/policy/query";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.policy.BoolparserPolicy", BoolparserPolicy],
    ["/zrchain.policy.Action", Action],
    ["/zrchain.policy.MsgRemoveSignMethod", MsgRemoveSignMethod],
    ["/zrchain.policy.NoData", NoData],
    ["/zrchain.policy.MsgRemoveMultiGrant", MsgRemoveMultiGrant],
    ["/zrchain.policy.QueryActionsRequest", QueryActionsRequest],
    ["/zrchain.policy.QueryPoliciesRequest", QueryPoliciesRequest],
    ["/zrchain.policy.QueryPoliciesResponse", QueryPoliciesResponse],
    ["/zrchain.policy.MsgApproveAction", MsgApproveAction],
    ["/zrchain.policy.KeyValue", KeyValue],
    ["/zrchain.policy.ActionResponse", ActionResponse],
    ["/zrchain.policy.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.policy.QuerySignMethodsByAddressResponse", QuerySignMethodsByAddressResponse],
    ["/zrchain.policy.QueryPoliciesByCreatorRequest", QueryPoliciesByCreatorRequest],
    ["/zrchain.policy.QueryPoliciesByCreatorResponse", QueryPoliciesByCreatorResponse],
    ["/zrchain.policy.GenesisState", GenesisState],
    ["/zrchain.policy.MsgNewPolicy", MsgNewPolicy],
    ["/zrchain.policy.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.policy.QueryPolicyByIdRequest", QueryPolicyByIdRequest],
    ["/zrchain.policy.QueryPolicyByIdResponse", QueryPolicyByIdResponse],
    ["/zrchain.policy.QuerySignMethodsByAddressRequest", QuerySignMethodsByAddressRequest],
    ["/zrchain.policy.MsgRemoveSignMethodResponse", MsgRemoveSignMethodResponse],
    ["/zrchain.policy.MsgAddMultiGrantResponse", MsgAddMultiGrantResponse],
    ["/zrchain.policy.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.policy.MsgNewPolicyResponse", MsgNewPolicyResponse],
    ["/zrchain.policy.MsgRevokeActionResponse", MsgRevokeActionResponse],
    ["/zrchain.policy.QueryActionDetailsByIdResponse", QueryActionDetailsByIdResponse],
    ["/zrchain.policy.MsgRevokeAction", MsgRevokeAction],
    ["/zrchain.policy.MsgApproveActionResponse", MsgApproveActionResponse],
    ["/zrchain.policy.MsgAddSignMethodResponse", MsgAddSignMethodResponse],
    ["/zrchain.policy.PolicyResponse", PolicyResponse],
    ["/zrchain.policy.QueryActionDetailsByIdRequest", QueryActionDetailsByIdRequest],
    ["/zrchain.policy.PolicyPacketData", PolicyPacketData],
    ["/zrchain.policy.MsgRemoveMultiGrantResponse", MsgRemoveMultiGrantResponse],
    ["/zrchain.policy.SignMethodPasskey", SignMethodPasskey],
    ["/zrchain.policy.MsgAddMultiGrant", MsgAddMultiGrant],
    ["/zrchain.policy.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.policy.AdditionalSignaturePasskey", AdditionalSignaturePasskey],
    ["/zrchain.policy.MsgAddSignMethod", MsgAddSignMethod],
    ["/zrchain.policy.Policy", Policy],
    ["/zrchain.policy.PolicyParticipant", PolicyParticipant],
    ["/zrchain.policy.Params", Params],
    ["/zrchain.policy.GenesisSignMethod", GenesisSignMethod],
    ["/zrchain.policy.QueryActionsResponse", QueryActionsResponse],
    
];

export { msgTypes }