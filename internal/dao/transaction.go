package dao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/model"
)

type TransactionDao struct{}

var Transaction = &TransactionDao{}

// Insert 添加交易记录
func (d *TransactionDao) Insert(ctx context.Context, transaction *model.Transaction) error {
	_, err := g.DB().Model("transaction").Ctx(ctx).Data(transaction).Insert()
	return err
}

// Update 更新交易记录
func (d *TransactionDao) Update(ctx context.Context, hash string, data g.Map) error {
	_, err := g.DB().Model("transaction").Ctx(ctx).Where("hash", hash).Data(data).Update()
	return err
}

// GetList 获取交易列表
func (d *TransactionDao) GetList(ctx context.Context, address string, page, pageSize int) (list []*model.Transaction, total int, err error) {
	m := g.DB().Model("transaction").
		Where("from_address", address).
		WhereOr("to_address", address)

	// 获取总数
	total, err = m.Ctx(ctx).Count()
	if err != nil {
		return nil, 0, err
	}

	list = make([]*model.Transaction, 0)
	err = m.Ctx(ctx).
		Page(page, pageSize).
		Order("id DESC").
		Scan(&list)

	return list, total, err
}

// GetByHash 根据哈希获取交易
func (d *TransactionDao) GetByHash(ctx context.Context, hash string) (*model.Transaction, error) {
	var tx *model.Transaction
	err := g.DB().Model("transaction").Ctx(ctx).Where("hash", hash).Scan(&tx)
	return tx, err
}
