package account

import (
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/account/balance"
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/account/create"
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/account/fund"
	"github.com/jeffcogswell/golembase-op-geth/cmd/golembase/account/importkey"
	"github.com/urfave/cli/v2"
)

func Account() *cli.Command {
	return &cli.Command{
		Name:  "account",
		Usage: "Manage accounts",
		Subcommands: []*cli.Command{
			create.Create(),
			fund.FundAccount(),
			balance.AccountBalance(),
			importkey.ImportAccount(),
		},
	}
}
