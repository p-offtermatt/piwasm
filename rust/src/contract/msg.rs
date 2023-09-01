use cosmwasm_std::Uint128;

use super::{neutron_stdlib::RequestPacketTimeoutHeight, wasm_stdlib::Addr};

pub struct InstantiateMsg {
    pub data: String,
}

pub enum ExecuteMsg {
    Send(ExecuteMsg_Send),
}

pub struct ExecuteMsg_Send {
    pub channel: String,
    pub to: Addr,
    pub denom: String,
    pub amount: Uint128,
    pub timeout_height: RequestPacketTimeoutHeight,
}
