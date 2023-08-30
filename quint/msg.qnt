// -*- mode: Bluespec; -*-

// message types specific to this contract.
module msg {

    // Quint does not support empty types, so we need to add a dummy field
    type InstantiateMsg = {data: str}

    // a utility function to get a default InstantiateMsg
    pure def GetInstantiateMsg: InstantiateMsg = {data: "Hello, World!"}

    type ExecuteMsg_Send = {
        channel: str,
        to: str,
        denom: str,
        amount: int,
        // no support for options, so just int here
        timeout_height: int
    }
}