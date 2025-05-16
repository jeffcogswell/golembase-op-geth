package housekeepingtx

import (
	"fmt"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/core/tracing"
	"github.com/jeffcogswell/golembase-op-geth/core/types"
	"github.com/jeffcogswell/golembase-op-geth/core/vm"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/address"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storagetx"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/entity"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/entity/entityexpiration"
)

func ExecuteTransaction(blockNumber uint64, txHash common.Hash, db vm.StateDB) ([]*types.Log, error) {

	// create the golem base storage processor address if it doesn't exist
	// this is needed to be able to use the state access interface
	if !db.Exist(address.GolemBaseStorageProcessorAddress) {
		db.CreateAccount(address.GolemBaseStorageProcessorAddress)
		db.CreateContract(address.GolemBaseStorageProcessorAddress)
		db.SetNonce(address.GolemBaseStorageProcessorAddress, 1, tracing.NonceChangeNewContract)
	}

	logs := []*types.Log{}

	deleteEntity := func(toDelete common.Hash) error {

		err := entity.Delete(db, toDelete)
		if err != nil {
			return fmt.Errorf("failed to delete entity: %w", err)
		}

		// create the log for the created entity
		log := &types.Log{
			Address:     address.GolemBaseStorageProcessorAddress, // Set the appropriate address if needed
			Topics:      []common.Hash{storagetx.GolemBaseStorageEntityDeleted, toDelete},
			Data:        []byte{},
			BlockNumber: blockNumber,
		}

		logs = append(logs, log)

		return nil
	}

	for key := range entityexpiration.IteratorOfEntitiesToExpireAtBlock(db, blockNumber) {
		err := deleteEntity(key)
		if err != nil {
			return nil, fmt.Errorf("failed to delete entity %s: %w", key.Hex(), err)
		}
	}

	entityexpiration.ClearEntitiesToExpireAtBlock(db, blockNumber)

	return logs, nil
}
