package service

import (
	"context"
	"go-wallet-defi/internal/model"
)

// 定义接口管理器
type IWallet interface {
	Create(ctx context.Context, chain string) (*model.Wallet, error)
	Import(ctx context.Context, chain, mnemonic string) (*model.Wallet, error)
	GetList(ctx context.Context, page, pageSize int) ([]*model.Wallet, error)
	GetBalance(ctx context.Context, address string) (string, error)
}

var (
	localWallet IWallet
)

func Wallet() IWallet {
	if localWallet == nil {
		panic("implement not found for interface IWallet, forgot register?")
	}
	return localWallet
}

// RegisterWallet 注册接口实现类
func RegisterWallet(i IWallet) {
	localWallet = i
}
