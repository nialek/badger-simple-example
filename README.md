## Badger DB simple example

This example uses embedded [BadgerBD](http://github.com/dgraph-io/badger) as an underlying database for a key-value
store, encapsulating transactions.

I've originally intended to use this code in a work project, but decided to use Redis instead, since it has some
functionality missing in Badger (hence you can see that I implemented `store.rename()` before switching).

I felt that it would be unwise to just delete this, maybe I (or you) will use it later in a different project.