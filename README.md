# piwasm

This is a repo for the quintwasm project, done as part of PI week.
Find some background reading in the background folder.

## Goal

Our goal is to take a .qnt file that models a smart contract, and produce corresponding Rust code for a [CosmWasm](https://book.cosmwasm.com/) smart contract,
specifically with the [Neutron SDK](https://github.com/neutron-org/neutron-sdk/tree/main).

## Current State

We have:
1. Implemented an example contract in Quint
2. Implemented the parts of the standard libraries used by that contract in Quint
3. Done a manual translation from the Quint contract to Rust to check feasibility of automatic translation
4. Started the work on an automatic translation: Much of the actual function bodies are not translated yet, but most boilerplate around them is. The automatic translation goes from the intermediate .json output dumped by the Quint typechecker to an internal representation in our Golang parser and from there to Rust.

In principle, it seems possible to extend our translation to cover the remaining fragments.

## Some technical details

### Conventions

We need to impose some conventions on the Quint code. For now, we split the Quint file into three parts:
1. Entrypoints
2. Tests
3. Utilities

Entrypoints is the part that gets translated into the entry points (API) of the CosmWasm contract. We only allow pure function definitions here, nothing else.
The functions are expected to just return some form of result and a new contract state (they do not modify the contract state themselves - they are the functional layer).

Tests are stuff that is not meant to be part of the contract, but instead constructs that are useful to test the Quint model.
This is the only place where we allow non-pure definitions, and in particular the only place that has a state space.
Actions here typically just call the entry points of the contract, keep a state around, define some constants e.g. for addresses that are used, define invariants, ...

Utilities are functions/vals that are called from the entrypoints, but are not entrypoints themselves.
We allow pure vals and pure defs here - nothing is stateful, since the entire state is only in the tests. This can hence be seen as part of the functional layer.

### Caveats in the translation

Some things in Quint need to be different from the original Rust contracts, e.g. there are no sum types, thus no options in Quint.
This particularly affects our "standard libraries" implementation.
For now, we would aim to do a manual mapping of the standard library anyways and not translate it automatically, so
we can work around this by manually converting the Quint workarounds for sum types to options in Rust.

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

To generate quint output, run:

```
quint typecheck quint/ibc_transfer.qnt --out quint/ibc_transfer.json
```

To run the parser, run:

```
cd parser
go run . ../quint/ibc_transfer_types.json ../rust/src/contract/ibc_transfer.rs
```

## Problems

* No sum types - makes options annoying, but is manageable with workarounds
* No generic types - makes it hard to write generic code, e.g. SubMsg<T> is a thing in CosmWasm, but here we need SubMsg_Transfer, SubMsg_Foo, ...
* No inheritance: I don't know of a way to treat an object "like" an object of another class. again makes it difficult to emulate rust code
