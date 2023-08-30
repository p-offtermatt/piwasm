use super::wasm_stdlib::Addr;

pub struct InstantiateMsg {
    pub data: String,
}

pub fn GetInstantiateMsg() -> InstantiateMsg {
    InstantiateMsg {
        data: "Hello, world!".to_string(),
    }
}

pub struct ExecuteMsg_Send {
    pub channel: String,
    pub to: Addr,
    pub denom: String,
    pub amount: i64,
    pub timeout_height: i64,
}
