// -*- mode: Bluespec; -*-

module ibc_transfer {
    import msg.* from "./msg"
    import wasm_stdlib.* from "./wasm_stdlib"
    import neutron_stdlib.* from "./neutron_stdlib"

    // this is the stuff the contract actually stores in the blockchain
    type ContractStorage = {
        // stores the version of the contract
        contractVersion: ContractVersion,
        // for each sent message, stores
        // as a key the id  of the expected reply,
        // and as value the address of the sender.
        replyQueue: int -> str,
        // a running id for replies.
        runningId: int
    }

    var storage: ContractStorage

    pure val CONTRACT_NAME = "ibc_transfer"
    pure val CONTRACT_VERSION_STR = "0.1.0"

    // pure def instantiate(
    //     curStorage: ContractStorage,
    //     msgInfo: MsgInfo,
    //     msg: InstantiateMsg
    // ): (StdResult, ContractStorage) = {
    //     (Ok, 
    //     curStorage.with("contractVersion",
    //         {contract: CONTRACT_NAME, version: CONTRACT_VERSION_STR}))
    // }

    pure def execute_send(
        msgInfo: MsgInfo,
        env: Env,
        msg: ExecuteMsg_Send,
        curStorage: ContractStorage
    ): (NeutronResult, ContractStorage) = {
        val sender = msgInfo.sender
        val recipient = msg.to
        val coin = {
            denom: msg.denom,
            amount: msg.amount
        }
        val transferMessage = {
            source_port: "transfer",
            source_channel: msg.channel,
            sender: env.contract.address,
            receiver: recipient,
            token: coin,
            timeout_height: msg.timeout_height,
            timeout_timestamp: 0,
            memo: "",
            fee: get_min_fee
        }
        // update the storage in steps
        // increment runningId
        val s1 = curStorage.with(
            "runningId", curStorage.runningId+1)
        // store the reply address
        val newId = s1.runningId
        val newReplyQueue = s1.replyQueue.put(
            newId, sender)
        val s2 = s1.with("replyQueue", newReplyQueue)
        // construct the neutron result
        val neutronResult = {
            tag: "ok",
            messages: List({
                id: newId,
                msg: transferMessage,
                replyOn: "always"
            }),
            error: {msg: "noerror"}
        }
        (neutronResult, s2)
    }

    // reply goes here
}