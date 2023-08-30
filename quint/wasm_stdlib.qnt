module wasm_stdlib {
    type Addr = str

    type Coin = {
        denom: str,
        amount: uint128
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

    pure def set_contract_version_deterministic(name: str, version: str, willError: bool): StdResult = 
        if (willError) {
            val err = {msg: "deserialization error"}
            Err(err)
        } else {
            val res = {data: "data"}
            Ok(res)
        }

    // models the set_contract_version function in the wasm stdlib.
    // we do not model its logic, but instead nondeterministically return an error or a success.
    pure def set_contract_version(name: str, version: str): StdResult =
    nondet err = oneOf(Set(true, false))
    set_contract_version_deterministic(name, version, err)

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