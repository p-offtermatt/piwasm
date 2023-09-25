use serde::{Deserialize, Serialize};

pub type Denom = String;
pub type Addr = cosmwasm_std::Addr;
pub type Coin = cosmwasm_std::Coin;

// pub struct Coin {
//     pub denom: Denom,
//     pub amount: i64,
// }

pub type MsgInfo = cosmwasm_std::MessageInfo;

// pub struct MsgInfo {
//     pub sender: Addr,
//     pub funds: Vec<Coin>,
// }

#[derive(Debug, Clone, Default, PartialEq, Eq, Hash, Serialize, Deserialize)]
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

impl From<StdResult> for cosmwasm_std::StdResult<Result> {
    fn from(result: StdResult) -> Self {
        match result {
            StdResult::Ok(result) => Ok(result),
            StdResult::Err(error) => Err(cosmwasm_std::StdError::generic_err(error.msg)),
        }
    }
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
    pub id: u64,
    pub result: StdResult,
}

impl From<cosmwasm_std::Reply> for Reply {
    fn from(value: cosmwasm_std::Reply) -> Self {
        Reply {
            id: value.id,
            result: StdResult::Ok(Result {
                data: "TODO".to_string(),
            }),
        }
    }
}
