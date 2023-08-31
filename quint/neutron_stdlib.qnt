module neutron_stdlib {
    import wasm_stdlib.* from "./wasm_stdlib"

    type IbcFee = {
        recv_fee: List[Coin],
        ack_fee: List[Coin],
        timeout_fee: List[Coin]
    }

    type RequestPacketTimeoutHeight = {
        revision_number: int,
        revision_height: int,
    }

    type NeutronMsg_IbcTransfer = {
        // the port on which the packet will be sent
        source_port: str,
        // the channel by which the packet will be sent
        source_channel: str,
        // the tokens to be transferred
        token: Coin,
        // the sender address
        sender: str,
        // the recipient address on the destination chain
        receiver: str,
        // Timeout height relative to the current block height.
        // The timeout is disabled when set to 0.
        timeout_height: RequestPacketTimeoutHeight,
        // Timeout timestamp in absolute nanoseconds since unix epoch.
        // The timeout is disabled when set to 0.
        timeout_timestamp: u64,
        // Memo to be sent along with transaction.
        memo: str,
        // Fees to refund relayer for different kinds of `SudoMsg` transmission
        // Unused fee types will be returned to msg sender.
        fee: IbcFee,
    }

    // a message that wants a reply
    type SubMsg_IbcTransfer = {
        id: int,
        msg: NeutronMsg_IbcTransfer,
        replyOn: str // always, error, success or never
    }

    // tag should be either "ok" or "error"
    type NeutronResult = { tag: str, messages: List[SubMsg_IbcTransfer], error: str }

    // models the min_ntrn_ibc_fee and query_min_ibc_fee functions
    pure def get_min_fee(): IbcFee = 
        {
            // no recvFee
            recv_fee: List(),
            ack_fee: List({
                denom: "untrn",
                amount: 1250,
            }),
            timeout_fee: List({
                denom: "untrn",
                amount: 500,
            }),
        }
}
