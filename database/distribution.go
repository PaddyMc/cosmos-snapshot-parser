package database

import (
	"database/sql"
	"fmt"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// SaveValidatorRewards allows to save for the given height the given total amount of coins
func SaveValidatorRewards(
	db *sql.DB,
	validatorAddress string,
	coins distrtypes.ValidatorOutstandingRewards,
	height int64,
) error {
	stmt := `INSERT INTO validator_rewards (denom, amount, height, validator_address) VALUES`
	var params []interface{}
	for i, coin := range coins.Rewards {
		pi := i * 4
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", pi+1, pi+2, pi+3, pi+4)
		params = append(params, coin.Denom, coin.Amount.String(), height, validatorAddress)
		stmt = stmt[:len(stmt)-1]
		_, err := db.Exec(stmt, params...)
		if err != nil {
			return fmt.Errorf("error while storing validator_rewards: %s", err)
		}
	}

	return nil
}
