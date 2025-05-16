package entity

import (
	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/stateblob"
)

func DeletePayload(access StateAccess, key common.Hash) {
	hash := crypto.Keccak256Hash(PayloadSalt, key[:])
	stateblob.DeleteBlob(access, hash)
}
