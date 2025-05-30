package entitiesofowner

import (
	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/keyset"
	"github.com/holiman/uint256"
)

type StateAccess = storageutil.StateAccess

var OwnerEntitiesSalt = []byte("golemBase.entitiesOfOwner")

func AddEntity(db StateAccess, owner common.Address, entity common.Hash) error {
	ownerKey := crypto.Keccak256Hash(OwnerEntitiesSalt, owner.Bytes())
	return keyset.AddValue(db, ownerKey, entity)
}

func RemoveEntity(db StateAccess, owner common.Address, entity common.Hash) error {
	ownerKey := crypto.Keccak256Hash(OwnerEntitiesSalt, owner.Bytes())
	return keyset.RemoveValue(db, ownerKey, entity)
}

func Iterate(db StateAccess, owner common.Address) func(yield func(entity common.Hash) bool) {
	ownerKey := crypto.Keccak256Hash(OwnerEntitiesSalt, owner.Bytes())
	return keyset.Iterate(db, ownerKey)
}

func Count(db StateAccess, owner common.Address) *uint256.Int {
	ownerKey := crypto.Keccak256Hash(OwnerEntitiesSalt, owner.Bytes())
	return keyset.Size(db, ownerKey)
}
