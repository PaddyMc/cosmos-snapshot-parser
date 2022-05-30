# cosmos-snapshot-parser

### Summary

The cosmos-snapshot-parser is a tool that will parse the state of the of a cosmos-sdk datastore, it will then transfer that data into a relational database (psql).

### How to run

First install the binary:

```bash
make install
```

Then run:
```bash
cosmos-snapshot-parser parse \
	--account-prefix osmo
	--connection-string "postgresql://plural:plural@localhost:5432/chain?sslmode=disable" 
	--blocks 1000 
	--db-dir /mnt/<volume>/.osmosisd/data
```

### Config

#### Account Prefix (--account-prefix)

---

This is the account prefix for the chain state that is being parser (e.g osmo).

#### Connection String (--connection-string)

---

The connection string for the psql database (e.g postgresql://plural:plural@localhost:5432/chain?sslmode=disable)

#### Blocks (--blocks)

---

The number of blocks to search backward and parse into a relational database. This widly depends on the pruning config of the node that is being parse. The default pruning configuration of a cosmos node keeps the last 362880 heights. (e.g 362880)

#### Database directory (--db-dir)

---

The directory of the `application.db` and `blockstore.db` of the cosmos-sdk based application. Use the absolute path to avoid any file location issues.


### Standing on the shoulders of giants

Shout out to these projects that I copied code from:

 - https://github.com/binaryholdings/cosmprund
 - https://github.com/forbole/juno
 - https://github.com/forbole/bdjuno
 - https://github.com/allinbits/cosmos-cash

### Beware of Dragons

 - Do Not Use on a production validator state
 - This is alpha tech please validate data and submit an issue if discrepencies found
