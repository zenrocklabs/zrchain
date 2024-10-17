use cosmwasm_schema::{ 
    cw_serde, 
};

use cosmwasm_std::{ 
    Uint64,
    AnyMsg,
    Binary, 
};

use zenrock_bindings::{
    KeyType,
    WalletType,
};

#[cw_serde]
pub struct InstantiateMsg {}

#[cw_serde]
pub enum ExecuteMsg {
    NewWorkspaceRequest {
        admin_policy_id: Uint64,
        sign_policy_id: Uint64,
    },
    NewKeyRequest{
        workspace_addr: String,
        keyring_addr: String,
        key_type: KeyType,
    },
    NewSignDataRequest{
        key_id: Uint64,
        data_for_signing: Binary,
    },
    NewSignTransactionRequest{
        key_id: Uint64,
        wallet_type: WalletType,
        unsigned_transaction: Binary,
        metadata: Option<AnyMsg>,
    },
    AddWorkspaceOwnerRequest {
        workspace_addr: String,
        new_owner: String,
    },
} 
