package entity

import (
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/entity/create"
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/entity/delete"
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/entity/update"
	"github.com/urfave/cli/v2"
)

func Entity() *cli.Command {
	return &cli.Command{
		Name:  "entity",
		Usage: "Manage entities",
		Subcommands: []*cli.Command{
			create.Create(),
			delete.Delete(),
			update.Update(),
		},
	}
}
