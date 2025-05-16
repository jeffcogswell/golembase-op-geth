package annotationindex

import (
	"encoding/binary"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
)

var NumericAnnotationIndexSalt = []byte("golemBaseNumericAnnotation")

func NumericAnnotationIndexKey(key string, value uint64) common.Hash {
	return crypto.Keccak256Hash(NumericAnnotationIndexSalt, []byte(key), AnnotationSeparator, binary.BigEndian.AppendUint64(nil, value))
}
