pub mod contract;

use contract::{
    ibc_transfer::ContractStorage,
    msg::{ExecuteMsg_Send, InstantiateMsg},
};
use cosmwasm_std::{
    entry_point, Binary, Deps, DepsMut, Empty, Env, MessageInfo, Response, StdError, StdResult,
    Storage,
};
use serde::{de::DeserializeOwned, Serialize};

const STORAGE_KEY: &[u8] = b"storage";

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> StdResult<Response> {
    let initial_storage = ContractStorage::default();
    let (result, storage) = contract::ibc_transfer::instantiate_helper(initial_storage, info, msg);

    match result {
        contract::wasm_stdlib::StdResult::Ok(result) => {
            save(deps.storage, STORAGE_KEY, &storage)?;

            Ok(Response::new().add_attribute("result", result.data))
        }
        contract::wasm_stdlib::StdResult::Err(error) => Err(StdError::generic_err(error.msg)),
    }
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg_Send,
) -> StdResult<Response> {
    let initial_storage = load::<ContractStorage>(deps.storage, STORAGE_KEY)?;
    let (result, storage) =
        contract::ibc_transfer::execute_send_helper(info, env, msg, initial_storage);

    match result {
        contract::neutron_stdlib::NeutronResult::Ok { messages } => {
            save(deps.storage, STORAGE_KEY, &storage)?;

            todo!();
        }
        contract::neutron_stdlib::NeutronResult::Error { error } => {
            Err(StdError::generic_err(error))
        }
    }
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
