package database

import (
	"database/sql"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	maxPostgreSQLParams = 65535
)

func SplitAccounts(accounts []types.AccountI, paramsNumber int) [][]types.AccountI {
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber
	slices := make([][]types.AccountI, len(accounts)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, account := range accounts {
		slices[sliceIndex] = append(slices[sliceIndex], account)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}

// SaveAccounts saves the given accounts inside the database
func SaveAccounts(db *sql.DB, accounts []types.AccountI) error {
	paramsNumber := 1
	slices := SplitAccounts(accounts, paramsNumber)

	for _, accounts := range slices {
		if len(accounts) == 0 {
			continue
		}

		// Store up-to-date data
		err := SaveAccountsDB(db, paramsNumber, accounts)
		if err != nil {
			return fmt.Errorf("error while storing accounts: %s", err)
		}
	}

	return nil
}

func SaveAccountsDB(db *sql.DB, paramsNumber int, accounts []types.AccountI) error {
	if len(accounts) == 0 {
		return nil
	}

	stmt := `INSERT INTO account (address) VALUES `
	var params []interface{}

	for i, account := range accounts {
		ai := i * paramsNumber
		stmt += fmt.Sprintf("($%d),", ai+1)
		params = append(params, account.GetAddress().String())
	}

	stmt = stmt[:len(stmt)-1]
	stmt += " ON CONFLICT DO NOTHING"
	_, err := db.Exec(stmt, params...)
	if err != nil {
		panic(err)
		return fmt.Errorf("error while storing accounts: %s", err)
	}

	return nil
}
