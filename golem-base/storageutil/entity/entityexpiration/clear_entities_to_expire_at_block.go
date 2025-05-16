package entityexpiration

import (
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/keyset"
	"github.com/holiman/uint256"
)

func ClearEntitiesToExpireAtBlock(access StateAccess, blockNumber uint64) {
	blockNumberBig := uint256.NewInt(blockNumber)
	expiredEntityKey := crypto.Keccak256Hash(BlockExpirationSalt, blockNumberBig.Bytes())
	keyset.Clear(access, expiredEntityKey)
}
