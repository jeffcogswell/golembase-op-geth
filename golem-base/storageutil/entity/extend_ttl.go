package entity

import (
	"fmt"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/entity/entityexpiration"
)

func ExtendTTL(
	access storageutil.StateAccess,
	entityKey common.Hash,
	numberOfBlocks uint64) (uint64, error) {

	entity, err := GetEntityMetaData(access, entityKey)
	if err != nil {
		return 0, err
	}

	err = entityexpiration.RemoveFromEntitiesToExpire(access, entity.ExpiresAtBlock, entityKey)
	if err != nil {
		return 0, fmt.Errorf("failed to remove from entities to expire at block %d: %w", entity.ExpiresAtBlock, err)
	}

	entity.ExpiresAtBlock += numberOfBlocks

	err = entityexpiration.AddToEntitiesToExpireAtBlock(access, entity.ExpiresAtBlock, entityKey)
	if err != nil {
		return 0, fmt.Errorf("failed to add to entities to expire at block %d: %w", entity.ExpiresAtBlock, err)
	}

	err = StoreEntityMetaData(access, entityKey, *entity)
	if err != nil {
		return 0, fmt.Errorf("failed to store entity meta data: %w", err)
	}

	return entity.ExpiresAtBlock, nil

}
