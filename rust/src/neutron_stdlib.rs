use crate::wasm_stdlib::*;

use std::vec::Vec;

pub struct IbcFee {
    pub recv_fee: Vec<Coin>,
    pub ack_fee: Vec<Coin>,
    pub timeout_fee: Vec<Coin>,
}

pub struct RequestPacketTimeoutHeight {
    pub revision_number: i32,
    pub revision_height: i32,
}

pub struct NeutronMsg_IbcTransfer {
    pub source_port: String,
    pub source_channel: String,
    pub token: Coin,
    pub sender: String,
    pub receiver: String,
    pub timeout_height: u64, // RequestPacketTimeoutHeight,
    pub timeout_timestamp: u64,
    pub memo: String,
    pub fee: IbcFee,
}

pub struct SubMsg_IbcTransfer {
    pub id: i32,
    pub msg: NeutronMsg_IbcTransfer,
    pub reply_on: String,
}

pub enum NeutronResult {
    Ok { messages: Vec<SubMsg_IbcTransfer> },
    Error { error: String },
}

pub fn get_min_fee() -> IbcFee {
    IbcFee {
        recv_fee: Vec::new(),
        ack_fee: vec![Coin {
            denom: "untrn".to_string(),
            amount: 1250,
        }],
        timeout_fee: vec![Coin {
            denom: "untrn".to_string(),
            amount: 500,
        }],
    }
}
