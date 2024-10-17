use cosmwasm_std::{
    entry_point,
    to_json_binary,
    Binary,
    Deps,
    DepsMut,
    Env,
    MessageInfo,
    Response,
    StdResult,
};

use crate::msg::{
    InstantiateMsg,
    ExecuteMsg
};

use crate::query::{
    QueryMsg,
};

use zenrock_bindings::{
    ZenRockMessage,
    ZenRockQuery,
    ZenRockQuerier,
};

#[entry_point]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> StdResult<Response> {
    Ok(Response::new())
}

#[entry_point]
pub fn execute(
    _deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    exec_msg: ExecuteMsg,
) -> StdResult<Response<ZenRockMessage>> {
    let sender = info.sender.to_string();
    let msg: ZenRockMessage;
    match exec_msg {
        ExecuteMsg::NewWorkspaceRequest {
            admin_policy_id,
            sign_policy_id,
        } => msg = ZenRockMessage::new_workspace_request(
            sender,
            admin_policy_id,
            sign_policy_id,
        ),
        ExecuteMsg::NewKeyRequest {
            workspace_addr,
            keyring_addr,
            key_type,
        } => msg = ZenRockMessage::new_key_request(
            sender,
            workspace_addr,
            keyring_addr,
            key_type,
        ),
        ExecuteMsg::NewSignDataRequest {
            key_id,
            data_for_signing,
        } => msg = ZenRockMessage::new_sign_request(
            sender,
            key_id,
            data_for_signing,
        ),
        ExecuteMsg::NewSignTransactionRequest {
            key_id,
            wallet_type,
            unsigned_transaction,
            metadata,
        } => msg = ZenRockMessage::new_sign_tx_request(
            sender,
            key_id,
            wallet_type,
            unsigned_transaction,
            metadata,
        ),
        ExecuteMsg::AddWorkspaceOwnerRequest {
            workspace_addr,
            new_owner,
        } => msg = ZenRockMessage::add_workspace_owner_request(
            sender,
            workspace_addr,
            new_owner
        ),
    }

    let res = Response::new().add_message(msg);
    Ok(res)
}

#[entry_point]
pub fn query(deps: Deps<ZenRockQuery>, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    let querier = ZenRockQuerier::new(&deps.querier);
    match msg {
        QueryMsg::KeyringByAddressQuery{
            keyring_addr
        } => {
            let response = querier.get_keyring_by_address(keyring_addr).unwrap();
            to_json_binary(&response)
        },
    }
}
