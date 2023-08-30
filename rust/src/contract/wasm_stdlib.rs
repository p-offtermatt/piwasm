use serde::{Deserialize, Serialize};

pub type Addr = cosmwasm_std::Addr;
pub type Denom = String;

pub struct Coin {
    pub denom: Denom,
    pub amount: i64,
}

pub type MsgInfo = cosmwasm_std::MessageInfo;

// pub struct MsgInfo {
//     pub sender: Addr,
//     pub funds: Vec<Coin>,
// }

#[derive(Debug, Clone, Default, Serialize, Deserialize)]
pub struct ContractVersion {
    pub contract: String,
    pub version: String,
}

pub struct Error {
    pub msg: String,
}

pub struct Result {
    pub data: String,
}

pub enum StdResult {
    Ok(Result),
    Err(Error),
}

pub fn Ok(res: Result) -> StdResult {
    StdResult::Ok(res)
}

pub fn Err(error: Error) -> StdResult {
    StdResult::Err(error)
}

pub type ContractInfo = cosmwasm_std::ContractInfo;
// pub struct ContractInfo {
//     pub address: Addr,
// }

pub type Env = cosmwasm_std::Env;
// pub struct Env {
//     pub contract: ContractInfo,
// }

pub struct Reply {
    pub id: i64,
    pub result: StdResult,
}
