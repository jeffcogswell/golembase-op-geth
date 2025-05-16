package entity

import (
	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/stateblob"
)

func GetPayload(access StateAccess, key common.Hash) []byte {
	hash := crypto.Keccak256Hash(PayloadSalt, key[:])
	return stateblob.GetBlob(access, hash)
}
