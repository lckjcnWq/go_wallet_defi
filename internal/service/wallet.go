package service

import (
	"context"
	"go-wallet-defi/internal/model"
)

type IWalletService interface {
	Create(ctx context.Context, chain string) (*model.Wallet, error)
	Import(ctx context.Context, chain, mnemonic string) (*model.Wallet, error)
	GetBalance(ctx context.Context, address string) (string, error)
}
