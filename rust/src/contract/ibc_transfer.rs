use std::collections::HashMap;
use std::collections::HashSet;

use serde::{Deserialize, Serialize};

use super::msg::*;
use super::neutron_stdlib::*;
use super::quint_stdlib::*;
use super::wasm_stdlib::*;

#[derive(Debug, Clone, Default, Serialize, Deserialize)]
pub struct ContractStorage {
    pub contractVersion: ContractVersion,
    pub replyQueue: HashMap<i64, Addr>,
    pub runningId: i64,
    pub successfulTransfers: HashSet<Addr>,
}

pub fn addresses() -> HashSet<Addr> {
    let mut addrs = HashSet::new();
    addrs.insert(Addr::unchecked("alice"));
    addrs.insert(Addr::unchecked("bob"));
    addrs.insert(Addr::unchecked("charlie"));
    addrs
}

pub fn tokens() -> HashSet<Denom> {
    let mut denoms = HashSet::new();
    denoms.insert("untrn".to_string());
    denoms.insert("uatom".to_string());
    denoms.insert("osmo".to_string());
    denoms
}

pub fn contractAddress() -> Addr {
    Addr::unchecked("ibc_transfer")
}

pub const CONTRACT_NAME: &str = "ibc_transfer";
pub const CONTRACT_VERSION_STR: &str = "0.1.0";

// pub fn getEnv() -> Env {
//     Env {
//         contract: ContractInfo {
//             address: contractAddress(),
//         },
//         ..Default::default()
//     }
// }

pub fn instantiate_helper(
    mut cur_storage: ContractStorage,
    msg_info: MsgInfo,
    msg: InstantiateMsg,
) -> (StdResult, ContractStorage) {
    let result = Result {
        data: "instantiated".to_string(),
    };
    cur_storage.contractVersion = ContractVersion {
        contract: CONTRACT_NAME.to_string(),
        version: CONTRACT_VERSION_STR.to_string(),
    };
    (Ok(result), cur_storage)
}

pub fn execute_send_helper(
    msgInfo: MsgInfo,
    env: Env,
    msg: ExecuteMsg_Send,
    mut curStorage: ContractStorage,
) -> (NeutronResult, ContractStorage) {
    let sender = &msgInfo.sender;
    let recipient = &msg.to;
    let coin = Coin {
        denom: msg.denom.clone(),
        amount: msg.amount,
    };
    let transferMessage = NeutronMsg_IbcTransfer {
        source_port: "transfer".to_string(),
        source_channel: msg.channel.clone(),
        sender: env.contract.address.clone(),
        receiver: recipient.clone(),
        token: coin,
        timeout_height: msg.timeout_height,
        timeout_timestamp: 0,
        memo: "".to_string(),
        fee: get_min_fee(),
    };
    curStorage.runningId += 1;
    let new_id = curStorage.runningId;
    curStorage.replyQueue.insert(new_id, sender.clone());
    let neutron_result = NeutronResult::Ok {
        messages: vec![SubMsg_IbcTransfer {
            id: new_id,
            msg: transferMessage,
            reply_on: "always".to_string(),
        }],
    };
    (neutron_result, curStorage)
}

pub fn reply_helper(
    env: Env,
    msg: Reply,
    mut curStorage: ContractStorage,
) -> (StdResult, ContractStorage) {
    if !curStorage.replyQueue.contains_key(&msg.id) {
        let error = Error {
            msg: "got reply to unknown transfer".to_string(),
        };

        (Err(error), curStorage)
    } else {
        let replyTo = curStorage.replyQueue.get(&msg.id).unwrap().clone();
        curStorage.replyQueue = mapRemove(&curStorage.replyQueue, &msg.id);
        curStorage.successfulTransfers.insert(replyTo);

        let result = Result {
            data: "got reply to successful transfer".to_string(),
        };

        (Ok(result), curStorage)
    }
}
