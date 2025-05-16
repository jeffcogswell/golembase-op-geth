package storagetx

import (
	"fmt"
	"math/big"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/core/types"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/address"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/entity"
	"github.com/jeffcogswell/golembase-op-geth/log"
	"github.com/jeffcogswell/golembase-op-geth/rlp"
	"github.com/holiman/uint256"
)

//go:generate go run ../../rlp/rlpgen -type StorageTransaction -out gen_storage_transaction_rlp.go

// GolemBaseStorageEntityCreated is the event signature for entity creation logs.
var GolemBaseStorageEntityCreated = crypto.Keccak256Hash([]byte("GolemBaseStorageEntityCreated(uint256,uint256)"))

// GolemBaseStorageEntityDeleted is the event signature for entity deletion logs.
var GolemBaseStorageEntityDeleted = crypto.Keccak256Hash([]byte("GolemBaseStorageEntityDeleted(uint256)"))

// GolemBaseStorageEntityUpdated is the event signature for entity update logs.
var GolemBaseStorageEntityUpdated = crypto.Keccak256Hash([]byte("GolemBaseStorageEntityUpdated(uint256,uint256)"))

// GolemBaseStorageEntityTTLExtended is the event signature for extending TTL of an entity.
var GolemBaseStorageEntityTTLExtended = crypto.Keccak256Hash([]byte("GolemBaseStorageEntityTTLExptended(uint256,uint256)"))

// StorageTransaction represents a transaction that can be applied to the storage layer.
// It contains a list of Create operations, a list of Update operations and a list of Delete operations.
//
// Semantics of the transaction operations are as follows:
//   - Create: adds new entities to the storage layer. Each entity has a TTL (number of blocks), a payload and a list of annotations. The Key of the entity is derived from the payload content, the transaction hash where the entity was created and the index of the create operation in the transaction.
//   - Update: updates existing entities. Each entity has a key, a TTL (number of blocks), a payload and a list of annotations. If the entity does not exist, the operation fails, failing the whole transaction.
//   - Delete: removes entities from the storage layer. If the entity does not exist, the operation fails, failing back the whole transaction.
//
// The transaction is atomic, meaning that all operations are applied or none are.
//
// Annotations are key-value pairs where the key is a string and the value is either a string or a number.
// The key-value pairs are used to build indexes and to query the storage layer.
// Same key can have both string and numeric annotation, but not multiple values of the same type.
type StorageTransaction struct {
	Create []Create      `json:"create"`
	Update []Update      `json:"update"`
	Delete []common.Hash `json:"delete"`
	Extend []ExtendTTL   `json:"extend"`
}

type Create struct {
	TTL                uint64                     `json:"ttl"`
	Payload            []byte                     `json:"payload"`
	StringAnnotations  []entity.StringAnnotation  `json:"stringAnnotations"`
	NumericAnnotations []entity.NumericAnnotation `json:"numericAnnotations"`
}

type Update struct {
	EntityKey          common.Hash                `json:"entityKey"`
	TTL                uint64                     `json:"ttl"`
	Payload            []byte                     `json:"payload"`
	StringAnnotations  []entity.StringAnnotation  `json:"stringAnnotations"`
	NumericAnnotations []entity.NumericAnnotation `json:"numericAnnotations"`
}

type ExtendTTL struct {
	EntityKey      common.Hash `json:"entityKey"`
	NumberOfBlocks uint64      `json:"numberOfBlocks"`
}

func (tx *StorageTransaction) Run(blockNumber uint64, txHash common.Hash, sender common.Address, access storageutil.StateAccess) (_ []*types.Log, err error) {

	defer func() {
		if err != nil {
			log.Error("failed to run storage transaction", "error", err)
		}
	}()

	logs := []*types.Log{}

	storeEntity := func(key common.Hash, ap *entity.EntityMetaData, payload []byte, emitLogs bool) error {

		err := entity.Store(access, key, sender, *ap, payload)
		if err != nil {
			return fmt.Errorf("failed to store entity: %w", err)
		}

		if emitLogs {
			expiresAtBlockNumberBig := uint256.NewInt(ap.ExpiresAtBlock)

			data := make([]byte, 32)
			expiresAtBlockNumberBig.PutUint256(data[:32])

			// create the log for the created entity
			log := &types.Log{
				Address:     address.GolemBaseStorageProcessorAddress,
				Topics:      []common.Hash{GolemBaseStorageEntityCreated, key},
				Data:        data,
				BlockNumber: blockNumber,
			}
			logs = append(logs, log)
		}

		return nil

	}

	for i, create := range tx.Create {
		// Convert i to a big integer and pad to 32 bytes
		bigI := big.NewInt(int64(i))
		paddedI := common.LeftPadBytes(bigI.Bytes(), 32)

		key := crypto.Keccak256Hash(txHash.Bytes(), create.Payload, paddedI)

		ap := &entity.EntityMetaData{
			Owner:              sender,
			ExpiresAtBlock:     blockNumber + create.TTL,
			StringAnnotations:  create.StringAnnotations,
			NumericAnnotations: create.NumericAnnotations,
		}

		err := storeEntity(key, ap, create.Payload, true)

		if err != nil {
			return nil, err
		}

	}

	deleteEntity := func(toDelete common.Hash, emitLogs bool) error {

		err := entity.Delete(access, toDelete)
		if err != nil {
			return fmt.Errorf("failed to delete entity: %w", err)
		}

		if emitLogs {

			// create the log for the created entity
			log := &types.Log{
				Address:     address.GolemBaseStorageProcessorAddress,
				Topics:      []common.Hash{GolemBaseStorageEntityDeleted, toDelete},
				Data:        []byte{},
				BlockNumber: blockNumber,
			}

			logs = append(logs, log)
		}

		return nil

	}

	for _, toDelete := range tx.Delete {
		err := deleteEntity(toDelete, true)
		if err != nil {
			return nil, err
		}
	}

	for _, update := range tx.Update {

		oldMetaData, err := entity.GetEntityMetaData(access, update.EntityKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get entity meta data for update %s: %w", update.EntityKey.Hex(), err)
		}

		err = deleteEntity(update.EntityKey, false)
		if err != nil {
			return nil, err
		}

		ap := &entity.EntityMetaData{
			ExpiresAtBlock:     blockNumber + update.TTL,
			StringAnnotations:  update.StringAnnotations,
			NumericAnnotations: update.NumericAnnotations,
			Owner:              oldMetaData.Owner,
		}

		err = storeEntity(update.EntityKey, ap, update.Payload, false)

		if err != nil {
			return nil, err
		}

		expiresAtBlockNumberBig := uint256.NewInt(ap.ExpiresAtBlock)
		data := make([]byte, 32)
		expiresAtBlockNumberBig.PutUint256(data[:32])

		logs = append(logs, &types.Log{
			Address:     address.GolemBaseStorageProcessorAddress,
			Topics:      []common.Hash{GolemBaseStorageEntityUpdated, update.EntityKey},
			Data:        data,
			BlockNumber: blockNumber,
		})

	}

	for _, extend := range tx.Extend {
		newExpiresAtBlock, err := entity.ExtendTTL(access, extend.EntityKey, extend.NumberOfBlocks)
		if err != nil {
			return nil, err
		}

		oldExpiresAtBlock := newExpiresAtBlock - extend.NumberOfBlocks

		oldExpiresAtBlockBig := uint256.NewInt(oldExpiresAtBlock)
		newExpiresAtBlockBig := uint256.NewInt(newExpiresAtBlock)

		data := make([]byte, 64)

		oldExpiresAtBlockBig.PutUint256(data[:32])
		newExpiresAtBlockBig.PutUint256(data[32:])

		logs = append(logs, &types.Log{
			Address:     address.GolemBaseStorageProcessorAddress,
			Topics:      []common.Hash{GolemBaseStorageEntityTTLExtended, extend.EntityKey},
			Data:        data,
			BlockNumber: blockNumber,
		})
	}

	return logs, nil
}

func ExecuteTransaction(d []byte, blockNumber uint64, txHash common.Hash, sender common.Address, access storageutil.StateAccess) ([]*types.Log, error) {
	tx := &StorageTransaction{}
	err := rlp.DecodeBytes(d, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode storage transaction: %w", err)
	}
	logs, err := tx.Run(blockNumber, txHash, sender, access)
	if err != nil {
		log.Error("Failed to run storage transaction", "error", err)
		return nil, fmt.Errorf("failed to run storage transaction: %w", err)
	}
	return logs, nil
}
