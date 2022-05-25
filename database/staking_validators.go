package database

import (
	"database/sql"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// SaveValidatorData saves properly the information about the given validator.
func SaveValidatorData(db *sql.DB, validator types.Validator, height int64) error {
	return SaveValidatorsData(db, []types.Validator{validator}, height)
}

// TODO: use in config
var conPub = "elestoconsvalpub"

// SaveValidatorsData allows the bulk saving of a list of validators.
func SaveValidatorsData(db *sql.DB, validators []types.Validator, height int64) error {
	if len(validators) == 0 {
		return nil
	}
	selfDelegationAccQuery := `
INSERT INTO account (address) VALUES `
	var selfDelegationParam []interface{}

	validatorQuery := `
INSERT INTO validator (consensus_address, consensus_pubkey) VALUES `
	var validatorParams []interface{}

	validatorInfoQuery := `
	INSERT INTO validator_info (consensus_address, operator_address, self_delegate_address, max_change_rate, max_rate, height)
	VALUES `
	var validatorInfoParams []interface{}

	for i, validator := range validators {
		vp := i * 2 // Starting position for validator params
		vi := i * 6 // Starting position for validator info params

		selfDelegationAccQuery += fmt.Sprintf("($%d),", i+1)
		selfDelegationParam = append(selfDelegationParam,
			validator.OperatorAddress)

		valPubKey, _ := bech32.ConvertAndEncode(
			conPub,
			validator.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey).Bytes(),
		)

		validatorQuery += fmt.Sprintf("($%d,$%d),", vp+1, vp+2)
		cons, _ := validator.GetConsAddr()
		validatorParams = append(
			validatorParams,
			cons.String(),
			valPubKey,
		)

		validatorInfoQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1,
			vi+2,
			vi+3,
			vi+4,
			vi+5,
			vi+6,
		)
		validatorInfoParams = append(
			validatorInfoParams,
			cons.String(),
			validator.OperatorAddress,
			sdk.AccAddress(validator.GetOperator()).String(),
			validator.Commission.MaxChangeRate.String(),
			validator.Commission.MaxRate.String(),
			height,
		)
	}

	selfDelegationAccQuery = selfDelegationAccQuery[:len(selfDelegationAccQuery)-1] // Remove trailing ","
	selfDelegationAccQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Exec(selfDelegationAccQuery, selfDelegationParam...)
	if err != nil {
		return fmt.Errorf("error while storing accounts: %s", err)
	}

	validatorQuery = validatorQuery[:len(validatorQuery)-1] // Remove trailing ","
	validatorQuery += " ON CONFLICT DO NOTHING"
	result, err := db.Exec(validatorQuery, validatorParams...)
	if err != nil {
		return fmt.Errorf("error while storing valdiators: %s", err)
	}
	fmt.Println(result)

	// Remove the trailing ","
	validatorInfoQuery = validatorInfoQuery[:len(validatorInfoQuery)-1]
	validatorInfoQuery += `
	ON CONFLICT (consensus_address) DO UPDATE
		SET consensus_address = excluded.consensus_address,
			operator_address = excluded.operator_address,
			self_delegate_address = excluded.self_delegate_address,
			max_change_rate = excluded.max_change_rate,
			max_rate = excluded.max_rate,
			height = excluded.height
	WHERE validator_info.height <= excluded.height`
	_, err = db.Exec(validatorInfoQuery, validatorInfoParams...)
	if err != nil {
		return fmt.Errorf("error while storing validator infos: %s", err)
	}

	return nil
}

// SaveValidatorCommission saves a single validator commission.
// It assumes that the delegator address is already present inside the
// proper database table.
func SaveValidatorCommissionData(db *sql.DB, validators []types.Validator, height int64) error {
	for _, val := range validators {
		cons, _ := val.GetConsAddr()
		// Update the current value
		stmt := `
INSERT INTO validator_commission (validator_address, commission, min_self_delegation, height)
VALUES ($1, $2, $3, $4)
ON CONFLICT (validator_address) DO UPDATE
    SET commission = excluded.commission,
        min_self_delegation = excluded.min_self_delegation,
        height = excluded.height
WHERE validator_commission.height <= excluded.height`
		_, err := db.Exec(stmt,
			cons.String(),
			val.Commission.Rate.String(),
			val.MinSelfDelegation.String(),
			height,
		)
		if err != nil {
			return fmt.Errorf("error while storing validator commission: %s", err)
		}
	}

	return nil
}
