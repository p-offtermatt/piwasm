use crate::msg::*;
use crate::neutron_stdlib::*;
use crate::quint_stdlib::*;
use crate::wasm_stdlib::*;

pub struct ContractStorage {
    pub contract_version: ContractVersion,
    pub reply_queue: std::collections::HashMap<i32, String>,
    pub running_id: i32,
    pub successful_transfers: std::collections::HashSet<String>,
}

pub const CONTRACT_NAME: &str = "ibc_transfer";
pub const CONTRACT_VERSION_STR: &str = "0.1.0";

pub fn instantiate(
    mut cur_storage: ContractStorage,
    msg_info: MsgInfo,
    msg: InstantiateMsg,
) -> (StdResult, ContractStorage) {
    let result = Result {
        data: "instantiated".to_string(),
    };
    cur_storage.contract_version = ContractVersion {
        contract: CONTRACT_NAME.to_string(),
        version: CONTRACT_VERSION_STR.to_string(),
    };
    (Ok(result), cur_storage)
}

pub fn execute_send(
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
    curStorage.running_id += 1;
    let new_id = curStorage.running_id;
    curStorage.reply_queue.insert(new_id, sender.clone());
    let neutron_result = NeutronResult::Ok {
        messages: vec![SubMsg_IbcTransfer {
            id: new_id,
            msg: transferMessage,
            reply_on: "always".to_string(),
        }],
    };
    (neutron_result, curStorage)
}

pub fn reply(
    env: Env,
    msg: Reply,
    mut curStorage: ContractStorage,
) -> (StdResult, ContractStorage) {
    let id = msg.id;
    let reply_to = curStorage.reply_queue.get(&id).unwrap().clone();
    curStorage.reply_queue.remove(&id);
    curStorage.successful_transfers.insert(reply_to);
    let result = Result {
        data: "got reply to successful transfer".to_string(),
    };
    (Ok(result), curStorage)
}
