package parser

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/tx"
)

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
