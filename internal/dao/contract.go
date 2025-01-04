package dao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/model"
)

type ContractDao struct{}

var Contract = &ContractDao{}

// Insert 插入合约信息
func (d *ContractDao) Insert(ctx context.Context, contract *model.Contract) error {
	_, err := g.DB().Model("contract").Ctx(ctx).Data(contract).Insert()
	return err
}

// Update 更新合约信息
func (d *ContractDao) Update(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("contract").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// GetByAddress 根据地址获取合约
func (d *ContractDao) GetByAddress(ctx context.Context, address, network string) (*model.Contract, error) {
	var contract *model.Contract
	err := g.DB().Model("contract").
		Ctx(ctx).
		Where("address", address).
		Where("network", network).
		Where("status", 1).
		Scan(&contract)
	return contract, err
}

// InsertEvent 插入合约事件
func (d *ContractDao) InsertEvent(ctx context.Context, event *model.ContractEvent) error {
	_, err := g.DB().Model("contract_event").Ctx(ctx).Data(event).Insert()
	return err
}

// GetEvents 获取合约事件
func (d *ContractDao) GetEvents(ctx context.Context, contractId uint64, eventName string, fromBlock, toBlock int64, page, pageSize int) ([]*model.ContractEvent, int, error) {
	m := g.DB().Model("contract_event").
		Where("contract_id", contractId)

	if eventName != "" {
		m = m.Where("name", eventName)
	}
	if fromBlock > 0 {
		m = m.Where("block_number >= ?", fromBlock)
	}
	if toBlock > 0 {
		m = m.Where("block_number <= ?", toBlock)
	}

	total, err := m.Ctx(ctx).Count()
	if err != nil {
		return nil, 0, err
	}

	var list []*model.ContractEvent
	err = m.Ctx(ctx).Page(page, pageSize).Order("id DESC").Scan(&list)

	return list, total, err
}

// InsertCall 插入合约调用
func (d *ContractDao) InsertCall(ctx context.Context, call *model.ContractCall) error {
	_, err := g.DB().Model("contract_call").Ctx(ctx).Data(call).Insert()
	return err
}

// UpdateCall 更新合约调用
func (d *ContractDao) UpdateCall(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("contract_call").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}
