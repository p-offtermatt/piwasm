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
    pub to: String,
    pub denom: String,
    pub amount: u128,
    pub timeout_height: u64,
}
