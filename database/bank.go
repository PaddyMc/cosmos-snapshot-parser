package database

import (
	"database/sql"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SaveSupply allows to save for the given height the given total amount of coins
func SaveSupply(db *sql.DB, coins []sdk.Coin, height int64) error {
	stmt := `INSERT INTO supply (denom, amount, height) VALUES`
	var params []interface{}
	for i, coin := range coins {
		pi := i * 3
		stmt += fmt.Sprintf("($%d,$%d,$%d),", pi+1, pi+2, pi+3)
		params = append(params, coin.Denom, coin.Amount.String(), height)
	}

	stmt = stmt[:len(stmt)-1]
	fmt.Println(stmt)
	_, err := db.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing supply: %s", err)
	}

	return nil
}
