pub type Addr = String;
pub type Denom = String;

pub struct Coin {
    pub denom: Denom,
    pub amount: i64,
}

pub struct MsgInfo {
    pub sender: Addr,
    pub funds: Vec<Coin>,
}

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

pub struct ContractInfo {
    pub address: Addr,
}

pub struct Env {
    pub contract: ContractInfo,
}

pub struct Reply {
    pub id: i64,
    pub result: StdResult,
}
