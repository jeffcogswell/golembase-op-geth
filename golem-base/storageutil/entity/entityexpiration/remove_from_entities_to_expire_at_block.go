package entityexpiration

import (
	"fmt"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/keyset"
	"github.com/holiman/uint256"
)

func RemoveFromEntitiesToExpire(access StateAccess, blockNumber uint64, entityKey common.Hash) error {
	expiresAtBlockNumberBig := uint256.NewInt(blockNumber)
	expiredEntityKey := crypto.Keccak256Hash(BlockExpirationSalt, expiresAtBlockNumberBig.Bytes())
	err := keyset.RemoveValue(access, expiredEntityKey, entityKey)
	if err != nil {
		return fmt.Errorf("failed to remove the entity from the key list: %w", err)
	}

	return nil
}
