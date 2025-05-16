package etlworld

import "github.com/jeffcogswell/golembase-op-geth/golem-base/etl/mongodb/mongogolem"

func (w *ETLWorld) AccessMongodb(fn func(*mongogolem.MongoGolem) error) error {
	return fn(w.mongoDriver)
}
