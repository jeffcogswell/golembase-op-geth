package storageutil

import (
	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/address"
)

type StateAccess interface {
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash) common.Hash
}

var GolemDBAddress = address.GolemBaseStorageProcessorAddress
