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

// GetList 获取钱包列表
func (d *WalletDao) GetList(ctx context.Context, page, pageSize int) (list []*model.Wallet, err error) {
	m := g.DB().Model("wallet")

	// 获取总数
	_, err = m.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}

	// 不需要返回敏感信息
	list = make([]*model.Wallet, 0)
	err = m.Ctx(ctx).
		Fields("id, user_id, address, chain, wallet_type, created_at, updated_at").
		Page(page, pageSize).
		Order("created_at DESC").
		Scan(&list)

	return list, err
}

// GetListByUserId 根据用户ID获取钱包列表
func (d *WalletDao) GetListByUserId(ctx context.Context, userId uint64, page, pageSize int) (list []*model.Wallet, total int, err error) {
	m := g.DB().Model("wallet").
		Where("user_id", userId)

	// 获取总数
	total, err = m.Ctx(ctx).Count()
	if err != nil {
		return nil, 0, err
	}

	// 不需要返回敏感信息
	list = make([]*model.Wallet, 0)
	err = m.Ctx(ctx).
		Fields("id, user_id, address, chain, wallet_type, created_at, updated_at").
		Page(page, pageSize).
		Order("created_at DESC").
		Scan(&list)

	return list, total, err
}

func (d *WalletDao) CheckAddressExists(ctx context.Context, address string) (bool, error) {
	count, err := g.DB().Model("wallet").Ctx(ctx).Where("address", address).Count()
	return count > 0, err
}

func (d *WalletDao) List(ctx context.Context, page, pageSize int) ([]*model.Wallet, error) {
	var wallets []*model.Wallet
	err := g.DB().Model("wallet").Ctx(ctx).Page(page, pageSize).Order("id DESC").Scan(&wallets)
	return wallets, err
}
