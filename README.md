# piwasm

This is a repo for the quintwasm project, done as part of PI week.
Find some background reading in the background folder.

## Problems

* No sum types - makes options annoying, but is manageable with workarounds
* No generic types - makes it hard to write generic code, e.g. SubMsg<T> is a thing in CosmWasm, but here we need SubMsg_Transfer, SubMsg_Foo, ...
* No inheritance: I don't know of a way to treat an object "like" an object of another class. again makes it difficult to emulate rust code