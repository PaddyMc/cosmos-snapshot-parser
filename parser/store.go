package parser

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	"github.com/syndtr/goleveldb/leveldb/opt"
	tmstore "github.com/tendermint/tendermint/store"
	db "github.com/tendermint/tm-db"
)

func LoadDataStores(dbDir string, keys map[string]*types.KVStoreKey) (
	appStore *rootmulti.Store,
	blockStore *tmstore.BlockStore,
) {
	o := opt.Options{
		DisableSeeksCompaction: true,
	}
	// Get the application store from a directory
	appDB, err := db.NewGoLevelDBWithOpts("application", dbDir, &o)
	appStore = rootmulti.NewStore(appDB)
	if err != nil {
		panic(err)
	}

	// Get the block store from a directory
	blockStoreDB, err := db.NewGoLevelDBWithOpts("blockstore", dbDir, &o)
	if err != nil {
		panic(err)
	}
	blockStore = tmstore.NewBlockStore(blockStoreDB)

	for _, value := range keys {
		appStore.MountStoreWithDB(value, sdk.StoreTypeIAVL, nil)
	}

	// Load the latest version in the state
	err = appStore.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	return
}

func CreateKeepers(marshaler *codec.ProtoCodec) (
	pk *paramskeeper.Keeper,
	ak *authkeeper.AccountKeeper,
	bk *bankkeeper.BaseKeeper,
	sk *stakingkeeper.Keeper,
	mk *mintkeeper.Keeper,
	dk *distrkeeper.Keeper,
	slk *slashingkeeper.Keeper,
	keys map[string]*types.KVStoreKey,
) {

	// todo allow for other keys to be mounted
	keys = types.NewKVStoreKeys(
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

	paramsKeeper := paramskeeper.NewKeeper(
		marshaler,
		nil,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.StoreKey],
	)

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

	MintKeeper := mintkeeper.NewKeeper(
		marshaler,
		keys[minttypes.StoreKey],
		paramsKeeper.Subspace(minttypes.ModuleName),
		&StakingKeeper,
		AccountKeeper,
		BankKeeper,
		authtypes.FeeCollectorName,
	)

	DistrKeeper := distrkeeper.NewKeeper(
		marshaler,
		keys[distrtypes.StoreKey],
		paramsKeeper.Subspace(distrtypes.ModuleName),
		AccountKeeper,
		BankKeeper,
		&StakingKeeper,
		authtypes.FeeCollectorName,
		blockedAddrs,
	)

	SlashingKeeper := slashingkeeper.NewKeeper(
		marshaler,
		keys[slashingtypes.StoreKey],
		&StakingKeeper,
		paramsKeeper.Subspace(slashingtypes.ModuleName),
	)

	return &paramsKeeper,
		&AccountKeeper,
		&BankKeeper,
		&StakingKeeper,
		&MintKeeper,
		&DistrKeeper,
		&SlashingKeeper,
		keys
}
