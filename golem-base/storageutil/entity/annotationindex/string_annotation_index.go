package annotationindex

import (
	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
)

var StringAnnotationIndexSalt = []byte("golemBaseStringAnnotation")
var AnnotationSeparator = []byte("|")

func StringAnnotationIndexKey(key, value string) common.Hash {
	return crypto.Keccak256Hash(StringAnnotationIndexSalt, []byte(key), AnnotationSeparator, []byte(value))
}
