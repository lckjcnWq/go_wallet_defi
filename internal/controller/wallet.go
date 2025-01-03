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

func (c *WalletController) Import(ctx context.Context, req *v1.ImportWalletReq) (res *v1.ImportWalletRes, err error) {
	wallet, err := service.Wallet().Import(ctx, req.Chain, req.Mnemonic)
	if err != nil {
		return nil, err
	}

	return &v1.ImportWalletRes{
		Address: wallet.Address,
	}, nil
}

func (c *WalletController) List(ctx context.Context, req *v1.ListWalletReq) (res *v1.ListWalletRes, err error) {
	wallets, err := service.Wallet().GetList(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]v1.WalletInfo, 0, len(wallets))
	for _, w := range wallets {
		list = append(list, v1.WalletInfo{
			Id:      w.Id,
			Address: w.Address,
			Chain:   w.Chain,
			//Balance:   w.Balance,
			CreatedAt: w.CreatedAt,
		})
	}

	return &v1.ListWalletRes{
		List: list,
	}, nil
}

func (c *WalletController) GetBalance(ctx context.Context, req *v1.GetBalanceReq) (res *v1.GetBalanceRes, err error) {
	balance, err := service.Wallet().GetBalance(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	return &v1.GetBalanceRes{
		Balance: balance,
	}, nil
}
