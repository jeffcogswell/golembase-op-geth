package entity

import (
	"fmt"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/entity/allentities"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/stateblob"
	"github.com/jeffcogswell/golembase-op-geth/rlp"
)

var EntityMetaDataSalt = []byte("golemBaseEntityMetaData")

func GetEntityMetaData(access StateAccess, key common.Hash) (*EntityMetaData, error) {

	if !allentities.Contains(access, key) {
		return nil, fmt.Errorf("entity %s not found", key.Hex())
	}

	hash := crypto.Keccak256Hash(EntityMetaDataSalt, key[:])
	d := stateblob.GetBlob(access, hash)

	emd := EntityMetaData{}
	err := rlp.DecodeBytes(d, &emd)
	if err != nil {
		return nil, err
	}

	return &emd, nil

}
