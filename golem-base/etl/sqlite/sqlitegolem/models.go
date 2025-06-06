// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlitegolem

type Entity struct {
	Key          string
	ExpiresAt    int64
	Payload      []byte
	OwnerAddress string
}

type NumericAnnotation struct {
	EntityKey     string
	AnnotationKey string
	Value         int64
}

type ProcessingStatus struct {
	Network                  string
	LastProcessedBlockNumber int64
	LastProcessedBlockHash   string
}

type StringAnnotation struct {
	EntityKey     string
	AnnotationKey string
	Value         string
}
