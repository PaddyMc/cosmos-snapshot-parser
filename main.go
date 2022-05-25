package main

import (
	"database/sql"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/plural-labs/cosmos-snapshot-parser/database"

	// store
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	// types :tear:
	ct "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"

	ibchost "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	"github.com/syndtr/goleveldb/leveldb/opt"
	tmstore "github.com/tendermint/tendermint/store"
	db "github.com/tendermint/tm-db"

	// keepers
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/rs/zerolog/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// Various prefixes for accounts and public keys
var (
	AccountAddressPrefix   = "elesto"
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
)

// SetConfig initialize the configuration instance for the sdk
func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}

func main() {
	SetConfig()
	dbDir := "/home/ghost/.elesto/data"

	o := opt.Options{
		DisableSeeksCompaction: true,
	}

	psql, err := database.GetDBConnection()
	if err != nil {
		panic(err)
	}

	// Get ApplicationStore
	appDB, err := db.NewGoLevelDBWithOpts("application", dbDir, &o)
	if err != nil {
		panic(err)
	}
	// Get BlockStore
	blockStoreDB, err := db.NewGoLevelDBWithOpts("blockstore", dbDir, &o)
	if err != nil {
		panic(err)
	}
	blockStore := tmstore.NewBlockStore(blockStoreDB)

	// only mount keys from core sdk
	// todo allow for other keys to be mounted
	keys := types.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		ibchost.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	appStore := rootmulti.NewStore(appDB)
	pruning := storetypes.NewPruningOptions(100, 0, 10)
	appStore.SetPruning(pruning)

	interfaceRegistry := ct.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	distrtypes.RegisterInterfaces(interfaceRegistry)
	evidencetypes.RegisterInterfaces(interfaceRegistry)
	govtypes.RegisterInterfaces(interfaceRegistry)
	slashingtypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)
	upgradetypes.RegisterInterfaces(interfaceRegistry)
	ibctransfertypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	for _, value := range keys {
		appStore.MountStoreWithDB(value, sdk.StoreTypeIAVL, nil)
	}

	err = appStore.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	ctx := sdk.NewContext(
		appStore,
		tmproto.Header{ChainID: "elesto"},
		true,
		server.ZeroLogWrapper{log.Logger},
	)

	paramsKeeper := paramskeeper.NewKeeper(
		marshaler,
		nil,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.StoreKey],
	)

	// module account permissions
	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
	}

	allowedReceivingModAcc := map[string]bool{}

	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	AccountKeeper := authkeeper.NewAccountKeeper(
		marshaler,
		keys[authtypes.StoreKey],
		paramsKeeper.Subspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)

	BankKeeper := bankkeeper.NewBaseKeeper(
		marshaler,
		keys[banktypes.StoreKey],
		AccountKeeper,
		paramsKeeper.Subspace(banktypes.ModuleName),
		blockedAddrs,
	)

	StakingKeeper := stakingkeeper.NewKeeper(
		marshaler,
		keys[stakingtypes.StoreKey],
		AccountKeeper,
		BankKeeper,
		paramsKeeper.Subspace(stakingtypes.ModuleName),
	)

	blockHeight := blockStore.Height()
	fmt.Print(blockStore.Size())

	i := uint64(0)
	for i <= pruning.KeepRecent {
		bh := int64(blockHeight) - int64(i)
		fmt.Println(bh)
		err = appStore.LoadVersion(bh)
		if err != nil {
			panic(err)
		}
		ctx = sdk.NewContext(
			appStore,
			tmproto.Header{ChainID: "elesto"},
			true,
			server.ZeroLogWrapper{log.Logger},
		)

		GetAndSaveBlockData(blockStore, psql, marshaler, bh)
		GetAndSaveValidators(ctx, StakingKeeper, psql, bh)
		GetAndSaveValidatorCommission(ctx, StakingKeeper, psql, bh)
		GetAndSaveAccounts(ctx, AccountKeeper, psql)

		i++
	}

}

func GetAndSaveBlockData(
	blockstore *tmstore.BlockStore,
	db *sql.DB,
	marshaler *codec.ProtoCodec,
	bh int64,
) error {
	block := blockstore.LoadBlock(bh)

	err := database.SaveBlock(db, block)
	if err != nil {
		return err
	}

	for _, msg := range block.Data.Txs {
		transaction, err := UnmarshalTx(marshaler, msg)
		if err != nil {
			panic(err)
		}
		err = database.SaveTx(
			db,
			*marshaler,
			&transaction,
			fmt.Sprintf("%X", []byte(msg.Hash())),
			bh,
		)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func GetAndSaveValidators(
	ctx sdk.Context,
	StakingKeeper stakingkeeper.Keeper,
	db *sql.DB,
	bh int64,
) {
	vals := StakingKeeper.GetAllValidators(ctx)
	err := database.SaveValidatorsData(db, vals, bh)
	if err != nil {
		panic(err)
	}
}

func GetAndSaveValidatorCommission(
	ctx sdk.Context,
	StakingKeeper stakingkeeper.Keeper,
	db *sql.DB,
	bh int64,
) {
	vals := StakingKeeper.GetAllValidators(ctx)
	err := database.SaveValidatorCommissionData(db, vals, bh)
	if err != nil {
		panic(err)
	}
}

func GetAndSaveAccounts(
	ctx sdk.Context,
	AccountKeeper authkeeper.AccountKeeper,
	db *sql.DB,
) {
	var accounts []authtypes.AccountI

	AccountKeeper.IterateAccounts(ctx, func(acc authtypes.AccountI) (stop bool) {
		fmt.Println(acc)
		accounts = append(accounts, acc)
		return false
	})

	err := database.SaveAccounts(db, accounts)
	if err != nil {
		panic(err)
	}
}

func UnmarshalTx(marshaler *codec.ProtoCodec, txBytes []byte) (tx.Tx, error) {
	var raw tx.TxRaw

	err := marshaler.Unmarshal(txBytes, &raw)
	if err != nil {
		return tx.Tx{}, err
	}

	var body tx.TxBody

	err = marshaler.Unmarshal(raw.BodyBytes, &body)
	if err != nil {
		// TODO: fail gracefully...
		return tx.Tx{}, err
	}

	var authInfo tx.AuthInfo

	err = marshaler.Unmarshal(raw.AuthInfoBytes, &authInfo)
	if err != nil {
		return tx.Tx{}, err
	}

	theTx := tx.Tx{
		Body:       &body,
		AuthInfo:   &authInfo,
		Signatures: raw.Signatures,
	}

	return theTx, nil

}
