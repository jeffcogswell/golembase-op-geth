package main

import (
	"log"
	"os"

	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/account"
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/blocks"
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/cat"
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/entity"
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/query"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "golembase CLI",
		Usage: "Golem Base",

		Commands: []*cli.Command{
			account.Account(),
			entity.Entity(),
			// create.Create(),
			blocks.Blocks(),
			cat.Cat(),
			query.Query(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
