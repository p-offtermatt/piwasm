# piwasm

This is a repo for the quintwasm project, done as part of PI week.
Find some background reading in the background folder.

## Usage

To load the model into REPL:
    
    ```
    cd quint
    quint -r ibc_transfer.qnt::ibc_transfer
    ```

To check the invariant:

```
quint run --invariant successfulTransfersWereSent ibc_transfer.qnt --max-samples=200
```


## Problems

* No sum types - makes options annoying, but is manageable with workarounds
* No generic types - makes it hard to write generic code, e.g. SubMsg<T> is a thing in CosmWasm, but here we need SubMsg_Transfer, SubMsg_Foo, ...
* No inheritance: I don't know of a way to treat an object "like" an object of another class. again makes it difficult to emulate rust code