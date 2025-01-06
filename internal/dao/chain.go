package dao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/model"
)

type ChainDao struct{}

var Chain = &ChainDao{}

// GetById 根据ID获取链信息
func (d *ChainDao) GetById(ctx context.Context, id uint64) (*model.Chain, error) {
	var chain *model.Chain
	err := g.DB().Model("chain").Where("id", id).Scan(&chain)
	return chain, err
}

// GetByChainId 根据ID获取链信息
func (d *ChainDao) GetByChainId(ctx context.Context, chainId uint64) (*model.Chain, error) {
	var chain *model.Chain
	err := g.DB().Model("chain").Where("chainId", chainId).Scan(&chain)
	return chain, err
}

// GetMapping 获取合约地址映射
func (d *ChainDao) GetMapping(ctx context.Context, fromChainId uint64, fromAddress string, toChainId uint64) (*model.ContractMapping, error) {
	var mapping *model.ContractMapping
	err := g.DB().Model("contract_mapping").
		Where("from_chain_id", fromChainId).
		Where("from_address", fromAddress).
		Where("to_chain_id", toChainId).
		Scan(&mapping)
	return mapping, err
}

// InsertCrossTransfer 插入跨链交易
func (d *ChainDao) InsertCrossTransfer(ctx context.Context, transfer *model.CrossTransfer) error {
	_, err := g.DB().Model("cross_transfer").Data(transfer).Insert()
	return err
}

// UpdateCrossTransfer 更新跨链交易
func (d *ChainDao) UpdateCrossTransfer(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("cross_transfer").Data(data).Where("id", id).Update()
	return err
}

// GetCrossTransferList 获取跨链交易列表
func (d *ChainDao) GetCrossTransferList(ctx context.Context, fromChainId, toChainId uint64, address string, status int, page, pageSize int) ([]*model.CrossTransfer, int, error) {
	m := g.DB().Model("cross_transfer")

	if fromChainId > 0 {
		m = m.Where("from_chain_id", fromChainId)
	}
	if toChainId > 0 {
		m = m.Where("to_chain_id", toChainId)
	}
	if address != "" {
		m = m.Where("from_address = ? OR to_address = ?", address, address)
	}
	if status > 0 {
		m = m.Where("status", status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []*model.CrossTransfer
	err = m.Page(page, pageSize).Order("id DESC").Scan(&list)

	return list, total, err
}
