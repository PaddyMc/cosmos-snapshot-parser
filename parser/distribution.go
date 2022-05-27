package parser

import (
	"database/sql"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/plural-labs/cosmos-snapshot-parser/database"
)

func GetAndSaveValidatorRewards(
	ctx sdk.Context,
	DistrKeeper distrkeeper.Keeper,
	db *sql.DB,
	bh int64,
) {
	DistrKeeper.IterateValidatorOutstandingRewards(ctx, func(val sdk.ValAddress, rewards distrtypes.ValidatorOutstandingRewards) (stop bool) {
		err := database.SaveValidatorRewards(db, val.String(), rewards, bh)
		if err != nil {
			panic(err)
		}
		return false
	})
}
