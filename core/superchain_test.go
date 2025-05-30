package core

import (
	"fmt"
	"testing"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/core/rawdb"
	"github.com/jeffcogswell/golembase-op-geth/superchain"
	"github.com/jeffcogswell/golembase-op-geth/triedb"
)

func TestOPStackGenesis(t *testing.T) {
	for id, cfg := range superchain.Chains {
		t.Run(fmt.Sprintf("chain-%s", cfg.Name), func(t *testing.T) {
			t.Parallel()
			_, err := LoadOPStackGenesis(id)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestRegistryChainConfigOverride(t *testing.T) {
	tests := []struct {
		name                 string
		overrides            *ChainOverrides
		setDenominator       *uint64
		expectedDenominator  uint64
		expectedRegolithTime *uint64
	}{
		{
			name:                 "ApplySuperchainUpgrades",
			overrides:            &ChainOverrides{ApplySuperchainUpgrades: true},
			setDenominator:       uint64ptr(50),
			expectedDenominator:  250,
			expectedRegolithTime: uint64ptr(0),
		},
		{
			name:                 "OverrideOptimismCanyon_denom_nil",
			overrides:            &ChainOverrides{OverrideOptimismCanyon: uint64ptr(1)},
			setDenominator:       nil,
			expectedDenominator:  250,
			expectedRegolithTime: nil,
		},
		{
			name:                 "OverrideOptimismCanyon_denom_0",
			overrides:            &ChainOverrides{OverrideOptimismCanyon: uint64ptr(1)},
			setDenominator:       uint64ptr(0),
			expectedDenominator:  250,
			expectedRegolithTime: nil,
		},
		{
			name:                 "OverrideOptimismCanyon_ignore_override",
			overrides:            &ChainOverrides{OverrideOptimismCanyon: uint64ptr(1)},
			setDenominator:       uint64ptr(100),
			expectedDenominator:  100,
			expectedRegolithTime: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db := rawdb.NewMemoryDatabase()
			genesis, err := LoadOPStackGenesis(10)
			if err != nil {
				t.Fatal(err)
			}
			if genesis.Config.RegolithTime == nil {
				t.Fatal("expected non-nil regolith time")
			}
			genesis.Config.RegolithTime = nil

			// initialize the DB
			tdb := triedb.NewDatabase(db, newDbConfig(rawdb.PathScheme))
			genesis.MustCommit(db, tdb)
			bl := genesis.ToBlock()
			rawdb.WriteCanonicalHash(db, bl.Hash(), 0)
			rawdb.WriteBlock(db, bl)

			if genesis.Config.Optimism == nil {
				t.Fatal("expected non nil Optimism config")
			}
			genesis.Config.Optimism.EIP1559DenominatorCanyon = tt.setDenominator
			// create chain config, even with incomplete genesis input: the chain config should be corrected
			chainConfig, _, _, err := SetupGenesisBlockWithOverride(db, tdb, genesis, tt.overrides)
			if err != nil {
				t.Fatal(err)
			}

			// check if we have a corrected chain config
			if tt.expectedRegolithTime == nil {
				if chainConfig.RegolithTime != nil {
					t.Fatal("expected regolith time to be nil")
				}
			} else if *chainConfig.RegolithTime != *tt.expectedRegolithTime {
				t.Fatalf("expected regolith time to be %d, but got %d", *tt.expectedRegolithTime, *chainConfig.RegolithTime)
			}

			if *chainConfig.Optimism.EIP1559DenominatorCanyon != tt.expectedDenominator {
				t.Fatalf("expected EIP1559DenominatorCanyon to be %d, but got %d", tt.expectedDenominator, *chainConfig.Optimism.EIP1559DenominatorCanyon)
			}
		})
	}
}

func TestOPMainnetGenesisDB(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	genesis, err := LoadOPStackGenesis(10)
	if err != nil {
		t.Fatal(err)
	}
	tdb := triedb.NewDatabase(db, newDbConfig(rawdb.PathScheme))
	genesis.MustCommit(db, tdb)
	bl := genesis.ToBlock()
	expected := common.HexToHash("0x7ca38a1916c42007829c55e69d3e9a73265554b586a499015373241b8a3fa48b")
	if blockHash := bl.Hash(); blockHash != expected {
		t.Fatalf("block hash mismatch: %s <> %s", blockHash, expected)
	}
	// This is written separately to the DB by Commit() and is thus tested explicitly here
	canonicalHash := rawdb.ReadCanonicalHash(db, 0)
	if canonicalHash != expected {
		t.Fatalf("canonical hash mismatch: %s <> %s", canonicalHash, expected)
	}
}
