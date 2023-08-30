// types and functions that are part of the wasm stdlib.
// this is only a subset, and some types are only partially defined here, e.g. with less fields than
// in the real cosmwasm lib.
module wasm_stdlib {
    type Addr = str

    // this is just str in the real cosmwasm lib, but makes the Quint code more readable.
    type Denom = str

    type Coin = {
        denom: Denom,
        amount: int
    }

    type MsgInfo = {
        sender: Addr,
        funds: List[Coin]
    }

    type ContractVersion = {
        // the crate name of the contract
        contract: str,
        // the version of the contract
        // for example, either semantic versioning like "v0.1.0" or a simple counter
        version: str
    }   

    type Error = {
        msg: str
    }
    type Result = {
        data: str
    }

    // tag is either "ok" or "error".
    // workaround for missing sum types in Quint
    type StdResult = { tag: str, error: Error, res: Result }

    pure def Ok(res: Result): StdResult = {
        tag: "ok",
        res: res,
        error: {msg: "no error"}
    }

    pure def Err(error: Error): StdResult = {
        tag: "error",
        error: error,
        res: {data: "error"}
    }

    type ContractInfo = {
        address: Addr
    }

    type Env = {
        contract: ContractInfo
    }

    type Reply = {
        id: int,
        result: StdResult
    }
}