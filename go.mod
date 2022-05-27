module github.com/plural-labs/cosmos-snapshot-parser

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.45.4
	github.com/cosmos/ibc-go/v2 v2.2.0
	github.com/dgraph-io/badger/v2 v2.2007.3 // indirect
	github.com/gin-gonic/gin v1.7.0 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lib/pq v1.10.4
	github.com/neilotoole/errgroup v0.1.5
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/rs/zerolog v1.26.0
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
	github.com/tendermint/tendermint v0.34.19
	github.com/tendermint/tm-db v0.6.6
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	google.golang.org/genproto v0.0.0-20220317150908-0efb43f6373e // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/tendermint/tm-db => github.com/tendermint/tm-db v0.6.7-0.20211116222540-a25e8a84a035

// Use Osmosis sdk
replace github.com/cosmos/cosmos-sdk => github.com/osmosis-labs/cosmos-sdk v0.45.1-0.20220524162204-830f277f8259

// Use Osmosis fast iavl
replace github.com/cosmos/iavl => github.com/osmosis-labs/iavl v0.17.3-osmo-v7

// Use osmosis fork of ibc-go
replace github.com/cosmos/ibc-go/v2 => github.com/osmosis-labs/ibc-go/v2 v2.0.2-osmo
