import { GeneratedType } from "@cosmjs/proto-signing";
import { QueryParamsRequest } from "./types/zrchain/identity/query";
import { QueryWorkspaceByAddressResponse } from "./types/zrchain/identity/query";
import { MsgAddKeyringParty } from "./types/zrchain/identity/tx";
import { MsgAddKeyringPartyResponse } from "./types/zrchain/identity/tx";
import { MsgRemoveKeyringPartyResponse } from "./types/zrchain/identity/tx";
import { MsgRemoveKeyringAdmin } from "./types/zrchain/identity/tx";
import { MsgUpdateWorkspaceResponse } from "./types/zrchain/identity/tx";
import { MsgDeactivateKeyringResponse } from "./types/zrchain/identity/tx";
import { GenesisState } from "./types/zrchain/identity/genesis";
import { QueryKeyringByAddressRequest } from "./types/zrchain/identity/query";
import { Workspace } from "./types/zrchain/identity/workspace";
import { NoData } from "./types/zrchain/identity/packet";
import { MsgNewChildWorkspaceResponse } from "./types/zrchain/identity/tx";
import { MsgNewKeyring } from "./types/zrchain/identity/tx";
import { MsgUpdateKeyringResponse } from "./types/zrchain/identity/tx";
import { QueryWorkspacesResponse } from "./types/zrchain/identity/query";
import { MsgDeactivateKeyring } from "./types/zrchain/identity/tx";
import { MsgNewChildWorkspace } from "./types/zrchain/identity/tx";
import { MsgRemoveWorkspaceOwnerResponse } from "./types/zrchain/identity/tx";
import { Keyring } from "./types/zrchain/identity/keyring";
import { MsgAddWorkspaceOwnerResponse } from "./types/zrchain/identity/tx";
import { MsgNewKeyringResponse } from "./types/zrchain/identity/tx";
import { Params } from "./types/zrchain/identity/params";
import { MsgAddWorkspaceOwner } from "./types/zrchain/identity/tx";
import { MsgNewWorkspaceResponse } from "./types/zrchain/identity/tx";
import { MsgAppendChildWorkspace } from "./types/zrchain/identity/tx";
import { MsgNewWorkspace } from "./types/zrchain/identity/tx";
import { MsgAppendChildWorkspaceResponse } from "./types/zrchain/identity/tx";
import { MsgUpdateKeyring } from "./types/zrchain/identity/tx";
import { MsgAddKeyringAdmin } from "./types/zrchain/identity/tx";
import { MsgRemoveKeyringAdminResponse } from "./types/zrchain/identity/tx";
import { MsgRemoveWorkspaceOwner } from "./types/zrchain/identity/tx";
import { QueryWorkspacesRequest } from "./types/zrchain/identity/query";
import { QueryWorkspaceByAddressRequest } from "./types/zrchain/identity/query";
import { MsgUpdateParamsResponse } from "./types/zrchain/identity/tx";
import { MsgAddKeyringAdminResponse } from "./types/zrchain/identity/tx";
import { QueryParamsResponse } from "./types/zrchain/identity/query";
import { QueryKeyringsRequest } from "./types/zrchain/identity/query";
import { QueryKeyringsResponse } from "./types/zrchain/identity/query";
import { QueryKeyringByAddressResponse } from "./types/zrchain/identity/query";
import { MsgUpdateWorkspace } from "./types/zrchain/identity/tx";
import { MsgRemoveKeyringParty } from "./types/zrchain/identity/tx";
import { IdentityPacketData } from "./types/zrchain/identity/packet";
import { MsgUpdateParams } from "./types/zrchain/identity/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/zrchain.identity.QueryParamsRequest", QueryParamsRequest],
    ["/zrchain.identity.QueryWorkspaceByAddressResponse", QueryWorkspaceByAddressResponse],
    ["/zrchain.identity.MsgAddKeyringParty", MsgAddKeyringParty],
    ["/zrchain.identity.MsgAddKeyringPartyResponse", MsgAddKeyringPartyResponse],
    ["/zrchain.identity.MsgRemoveKeyringPartyResponse", MsgRemoveKeyringPartyResponse],
    ["/zrchain.identity.MsgRemoveKeyringAdmin", MsgRemoveKeyringAdmin],
    ["/zrchain.identity.MsgUpdateWorkspaceResponse", MsgUpdateWorkspaceResponse],
    ["/zrchain.identity.MsgDeactivateKeyringResponse", MsgDeactivateKeyringResponse],
    ["/zrchain.identity.GenesisState", GenesisState],
    ["/zrchain.identity.QueryKeyringByAddressRequest", QueryKeyringByAddressRequest],
    ["/zrchain.identity.Workspace", Workspace],
    ["/zrchain.identity.NoData", NoData],
    ["/zrchain.identity.MsgNewChildWorkspaceResponse", MsgNewChildWorkspaceResponse],
    ["/zrchain.identity.MsgNewKeyring", MsgNewKeyring],
    ["/zrchain.identity.MsgUpdateKeyringResponse", MsgUpdateKeyringResponse],
    ["/zrchain.identity.QueryWorkspacesResponse", QueryWorkspacesResponse],
    ["/zrchain.identity.MsgDeactivateKeyring", MsgDeactivateKeyring],
    ["/zrchain.identity.MsgNewChildWorkspace", MsgNewChildWorkspace],
    ["/zrchain.identity.MsgRemoveWorkspaceOwnerResponse", MsgRemoveWorkspaceOwnerResponse],
    ["/zrchain.identity.Keyring", Keyring],
    ["/zrchain.identity.MsgAddWorkspaceOwnerResponse", MsgAddWorkspaceOwnerResponse],
    ["/zrchain.identity.MsgNewKeyringResponse", MsgNewKeyringResponse],
    ["/zrchain.identity.Params", Params],
    ["/zrchain.identity.MsgAddWorkspaceOwner", MsgAddWorkspaceOwner],
    ["/zrchain.identity.MsgNewWorkspaceResponse", MsgNewWorkspaceResponse],
    ["/zrchain.identity.MsgAppendChildWorkspace", MsgAppendChildWorkspace],
    ["/zrchain.identity.MsgNewWorkspace", MsgNewWorkspace],
    ["/zrchain.identity.MsgAppendChildWorkspaceResponse", MsgAppendChildWorkspaceResponse],
    ["/zrchain.identity.MsgUpdateKeyring", MsgUpdateKeyring],
    ["/zrchain.identity.MsgAddKeyringAdmin", MsgAddKeyringAdmin],
    ["/zrchain.identity.MsgRemoveKeyringAdminResponse", MsgRemoveKeyringAdminResponse],
    ["/zrchain.identity.MsgRemoveWorkspaceOwner", MsgRemoveWorkspaceOwner],
    ["/zrchain.identity.QueryWorkspacesRequest", QueryWorkspacesRequest],
    ["/zrchain.identity.QueryWorkspaceByAddressRequest", QueryWorkspaceByAddressRequest],
    ["/zrchain.identity.MsgUpdateParamsResponse", MsgUpdateParamsResponse],
    ["/zrchain.identity.MsgAddKeyringAdminResponse", MsgAddKeyringAdminResponse],
    ["/zrchain.identity.QueryParamsResponse", QueryParamsResponse],
    ["/zrchain.identity.QueryKeyringsRequest", QueryKeyringsRequest],
    ["/zrchain.identity.QueryKeyringsResponse", QueryKeyringsResponse],
    ["/zrchain.identity.QueryKeyringByAddressResponse", QueryKeyringByAddressResponse],
    ["/zrchain.identity.MsgUpdateWorkspace", MsgUpdateWorkspace],
    ["/zrchain.identity.MsgRemoveKeyringParty", MsgRemoveKeyringParty],
    ["/zrchain.identity.IdentityPacketData", IdentityPacketData],
    ["/zrchain.identity.MsgUpdateParams", MsgUpdateParams],
    
];

export { msgTypes }