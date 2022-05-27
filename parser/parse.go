package parser

import (
	"database/sql"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/rs/zerolog/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmstore "github.com/tendermint/tendermint/store"

	"github.com/PaddyMc/cosmos-snapshot-parser/database"
)

func Parse(
	accountPrefix, dataDir, connectionString string,
	numberOfBlocks uint64,
	marshaler *codec.ProtoCodec,
) error {
	SetConfig(accountPrefix)
	dbDir := dataDir
	connStr := connectionString
	numOfBlocksToParse := numberOfBlocks

	// pruning defies what the ranges for parser to parse
	// this should match the pruning config of the node
	// TODO: for iteration 1 we don't want to worry about this
	//pruning := storetypes.NewPruningOptions(1000, 0, 10)

	// Create all the keepers needed to parse the application store
	_,
		AccountKeeper,
		BankKeeper,
		StakingKeeper,
		_,
		DistrKeeper,
		_,
		keys := CreateKeepers(marshaler)

	// Load the block and application stores
	appStore, blockStore := LoadDataStores(dbDir, keys)

	// Get the database connection (psql)
	psql, err := database.GetDBConnection(connStr)
	if err != nil {
		panic(err)
	}

	// Run the parsing strat...
	err = strat(
		numOfBlocksToParse,
		appStore,
		blockStore,
		psql,
		marshaler,
		AccountKeeper,
		BankKeeper,
		StakingKeeper,
		DistrKeeper,
	)

	return nil
}

func strat(
	numberOfBlocksToParse uint64,
	appStore *rootmulti.Store,
	blockStore *tmstore.BlockStore,
	psql *sql.DB,
	marshaler *codec.ProtoCodec,
	AccountKeeper *authkeeper.AccountKeeper,
	BankKeeper *bankkeeper.BaseKeeper,
	StakingKeeper *stakingkeeper.Keeper,
	DistrKeeper *distrkeeper.Keeper,
) error {
	ctx := sdk.NewContext(
		appStore,
		tmproto.Header{ChainID: ""},
		true,
		// TODO: fix this error...
		server.ZeroLogWrapper{log.Logger},
	)

	blockHeight := blockStore.Height()
	// Does this return the correct value?
	// fmt.Print(blockStore.Size())

	// TODO: fail gracefully here...
	//if pruning.KeepRecent < numOfBlocksToParse {
	//	panic("no enough blocks mate...")
	//}

	// NOTE: We can optimise here using threads, for
	// v1 (MVP) we run synchronously

	// We load the accounts first...
	GetAndSaveAccounts(ctx, AccountKeeper, psql)

	// Then we get the validators at the highest height
	GetAndSaveValidators(ctx, StakingKeeper, psql, blockHeight)
	GetAndSaveValidatorCommission(ctx, StakingKeeper, psql, blockHeight)

	// wow, such for loop...
	i := uint64(0)
	for i <= numberOfBlocksToParse {
		bh := int64(blockHeight) - int64(i)
		fmt.Println(bh)
		err := appStore.LoadVersion(bh)
		if err != nil {
			return err
		}
		ctx = sdk.NewContext(
			appStore,
			tmproto.Header{ChainID: ""},
			true,
			// TODO: fix this error...
			server.ZeroLogWrapper{log.Logger},
		)

		// TODO: params use amino encoding...
		// paramsd := DistrKeeper.GetParams(ctx)
		// fmt.Println(paramsd)

		// This is the basic strat for v1 for each height:
		// 	- Get Blocks
		GetAndSaveBlockData(blockStore, psql, marshaler, bh)
		//	- Get Validator rewards
		GetAndSaveValidatorRewards(ctx, *DistrKeeper, psql, bh)
		//	- Get Supply
		GetAndSaveSupply(ctx, *BankKeeper, psql, bh)
		//	- Get Validator Power
		GetAndSaveValidatorPower(ctx, StakingKeeper, psql, bh)

		i++
	}

	return nil
}
