use cosmwasm_std::SubMsg;
use neutron_sdk::bindings::msg::NeutronMsg;

use super::wasm_stdlib::*;

use std::vec::Vec;

pub type IbcFee = neutron_sdk::bindings::msg::IbcFee;

// pub struct IbcFee {
//     pub recv_fee: Vec<Coin>,
//     pub ack_fee: Vec<Coin>,
//     pub timeout_fee: Vec<Coin>,
// <Result>}

pub type RequestPacketTimeoutHeight = neutron_sdk::sudo::msg::RequestPacketTimeoutHeight;

// pub struct RequestPacketTimeoutHeight {
//     pub revision_number: i64,
//     pub revision_height: i64,
// }

pub struct NeutronMsg_IbcTransfer {
    pub source_port: String,
    pub source_channel: String,
    pub token: Coin,
    pub sender: Addr,
    pub receiver: Addr,
    pub timeout_height: RequestPacketTimeoutHeight,
    pub timeout_timestamp: u64,
    pub memo: String,
    pub fee: IbcFee,
}

impl From<NeutronMsg_IbcTransfer> for NeutronMsg {
    fn from(msg: NeutronMsg_IbcTransfer) -> Self {
        NeutronMsg::IbcTransfer {
            source_port: msg.source_port,
            source_channel: msg.source_channel,
            token: msg.token.into(),
            sender: msg.sender.to_string(),
            receiver: msg.receiver.to_string(),
            timeout_height: msg.timeout_height.into(),
            timeout_timestamp: msg.timeout_timestamp,
            memo: msg.memo,
            fee: msg.fee.into(),
        }
    }
}

pub struct SubMsg_IbcTransfer {
    pub id: u64,
    pub msg: NeutronMsg_IbcTransfer,
    pub reply_on: String,
}

impl From<SubMsg_IbcTransfer> for SubMsg<NeutronMsg> {
    fn from(msg: SubMsg_IbcTransfer) -> Self {
        match msg.reply_on.as_str() {
            "always" => SubMsg::reply_always(NeutronMsg::from(msg.msg), msg.id),
            "error" => SubMsg::reply_on_error(NeutronMsg::from(msg.msg), msg.id),
            "success" => SubMsg::reply_on_success(NeutronMsg::from(msg.msg), msg.id),
            "never" => SubMsg::new(NeutronMsg::from(msg.msg)),
            _ => panic!("Invalid reply_on value"),
        }
    }
}

pub enum NeutronResult {
    Ok { messages: Vec<SubMsg_IbcTransfer> },
    Error { error: String },
}

impl From<NeutronResult> for cosmwasm_std::StdResult<Vec<SubMsg<NeutronMsg>>> {
    fn from(result: NeutronResult) -> Self {
        match result {
            NeutronResult::Ok { messages } => Ok(messages.into_iter().map(SubMsg::from).collect()),
            NeutronResult::Error { error } => Err(cosmwasm_std::StdError::generic_err(error)),
        }
    }
}

pub fn get_min_fee() -> IbcFee {
    IbcFee {
        recv_fee: Vec::new(),
        ack_fee: vec![Coin {
            denom: "untrn".to_string(),
            amount: 1250_u128.into(),
        }],
        timeout_fee: vec![Coin {
            denom: "untrn".to_string(),
            amount: 500_u128.into(),
        }],
    }
}
