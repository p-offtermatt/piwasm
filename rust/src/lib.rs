#![allow(unused_imports)]

pub mod contract;

use contract::ibc_transfer::{ContractStorage, ExecuteMsg_Send, InstantiateMsg};

use cosmwasm_std::{
    entry_point, DepsMut, Env, MessageInfo, Reply, Response, StdError, StdResult, Storage,
};
use neutron_sdk::bindings::msg::NeutronMsg;
use serde::{de::DeserializeOwned, Serialize};

const STORAGE_KEY: &[u8] = b"storage";

pub enum ExecuteMsg {
    Send(ExecuteMsg_Send),
}

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> StdResult<Response> {
    let initial_storage = ContractStorage::default();
    let (result, storage) = contract::ibc_transfer::instantiate(initial_storage, info, msg);
    let result = StdResult::from(result)?;

    save(deps.storage, STORAGE_KEY, &storage)?;
    Ok(Response::new().add_attribute("result", result.data))
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> StdResult<Response<NeutronMsg>> {
    match msg {
        ExecuteMsg::Send(msg) => execute_send(deps, env, info, msg),
    }
}

pub fn execute_send(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg_Send,
) -> StdResult<Response<NeutronMsg>> {
    let initial_storage = load::<ContractStorage>(deps.storage, STORAGE_KEY)?;
    let (result, storage) = contract::ibc_transfer::execute_send(info, env, msg, initial_storage);

    let messages = StdResult::from(result)?;

    save(deps.storage, STORAGE_KEY, &storage)?;

    let mut response = Response::new();
    for message in messages {
        response = response.add_submessage(message);
    }
    Ok(response)
}

#[entry_point]
pub fn reply(deps: DepsMut, env: Env, msg: Reply) -> StdResult<Response> {
    let initial_storage = load::<ContractStorage>(deps.storage, STORAGE_KEY)?;
    let (result, storage) = contract::ibc_transfer::reply(env, msg.into(), initial_storage);
    let result = StdResult::from(result)?;

    save(deps.storage, STORAGE_KEY, &storage)?;

    Ok(Response::new().add_attribute("result", result.data))
}

fn save<T: Serialize>(storage: &mut dyn Storage, key: &[u8], value: &T) -> StdResult<()> {
    let bytes = postcard::to_allocvec(value)
        .map_err(|e| StdError::generic_err(format!("Error serializing: {e}")))?;

    storage.set(key, bytes.as_slice());

    Ok(())
}

fn load<T: DeserializeOwned>(storage: &dyn Storage, key: &[u8]) -> StdResult<T> {
    let bytes = &storage
        .get(key)
        .ok_or_else(|| StdError::not_found(std::any::type_name::<T>()))?;

    postcard::from_bytes(bytes.as_slice())
        .map_err(|e| StdError::generic_err(format!("Error deserializing: {e}")))
}
