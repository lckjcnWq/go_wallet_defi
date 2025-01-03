package dao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/model"
)

type WalletDao struct{}

var Wallet = WalletDao{}

func (d *WalletDao) Insert(ctx context.Context, wallet *model.Wallet) error {
	_, err := g.DB().Model("wallet").Ctx(ctx).Data(wallet).Insert()
	return err
}

func (d *WalletDao) GetByAddress(ctx context.Context, address string) (*model.Wallet, error) {
	var wallet *model.Wallet
	err := g.DB().Model("wallet").Ctx(ctx).Where("address", address).Scan(&wallet)
	return wallet, err
}
