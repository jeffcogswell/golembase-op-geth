package entityexpiration

import (
	"fmt"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/keyset"
	"github.com/holiman/uint256"
)

type StateAccess = storageutil.StateAccess

var BlockExpirationSalt = []byte("golemBaseExpiresAtBlock")

func AddToEntitiesToExpireAtBlock(access StateAccess, blockNumber uint64, entityKey common.Hash) error {
	expiresAtBlockNumberBig := uint256.NewInt(blockNumber)
	expiredEntityKey := crypto.Keccak256Hash(BlockExpirationSalt, expiresAtBlockNumberBig.Bytes())
	err := keyset.AddValue(access, expiredEntityKey, entityKey)
	if err != nil {
		return fmt.Errorf("failed to append to key list: %w", err)
	}

	return nil
}
