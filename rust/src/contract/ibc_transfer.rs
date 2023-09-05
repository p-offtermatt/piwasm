use super::neutron_stdlib::*;
use super::wasm_stdlib::*;
use im::HashMap;
use im::HashSet;
use im::Vector;

pub struct InstantiateMsg {
    pub data: String,
}

pub struct ExecuteMsg_Send {
    pub channel: String,
    pub to: String,
    pub denom: String,
    pub amount: u64,
    pub timeout_height: u64,
}

pub fn GetInstantiateMsg() -> InstantiateMsg {
    InstantiateMsg {
        data: "Hello, World!".to_string(),
    }
}

pub const CONTRACT_NAME: String = { "ibc_transfer".to_string() };

pub const CONTRACT_VERSION_STR: String = { "0.1.0".to_string() };

pub struct ContractStorage {
    pub contractVersion: ContractVersion,
    pub replyQueue: HashMap<u64, String>,
    pub runningId: u64,
    pub successfulTransfers: HashSet<Addr>,
}

pub fn instantiate(
    mut curStorage: ContractStorage,
    mut msgInfo: MsgInfo,
    mut msg: InstantiateMsg,
) -> (StdResult, ContractStorage) {
    let result = Todo {
        data: "instantiated".to_string(),
    };
    (StdResult::Ok(result), {
        curStorage.contractVersion = Todo {
            contract: CONTRACT_NAME,
            version: CONTRACT_VERSION_STR,
        };
        curStorage
    })
}

pub fn reply(
    mut env: Env,
    mut msg: Reply,
    mut curStorage: ContractStorage,
) -> (StdResult, ContractStorage) {
    if !curStorage
        .replyQueue
        .keys()
        .collect::<HashSet<_>>()
        .contains_key(&msg.id)
    {
        let error = Todo {
            msg: "got reply to unknown transfer".to_string(),
        };
        (StdResult::Ok(error), curStorage)
    } else {
        let replyTo = curStorage.replyQueue.get(&msg.id).unwrap();
        let s1 = {
            curStorage.replyQueue = curStorage.replyQueue.remove(&msg.id);
            curStorage
        };
        let s2 = {
            s1.successfulTransfers = s1.successfulTransfers.union(im::hashset!(replyTo));
            s1
        };
        let result = Todo {
            data: "got reply to successful transfer".to_string(),
        };
        (StdResult::Ok(result), s2)
    }
}

pub fn execute_send(
    mut msgInfo: MsgInfo,
    mut env: Env,
    mut msg: ExecuteMsg_Send,
    mut curStorage: ContractStorage,
) -> (NeutronResult, ContractStorage) {
    let sender = msgInfo.sender;
    let recipient = msg.to;
    let coin = Todo {
        denom: msg.denom,
        amount: msg.amount,
    };
    let transferMessage = Todo {
        source_port: "transfer".to_string(),
        source_channel: msg.channel,
        sender: env.contract.address,
        receiver: recipient,
        token: coin,
        timeout_height: msg.timeout_height,
        timeout_timestamp: 0_u64,
        memo: "".to_string(),
        fee: get_min_fee,
    };
    let s1 = {
        curStorage.runningId = curStorage.runningId + 1_u64;
        curStorage
    };
    let newId = s1.runningId;
    let newReplyQueue = s1.replyQueue.update(newId, sender);
    let s2 = {
        s1.replyQueue = newReplyQueue;
        s1
    };
    let neutronResult = Todo {
        tag: "ok".to_string(),
        messages: im::vector!(Todo {
            id: newId,
            msg: transferMessage,
            replyOn: "always".to_string(),
        }),
        error: "no error".to_string(),
    };
    (neutronResult, s2)
}
