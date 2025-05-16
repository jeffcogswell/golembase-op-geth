package entity

import (
	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/stateblob"
)

var PayloadSalt = []byte("golemBasePayload")

func StorePayload(access StateAccess, key common.Hash, payload []byte) {
	hash := crypto.Keccak256Hash(PayloadSalt, key[:])
	stateblob.SetBlob(access, hash, payload)
}
