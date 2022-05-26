package parser

import (
	"database/sql"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/plural-labs/cosmos-snapshot-parser/database"
)

func GetAndSaveSupply(
	ctx sdk.Context,
	BankKeeper bankkeeper.BaseKeeper,
	db *sql.DB,
	bh int64,
) {
	var coins []sdk.Coin

	BankKeeper.IterateTotalSupply(ctx, func(coin sdk.Coin) (stop bool) {
		coins = append(coins, coin)
		return false
	})
	err := database.SaveSupply(db, coins, bh)
	if err != nil {
		panic(err)
	}
}
