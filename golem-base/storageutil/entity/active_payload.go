package entity

import "github.com/jeffcogswell/golembase-op-geth/common"

//go:generate go run ../../../rlp/rlpgen -type EntityMetaData -out gen_entity_meta_data_rlp.go

// EntityMetaData represents information about an entity that is currently active in the storage layer.
// This is what stored in the state.
// It contains a TTL (number of blocks) and a list of annotations.
// The Key of the entity is derived from the payload content and the transaction hash where the entity was created.

type EntityMetaData struct {
	ExpiresAtBlock     uint64              `json:"expiresAtBlock"`
	StringAnnotations  []StringAnnotation  `json:"stringAnnotations"`
	NumericAnnotations []NumericAnnotation `json:"numericAnnotations"`
	Owner              common.Address      `json:"owner"`
}

type StringAnnotation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NumericAnnotation struct {
	Key   string `json:"key"`
	Value uint64 `json:"value"`
}
