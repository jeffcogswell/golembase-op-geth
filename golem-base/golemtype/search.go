package golemtype

import (
	"github.com/jeffcogswell/golembase-op-geth/common"
)

type SearchResult struct {
	Key   common.Hash `json:"key"`
	Value []byte      `json:"value"`
}
