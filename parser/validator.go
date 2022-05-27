package parser

import (
	"database/sql"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/PaddyMc/cosmos-snapshot-parser/database"
)

func GetAndSaveValidators(
	ctx sdk.Context,
	StakingKeeper *stakingkeeper.Keeper,
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
	StakingKeeper *stakingkeeper.Keeper,
	db *sql.DB,
	bh int64,
) {
	vals := StakingKeeper.GetAllValidators(ctx)
	err := database.SaveValidatorCommissionData(db, vals, bh)
	if err != nil {
		panic(err)
	}
}

func GetAndSaveValidatorPower(
	ctx sdk.Context,
	StakingKeeper *stakingkeeper.Keeper,
	db *sql.DB,
	bh int64,
) {
	vals := StakingKeeper.GetAllValidators(ctx)
	err := database.SaveValidatorsVotingPowers(db, vals, bh)
	if err != nil {
		panic(err)
	}
}
