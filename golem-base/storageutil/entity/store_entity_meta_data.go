package entity

import (
	"bytes"
	"fmt"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/storageutil/stateblob"
	"github.com/jeffcogswell/golembase-op-geth/rlp"
)

func StoreEntityMetaData(access StateAccess, key common.Hash, emd EntityMetaData) error {
	hash := crypto.Keccak256Hash(EntityMetaDataSalt, key[:])

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, &emd)
	if err != nil {
		return fmt.Errorf("failed to encode entity meta data: %w", err)
	}

	stateblob.SetBlob(access, hash, buf.Bytes())
	return nil
}
