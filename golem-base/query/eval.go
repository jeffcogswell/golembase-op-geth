package query

import "github.com/jeffcogswell/golembase-op-geth/common"

type DataSource interface {
	GetKeysForStringAnnotation(annotation string, value string) ([]common.Hash, error)
	GetKeysForNumericAnnotation(annotation string, value uint64) ([]common.Hash, error)
}

type Evaluator interface {
	Evaulate(ds DataSource) ([]common.Hash, error)
}
