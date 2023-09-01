use std::collections::HashMap;
use std::collections::HashSet;
use super::neutron_stdlib::*;
use super::wasm_stdlib::*;

struct InstantiateMsg {
    data: String,
}

struct ExecuteMsg_Send {
    channel: String,
    to: String,
    denom: String,
    amount: u64,
    timeout_height: u64,
}

fn GetInstantiateMsg() -> InstantiateMsg {
}

pub const CONTRACT_NAME: String = {
}
;

pub const CONTRACT_VERSION_STR: String = {
}
;

struct ContractStorage {
    contractVersion: ContractVersion,
    replyQueue: HashMap::<u64, String>,
    runningId: u64,
    successfulTransfers: HashSet::<Addr>,
}

fn instantiate(curStorage: ContractStorage, msgInfo: MsgInfo, msg: InstantiateMsg) -> (StdResult, ContractStorage) {
}

fn reply(env: Env, msg: Reply, curStorage: ContractStorage) -> (StdResult, ContractStorage) {
}

fn execute_send(msgInfo: MsgInfo, env: Env, msg: ExecuteMsg_Send, curStorage: ContractStorage) -> (NeutronResult, ContractStorage) {
}

