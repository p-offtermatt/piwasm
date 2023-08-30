// -*- mode: Bluespec; -*-

module ibc_transfer {
    import msg.* from "./msg"
    import wasm_stdlib.* from "./wasm_stdlib"
    import neutron_stdlib.* from "./neutron_stdlib"
    import quint_stdlib.* from "./quint_stdlib"

    // MODEL PARAMETERS

    // the address space
    pure val addresses = Set("alice", "bob", "charlie")

    // the possible token denoms
    pure val tokens = Set("untrn", "uatom", "osmo")

    pure val contractAddress = "ibc_transfer"

    // this is the stuff the contract actually stores in the blockchain
    type ContractStorage = {
        // stores the version of the contract
        contractVersion: ContractVersion,
        // for each sent message, stores
        // as a key the id  of the expected reply,
        // and as value the address of the sender.
        replyQueue: int -> str,
        // a running id for replies.
        runningId: int,
        // keep a set of addresses that have had successful ibc transfers acknowledged
        successfulTransfers: Set[str]
    }

    var storage: ContractStorage

    pure val CONTRACT_NAME = "ibc_transfer"
    pure val CONTRACT_VERSION_STR = "0.1.0"

    pure def instantiate(
        curStorage: ContractStorage,
        msgInfo: MsgInfo,
        msg: InstantiateMsg
    ): (StdResult, ContractStorage) = {
        val result = {data: "instantiated"}
        (Ok(result),
        curStorage.with("contractVersion",
            {contract: CONTRACT_NAME, version: CONTRACT_VERSION_STR}))
    }

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

    pure def reply(env: Env, msg: Reply, curStorage: ContractStorage): StdResult = {
        val id = msg.id
        val replyTo = curStorage.replyQueue.get(id)
        // remove the id from the reply queue
        val s1 = curStorage.with("replyQueue", curStorage.replyQueue.mapRemove(id))
        // add address to the successful transfers
        val s2 = s1.with("successfulTransfers", s1.successfulTransfers.union(Set(replyTo)))
        val result = {data: "got reply to successful transfer"}
        Ok(result)
    }

    action init: bool =
        val emptyStorage = {
            contractVersion: {contract: "", version: ""},
            replyQueue: Map(),
            runningId: 0,
            successfulTransfers: Set()
        }
        val funds = List({
            denom: "untrn",
            amount: 1000000
        })
        val msgInfo = {sender: "alice",
                       funds: funds}
        val instantiateResult = instantiate(emptyStorage, msgInfo, GetInstantiateMsg)
        all {
            storage' = instantiateResult._2,
            instantiateResult._1.tag == "ok"
        }

    action send: bool = 
        nondet sender = oneOf(addresses)
        nondet receiver = oneOf(addresses)
        nondet amount = oneOf(1.to(10))
        nondet denom = oneOf(tokens)
        nondet channel = oneOf(Set("channel-0", "channel-1"))
        nondet timeout_height = oneOf(0.to(10))
        val msgInfo = {sender: sender,
                       funds: List({
                           denom: denom,
                           amount: amount
                       })}
        val env = {contract: {address: contractAddress}}
        val msg = {to: receiver,
                   denom: denom,
                   amount: amount,
                   channel: channel,
                   timeout_height: timeout_height}
        val sendResult = execute_send(msgInfo, env, msg, storage)
        all {
            storage' = sendResult._2,
            sendResult._1.tag == "ok"
        }
}