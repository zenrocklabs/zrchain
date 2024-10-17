use cosmwasm_std::{
    entry_point,
    DepsMut,
    Env,
    MessageInfo,
    Response,
    StdResult,
    Uint64,
};

use crate::msg::{
    InstantiateMsg,
};

use zenrock_bindings::msg::{
    SudoMsg,
    ZenrockError,
    SudoResponses,
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
pub fn sudo(deps: DepsMut, _env: Env, msg: SudoMsg) -> Result<Response, ZenrockError> {
    match msg {
        SudoMsg::ParseInput{
            input,
        } => {
            // set ZENROCK_WASM_DEBUG=true to print the debug logs
            deps.api.debug(&format!("input {:?}", input));

            let res = SudoResponses::parse_input_response(Uint64::from(123u64));
            
            // return Ok response on success, Err on failure
            Ok(Response::new().set_data(res))

            // Err(ZenrockError::Unauthorized {})
        }
    }
}