package useraccount

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
)

type UserAccount struct {
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
}

func Load() (*UserAccount, error) {
	privageKeyPath, err := xdg.ConfigFile(PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get config file path: %w", err)
	}

	privageKeyBytes, err := os.ReadFile(privageKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := crypto.ToECDSA(privageKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize private key: %w", err)
	}

	return &UserAccount{
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}, nil

}
