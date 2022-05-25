package database

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	tx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/lib/pq"
	"github.com/tendermint/tendermint/types"
)

var conAddr = "elestovalcons"

// createPartitionIfNotExists creates a new partition having the given partition id if not existing
func createPartitionIfNotExists(
	db *sql.DB,
	table string,
	partitionID int64,
) error {
	partitionTable := fmt.Sprintf("%s_%d", table, partitionID)

	stmt := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s PARTITION OF %s FOR VALUES IN (%d)",
		partitionTable,
		table,
		partitionID,
	)
	_, err := db.Exec(stmt)

	if err != nil {
		return err
	}

	return nil
}

// SaveBlock implements database.Database
func SaveBlock(db *sql.DB, block *types.Block) error {
	sqlStatement := `
INSERT INTO block (height, hash, num_txs, total_gas, proposer_address, timestamp)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`

	valConsAddr, _ := bech32.ConvertAndEncode(
		conAddr,
		block.ProposerAddress.Bytes(),
	)

	proposerAddress := sql.NullString{
		Valid:  len(block.ProposerAddress) != 0,
		String: valConsAddr,
	}

	bHash := block.Hash()

	_, err := db.Exec(
		sqlStatement,
		block.Height,
		bHash.String(),
		len(block.Data.Txs),
		// TODO: calculate gas
		100,
		proposerAddress,
		block.Time,
	)
	return err
}

// SaveTx implements database.Database
func SaveTx(db *sql.DB, marshaler codec.ProtoCodec, tx *tx.Tx, hash string, height int64) error {
	var partitionID int64

	partitionSize := 0
	if partitionSize > 0 {
		partitionID = height / int64(partitionSize)
		err := createPartitionIfNotExists(db, "transaction", partitionID)
		if err != nil {
			return err
		}
	}

	return saveTxInsidePartition(db, marshaler, tx, hash, height, partitionID)
}

// saveTxInsidePartition stores the given transaction inside the partition having the given id
func saveTxInsidePartition(db *sql.DB, marshaler codec.ProtoCodec, tx *tx.Tx, hash string, height int64, partitionId int64) error {
	sqlStatement := `
INSERT INTO transaction
(hash, height, success, messages, memo, signatures, signer_infos, fee, gas_wanted, gas_used, raw_log, logs)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT DO NOTHING`

	var sigs = make([]string, len(tx.Signatures))
	for index, sig := range tx.Signatures {
		sigs[index] = base64.StdEncoding.EncodeToString(sig)
	}

	var msgs = make([]string, len(tx.Body.Messages))
	for index, msg := range tx.Body.Messages {
		bz, err := marshaler.MarshalJSON(msg)
		if err != nil {
			return err
		}
		msgs[index] = string(bz)
	}
	msgsBz := fmt.Sprintf("[%s]", strings.Join(msgs, ","))

	feeBz, err := marshaler.MarshalJSON(tx.AuthInfo.Fee)
	if err != nil {
		return fmt.Errorf("failed to JSON encode tx fee: %s", err)
	}

	var sigInfos = make([]string, len(tx.AuthInfo.SignerInfos))
	for index, info := range tx.AuthInfo.SignerInfos {
		bz, err := marshaler.MarshalJSON(info)
		if err != nil {
			return err
		}
		sigInfos[index] = string(bz)
	}
	sigInfoBz := fmt.Sprintf("[%s]", strings.Join(sigInfos, ","))

	//	logsBz, err := marshaler.MarshalJSON(tx.Logs)
	//	if err != nil {
	//		return err
	//	}

	_, err = db.Exec(sqlStatement,
		hash, height, true,
		msgsBz, tx.Body.Memo, pq.Array(sigs),
		sigInfoBz, string(feeBz),
		1, 1, "{}", string("{}"),
	)
	return err
}
