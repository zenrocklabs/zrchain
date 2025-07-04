import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgRemoveWorkspaceOwner } from "./types/zrchain/identity/tx";
import { QueryKeyringByAddressRequest } from "./types/zrchain/identity/query";
import { QueryKeyringByAddressResponse } from "./types/zrchain/identity/query";
import { MsgRemoveKeyringPartyResponse } from "./types/zrchain/identity/tx";
import { MsgUpdateWorkspaceResponse } from "./types/zrchain/identity/tx";
import { MsgAddKeyringParty } from "./types/zrchain/identity/tx";
import { IdentityPacketData } from "./types/zrchain/identity/packet";
import { NoData } from "./types/zrchain/identity/packet";
import { MsgRemoveWorkspaceOwnerResponse } from "./types/zrchain/identity/tx";
import { Workspace } from "./types/zrchain/identity/workspace";
import { MsgDeactivateKeyring } from "./types/zrchain/identity/tx";
import { MsgNewKeyring } from "./types/zrchain/identity/tx";
import { MsgNewWorkspace } from "./types/zrchain/identity/tx";
import { MsgAppendChildWorkspace } from "./types/zrchain/identity/tx";
import { MsgRemoveKeyringAdmin } from "./types/zrchain/identity/tx";
import { QueryWorkspacesRequest } from "./types/zrchain/identity/query";
import { QueryWorkspaceByAddressResponse } from "./types/zrchain/identity/query";
import { MsgRemoveKeyringParty } from "./types/zrchain/identity/tx";
import { MsgAddKeyringPartyResponse } from "./types/zrchain/identity/tx";
import { MsgAddKeyringAdminResponse } from "./types/zrchain/identity/tx";
import { MsgNewChildWorkspace } from "./types/zrchain/identity/tx";
import { MsgAddKeyringAdmin } from "./types/zrchain/identity/tx";
import { QueryParamsRequest } from "./types/zrchain/identity/query";
import { MsgUpdateWorkspace } from "./types/zrchain/identity/tx";
import { MsgAppendChildWorkspaceResponse } from "./types/zrchain/identity/tx";
import { MsgNewChildWorkspaceResponse } from "./types/zrchain/identity/tx";
import { MsgNewKeyringResponse } from "./types/zrchain/identity/tx";
import { Keyring } from "./types/zrchain/identity/keyring";
import { QueryWorkspacesResponse } from "./types/zrchain/identity/query";
import { QueryWorkspaceByAddressRequest } from "./types/zrchain/identity/query";
import { QueryKeyringsRequest } from "./types/zrchain/identity/query";
import { QueryKeyringsResponse } from "./types/zrchain/identity/query";
import { MsgAddWorkspaceOwner } from "./types/zrchain/identity/tx";
import { MsgUpdateParamsResponse } from "./types/zrchain/identity/tx";
import { MsgUpdateKeyringResponse } from "./types/zrchain/identity/tx";
import { Params } from "./types/zrchain/identity/params";
import { GenesisState } from "./types/zrchain/identity/genesis";
import { MsgNewWorkspaceResponse } from "./types/zrchain/identity/tx";
import { MsgAddWorkspaceOwnerResponse } from "./types/zrchain/identity/tx";
import { MsgUpdateParams } from "./types/zrchain/identity/tx";
import { QueryParamsResponse } from "./types/zrchain/identity/query";
import { MsgUpdateKeyring } from "./types/zrchain/identity/tx";
import { MsgRemoveKeyringAdminResponse } from "./types/zrchain/identity/tx";
import { MsgDeactivateKeyringResponse } from "./types/zrchain/identity/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.identity.MsgRemoveWorkspaceOwner", MsgRemoveWorkspaceOwner],
    ["/zrchain.identity.QueryKeyringByAddressRequest", QueryKeyringByAddressRequest],
    ["/zrchain.identity.QueryKeyringByAddressResponse", QueryKeyringByAddressResponse],
    ["/zrchain.identity.MsgRemoveKeyringPartyResponse", MsgRemoveKeyringPartyResponse],
    ["/zrchain.identity.MsgUpdateWorkspaceResponse", MsgUpdateWorkspaceResponse],
    ["/zrchain.identity.MsgAddKeyringParty", MsgAddKeyringParty],
    ["/zrchain.identity.IdentityPacketData", IdentityPacketData],
    ["/zrchain.identity.NoData", NoData],
    ["/zrchain.identity.MsgRemoveWorkspaceOwnerResponse", MsgRemoveWorkspaceOwnerResponse],
    ["/zrchain.identity.Workspace", Workspace],
    ["/zrchain.identity.MsgDeactivateKeyring", MsgDeactivateKeyring],
    ["/zrchain.identity.MsgNewKeyring", MsgNewKeyring],
    ["/zrchain.identity.MsgNewWorkspace", MsgNewWorkspace],
    ["/zrchain.identity.MsgAppendChildWorkspace", MsgAppendChildWorkspace],
    ["/zrchain.identity.MsgRemoveKeyringAdmin", MsgRemoveKeyringAdmin],
    ["/zrchain.identity.QueryWorkspacesRequest", QueryWorkspacesRequest],
    ["/zrchain.identity.QueryWorkspaceByAddressResponse", QueryWorkspaceByAddressResponse],
    ["/zrchain.identity.MsgRemoveKeyringParty", MsgRemoveKeyringParty],
    ["/zrchain.identity.MsgAddKeyringPartyResponse", MsgAddKeyringPartyResponse],
    ["/zrchain.identity.MsgAddKeyringAdminResponse", MsgAddKeyringAdminResponse],
    ["/zrchain.identity.MsgNewChildWorkspace", MsgNewChildWorkspace],
    ["/zrchain.identity.MsgAddKeyringAdmin", MsgAddKeyringAdmin],
    ["/zrchain.identity.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.identity.MsgUpdateWorkspace", MsgUpdateWorkspace],
    ["/zrchain.identity.MsgAppendChildWorkspaceResponse", MsgAppendChildWorkspaceResponse],
    ["/zrchain.identity.MsgNewChildWorkspaceResponse", MsgNewChildWorkspaceResponse],
    ["/zrchain.identity.MsgNewKeyringResponse", MsgNewKeyringResponse],
    ["/zrchain.identity.Keyring", Keyring],
    ["/zrchain.identity.QueryWorkspacesResponse", QueryWorkspacesResponse],
    ["/zrchain.identity.QueryWorkspaceByAddressRequest", QueryWorkspaceByAddressRequest],
    ["/zrchain.identity.QueryKeyringsRequest", QueryKeyringsRequest],
    ["/zrchain.identity.QueryKeyringsResponse", QueryKeyringsResponse],
    ["/zrchain.identity.MsgAddWorkspaceOwner", MsgAddWorkspaceOwner],
    ["/zrchain.identity.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.identity.MsgUpdateKeyringResponse", MsgUpdateKeyringResponse],
    ["/zrchain.identity.Params", Params],
    ["/zrchain.identity.GenesisState", GenesisState],
    ["/zrchain.identity.MsgNewWorkspaceResponse", MsgNewWorkspaceResponse],
    ["/zrchain.identity.MsgAddWorkspaceOwnerResponse", MsgAddWorkspaceOwnerResponse],
    ["/zrchain.identity.MsgUpdateParams", MsgUpdateParams],
    ["/zrchain.identity.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.identity.MsgUpdateKeyring", MsgUpdateKeyring],
    ["/zrchain.identity.MsgRemoveKeyringAdminResponse", MsgRemoveKeyringAdminResponse],
    ["/zrchain.identity.MsgDeactivateKeyringResponse", MsgDeactivateKeyringResponse],
    
];

export { msgTypes }