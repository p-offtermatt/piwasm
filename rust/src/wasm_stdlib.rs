pub type Addr = String;

pub struct Coin {
    pub denom: String,
    pub amount: u128,
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

pub fn set_contract_version_deterministic(
    name: String,
    version: String,
    will_error: bool,
) -> StdResult {
    if will_error {
        let err = Error {
            msg: "deserialization error".to_owned(),
        };
        Err(err)
    } else {
        let res = Result {
            data: "data".to_owned(),
        };
        Ok(res)
    }
}

pub fn set_contract_version(name: String, version: String) -> StdResult {
    let err = rand::random();
    set_contract_version_deterministic(name, version, err)
}

pub struct ContractInfo {
    pub address: Addr,
}

pub struct Env {
    pub contract: ContractInfo,
}

pub struct Reply {
    pub id: i32,
    pub result: StdResult,
}
