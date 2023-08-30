// -*- mode: Bluespec; -*-

// The contract follows this logic:
// There are three entry points, instantiate, send, and reply.
// The store contains four fields:
// - the contract version
// - a map from message ids of attempted transfers to senders
// - a running id
// - a set of addresses that have had successful ibc transfers acknowledged
// The instantiate action sets the contract version and the running id to 0.
// The send action sends an IBC transfer, and stores the id of the message, mapped to the sender of the transfer.
// The reply action is triggered when the contract receives a reply for a transfer message.
// It will add the sender of the transfer message to the set of successful transfers.

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

    // the address of this contract
    pure val contractAddress = "ibc_transfer"

    // the name of the contract
    pure val CONTRACT_NAME = "ibc_transfer"
    // the version of the contract
    pure val CONTRACT_VERSION_STR = "0.1.0"

    // the default environment
    pure def getEnv(): Env = {
        {contract: {address: contractAddress}}
    }

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

    // the contract storage - the only variable in the contract
    var storage: ContractStorage

    // a helper function for the instantiate action
    pure def instantiate_helper(
        curStorage: ContractStorage,
        msgInfo: MsgInfo,
        msg: InstantiateMsg
    ): (StdResult, ContractStorage) = {
        val result = {data: "instantiated"}
        (Ok(result),
        curStorage.with("contractVersion",
            {contract: CONTRACT_NAME, version: CONTRACT_VERSION_STR}))
    }

    // a helper function for the send action
    pure def execute_send_helper(
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

    // a helper function for the reply action
    pure def reply_helper(env: Env, msg: Reply, curStorage: ContractStorage): (StdResult, ContractStorage) = {
        if (not(curStorage.replyQueue.keys().contains(msg.id))) {
            val error = {msg: "got reply to unknown transfer"}
            (Err(error), curStorage)
        } else {
            val replyTo = curStorage.replyQueue.get(msg.id)
            // remove the id from the reply queue
            val s1 = curStorage.with("replyQueue", curStorage.replyQueue.mapRemove(msg.id))
            // add address to the successful transfers
            val s2 = s1.with("successfulTransfers", s1.successfulTransfers.union(Set(replyTo)))
            val result = {data: "got reply to successful transfer"}
            (Ok(result), s2)
        }
    }

    // init instantiates the contract
    // by starting from an empty storage, then 
    // instantiating it by calling instantiate_helper.
    // the sender and funds for the message are just given a fixed value.
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
        val instantiateResult = instantiate_helper(emptyStorage, msgInfo, GetInstantiateMsg)
        all {
            storage' = instantiateResult._2,
            instantiateResult._1.tag == "ok"
        }

    // deterministic version of send
    // where the nondeterministic parameters are given as arguments
    action send_deterministic(
        sender: Addr, receiver: Addr, amount: int, denom: Denom,
        channel: str,
        timeout_height: int
    ): bool =
        val msgInfo = {sender: sender,
                       funds: List({
                           denom: denom,
                           amount: amount
                       })}
        val msg = {to: receiver,
                   denom: denom,
                   amount: amount,
                   channel: channel,
                   timeout_height: timeout_height}
        val sendResult = execute_send_helper(msgInfo, getEnv, msg, storage)
        all {
            storage' = sendResult._2,
            sendResult._1.tag == "ok"
        }

    // the send action which nondeterministically picks
    // sender, receiver, amount, ..., then calls send_deterministic
    action send: bool = 
        nondet sender = oneOf(addresses)
        nondet receiver = oneOf(addresses)
        nondet amount = oneOf(1.to(10))
        nondet denom = oneOf(tokens)
        nondet channel = oneOf(Set("channel-0", "channel-1"))
        nondet timeout_height = oneOf(0.to(10))
        send_deterministic(sender, receiver, amount, denom, channel, timeout_height)
    
    // deterministic version of reply
    action reply_deterministic(id: int): bool =
        val replyMsg = {id: id, result: Ok({data: "success"})}
        val replyResult = reply_helper(getEnv, replyMsg, storage)
        all {
            storage' = replyResult._2,
            replyResult._1.tag == "ok"
        }

    // Reply to a message with a given id.
    action reply: bool =
        nondet id = oneOf(1.to(10))
        reply_deterministic(id)
}