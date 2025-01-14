package dao

import (
	"context"
	"go-wallet-defi/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GasDao struct{}

var Gas = &GasDao{}

// CreateBatchTransaction 创建批量交易
func (d *GasDao) CreateBatchTransaction(ctx context.Context, tx *model.BatchTransaction) error {
	_, err := g.DB().Model("batch_transaction").Ctx(ctx).Data(tx).Insert()
	return err
}

// UpdateBatchTransaction 更新批量交易
func (d *GasDao) UpdateBatchTransaction(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("batch_transaction").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// CreateGasPrice 创建Gas价格记录
func (d *GasDao) CreateGasPrice(ctx context.Context, price *model.GasPrice) error {
	_, err := g.DB().Model("gas_price").Ctx(ctx).Data(price).Insert()
	return err
}

// GetRecentGasPrices 获取近期Gas价格
func (d *GasDao) GetRecentGasPrices(ctx context.Context, chainId uint64, limit int) ([]*model.GasPrice, error) {
	var prices []*model.GasPrice
	err := g.DB().Model("gas_price").
		Ctx(ctx).
		Where("chain_id", chainId).
		Order("id DESC").
		Limit(limit).
		Scan(&prices)
	return prices, err
}

// UpdateGasStrategy 更新Gas策略
func (d *GasDao) UpdateGasStrategy(ctx context.Context, strategy *model.GasStrategy) error {
	_, err := g.DB().Model("gas_strategy").
		Ctx(ctx).
		Where("chain_id", strategy.ChainId).
		Where("type", strategy.Type).
		Data(strategy).
		Save()
	return err
}

// GetGasStrategy 获取Gas策略
func (d *GasDao) GetGasStrategy(ctx context.Context, chainId uint64, strategyType string) (*model.GasStrategy, error) {
	var strategy *model.GasStrategy
	err := g.DB().Model("gas_strategy").
		Ctx(ctx).
		Where("chain_id", chainId).
		Where("type", strategyType).
		Scan(&strategy)
	return strategy, err
}

// CreateMEVProtection 创建MEV防护记录
func (d *GasDao) CreateMEVProtection(ctx context.Context, protection *model.MEVProtection) error {
	_, err := g.DB().Model("mev_protection").Ctx(ctx).Data(protection).Insert()
	return err
}

// GetMEVProtections 获取MEV防护记录
func (d *GasDao) GetMEVProtections(ctx context.Context, txHash string) ([]*model.MEVProtection, error) {
	var protections []*model.MEVProtection
	err := g.DB().Model("mev_protection").
		Ctx(ctx).
		Where("tx_hash", txHash).
		Scan(&protections)
	return protections, err
}

// CreateTxAcceleration 创建交易加速
func (d *GasDao) CreateTxAcceleration(ctx context.Context, acceleration *model.TxAcceleration) error {
	_, err := g.DB().Model("tx_acceleration").Ctx(ctx).Data(acceleration).Insert()
	return err
}

// UpdateTxAcceleration 更新交易加速
func (d *GasDao) UpdateTxAcceleration(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("tx_acceleration").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}
