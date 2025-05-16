package testutil

import (
	"context"
	"fmt"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/core/types"
	"github.com/jeffcogswell/golembase-op-geth/golem-base/golemtype"
)

// World is the test world - it holds all the state that is shared between steps
type World struct {
	GethInstance     *GethInstance
	FundedAccount    *FundedAccount
	LastReceipt      *types.Receipt
	SearchResult     []golemtype.SearchResult
	CreatedEntityKey common.Hash
	LastError        error
}

func NewWorld(ctx context.Context, gethPath string) (*World, error) {
	geth, err := startGethInstance(ctx, gethPath)
	if err != nil {
		return nil, fmt.Errorf("failed to start geth instance: %w", err)
	}

	acc, err := geth.createAccountAndTransferFunds(ctx, EthToWei(100))
	if err != nil {
		return nil, fmt.Errorf("failed to create account and transfer funds: %w", err)
	}

	return &World{
		GethInstance:  geth,
		FundedAccount: acc,
	}, nil

}

func (w *World) Shutdown() {
	w.GethInstance.shutdown()
}

func (w *World) AddLogsToTestError(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w\n\nGeth Logs:\n%s", err, w.GethInstance.output.String())
}
