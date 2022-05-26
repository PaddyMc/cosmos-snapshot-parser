package parser

import (
	"database/sql"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/plural-labs/cosmos-snapshot-parser/database"
)

func GetAndSaveAccounts(
	ctx sdk.Context,
	AccountKeeper *authkeeper.AccountKeeper,
	db *sql.DB,
) {
	var accounts []authtypes.AccountI

	AccountKeeper.IterateAccounts(ctx, func(acc authtypes.AccountI) (stop bool) {
		accounts = append(accounts, acc)
		return false
	})

	err := database.SaveAccounts(db, accounts)
	if err != nil {
		panic(err)
	}
}
