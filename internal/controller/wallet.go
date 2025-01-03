package controller

import (
	"context"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/service"
)

type WalletController struct{}

func (c *WalletController) Create(ctx context.Context, req *v1.CreateWalletReq) (res *v1.CreateWalletRes, err error) {
	wallet, err := service.Wallet().Create(ctx, req.Chain)
	if err != nil {
		return nil, err
	}

	return &v1.CreateWalletRes{
		Address:  wallet.Address,
		Mnemonic: wallet.Mnemonic,
	}, nil
}
