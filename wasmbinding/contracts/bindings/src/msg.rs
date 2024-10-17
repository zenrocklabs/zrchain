use cosmwasm_schema::{ 
    cw_serde, 
};

use cosmwasm_std::{ 
    Uint64,
    AnyMsg,
    Binary, 
    CustomMsg, 
    CosmosMsg, 
    StdError,
};

use thiserror::Error;

#[cw_serde]
pub struct InstantiateMsg {}

#[cw_serde]
pub enum KeyType {
    Unspecified = 0,
    Ecdsa = 1,
    Ed25519 = 2,
    Eddsa = 3,
    Bitcoin = 4,
    Btc = 5,
}

#[cw_serde]
pub enum WalletType {
    Unspecified = 0,
    Native = 1,
    Evm = 2,
    BtcTestnet = 3,
    BtcMainnet = 4,
    BtcRegnet = 5,
    Solana = 6,
}

#[cw_serde]
pub enum ZenRockMessage {
    NewWorkspaceRequest {
        creator: String,
        admin_policy_id: Uint64,
        sign_policy_id: Uint64,
    },
    NewKeyRequest{
        creator: String,
        workspace_addr: String,
        keyring_addr: String,
        key_type: KeyType,
    },
    NewSignDataRequest{
        creator: String,
        key_id: Uint64,
        data_for_signing: Binary,
    },
    NewSignTransactionRequest{
        creator: String,
        key_id: Uint64,
        wallet_type: WalletType,
        unsigned_transaction: Binary,
        metadata: Option<AnyMsg>,
    },
    AddWorkspaceOwnerRequest {
        creator: String,
        workspace_addr: String,
        new_owner: String,
    }
} 

impl CustomMsg for ZenRockMessage {}

impl From<ZenRockMessage> for CosmosMsg<ZenRockMessage> {
    fn from(msg: ZenRockMessage) -> CosmosMsg<ZenRockMessage> {
        CosmosMsg::Custom(msg)
    }
}

impl ZenRockMessage {
    pub fn new_key_request( creator: String, workspace_addr: String, keyring_addr: String, key_type: KeyType) -> Self{
        return ZenRockMessage::NewKeyRequest {
            creator: creator,
            workspace_addr: workspace_addr,
            keyring_addr: keyring_addr,
            key_type: key_type,
        };
    }
    
    pub fn new_sign_request( creator: String, key_id: Uint64, data_for_signing: Binary) -> Self {
        return ZenRockMessage::NewSignDataRequest {
            creator: creator,
            key_id: key_id,
            data_for_signing: data_for_signing,
        };
    }
    
    pub fn new_sign_tx_request( creator: String, key_id: Uint64, wallet_type: WalletType, unsigned_transaction: Binary, metadata: Option<AnyMsg>) -> Self{
        return ZenRockMessage::NewSignTransactionRequest {
            creator: creator,
            key_id: key_id,
            wallet_type: wallet_type,
            unsigned_transaction: unsigned_transaction,
            metadata: metadata,
        };
    }
    
    pub fn new_workspace_request( creator: String, admin_policy_id: Uint64, sign_policy_id: Uint64) -> Self{
        return ZenRockMessage::NewWorkspaceRequest {
            creator: creator,
            admin_policy_id: admin_policy_id,
            sign_policy_id: sign_policy_id,
        };
    }
    
    pub fn add_workspace_owner_request(creator: String, workspace_addr: String, new_owner: String) -> Self {
        return ZenRockMessage::AddWorkspaceOwnerRequest {
            creator,
            workspace_addr,
            new_owner
        };
    }
}

#[derive(Error, Debug, PartialEq)]
pub enum ZenrockError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("custom error: {msg:?}")]
    CustomError { msg: String },
}

#[cw_serde]
pub enum SudoResponses {
    ParseInputResponse {
        value: Uint64,
    }
}

#[cw_serde]
pub enum SudoMsg {
    ParseInput {
        input: Binary,
    }
}

impl SudoResponses {
    pub fn parse_input_response (value: Uint64) -> Vec<u8> {
        let response = SudoResponses::ParseInputResponse {
            value: value,
        };
        serde_json::to_string(&response).unwrap().into_bytes()
    }
}