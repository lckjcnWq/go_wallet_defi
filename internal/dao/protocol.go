package dao

import (
	"context"
	"go-wallet-defi/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type ProtocolDao struct{}

var Protocol = &ProtocolDao{}

// CreateSwapTransaction 创建DEX交易
func (d *ProtocolDao) CreateSwapTransaction(ctx context.Context, tx *model.SwapTransaction) error {
	_, err := g.DB().Model("swap_transaction").Data(tx).Insert()
	return err
}

// UpdateSwapTransaction 更新DEX交易
func (d *ProtocolDao) UpdateSwapTransaction(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("swap_transaction").Where("id", id).Data(data).Update()
	return err
}

// UpdateLendingPosition 更新借贷仓位
func (d *ProtocolDao) UpdateLendingPosition(ctx context.Context, position *model.LendingPosition) error {
	_, err := g.DB().Model("lending_position").
		Where("chain_id", position.ChainId).
		Where("protocol", position.Protocol).
		Where("address", position.Address).
		Where("token", position.Token).
		Data(position).Save()
	return err
}

// GetLendingPositions 获取借贷仓位
func (d *ProtocolDao) GetLendingPositions(ctx context.Context, address string) ([]*model.LendingPosition, error) {
	var positions []*model.LendingPosition
	err := g.DB().Model("lending_position").Where("address", address).Scan(&positions)
	return positions, err
}

// CreateNFTTransaction 创建NFT交易
func (d *ProtocolDao) CreateNFTTransaction(ctx context.Context, tx *model.NFTTransaction) error {
	_, err := g.DB().Model("nft_transaction").Data(tx).Insert()
	return err
}

// UpdateNFTTransaction 更新NFT交易
func (d *ProtocolDao) UpdateNFTTransaction(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("nft_transaction").Where("id", id).Data(data).Update()
	return err
}

// CreateBridgeTransaction 创建跨链交易
func (d *ProtocolDao) CreateBridgeTransaction(ctx context.Context, tx *model.BridgeTransaction) error {
	_, err := g.DB().Model("bridge_transaction").Data(tx).Insert()
	return err
}

// UpdateBridgeTransaction 更新跨链交易
func (d *ProtocolDao) UpdateBridgeTransaction(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("bridge_transaction").Where("id", id).Data(data).Update()
	return err
}

// CreateAggregatorTransaction 创建聚合交易
func (d *ProtocolDao) CreateAggregatorTransaction(ctx context.Context, tx *model.AggregatorTransaction) error {
	_, err := g.DB().Model("aggregator_transaction").Data(tx).Insert()
	return err
}

// UpdateAggregatorTransaction 更新聚合交易
func (d *ProtocolDao) UpdateAggregatorTransaction(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("aggregator_transaction").Where("id", id).Data(data).Update()
	return err
}
