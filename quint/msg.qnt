// -*- mode: Bluespec; -*-

module msg {

    type InstantiateMsg = {data: str}

    pure def GetInstantiateMsg: InstantiateMsg = {data: "Hello, World!"}

    type ExecuteMsg_Send = {
        channel: str,
        to: str,
        denom: str,
        amount: uint128,
        // no support for options, so just uint here
        timeout_height: uint64
    }
}