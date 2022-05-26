package parser

import (
	"database/sql"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	tmstore "github.com/tendermint/tendermint/store"

	"github.com/plural-labs/cosmos-snapshot-parser/database"
)

func GetAndSaveBlockData(
	blockstore *tmstore.BlockStore,
	db *sql.DB,
	marshaler *codec.ProtoCodec,
	bh int64,
) error {
	block := blockstore.LoadBlock(bh)

	err := database.SaveBlock(db, block)
	if err != nil {
		return err
	}

	for _, msg := range block.Data.Txs {
		// NOTE: this is where the codec throw an error if not
		// correctly initialized
		transaction, err := UnmarshalTx(marshaler, msg)
		if err != nil {
			panic(err)
		}
		err = database.SaveTx(
			db,
			*marshaler,
			&transaction,
			fmt.Sprintf("%X", msg.Hash()),
			bh,
		)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
