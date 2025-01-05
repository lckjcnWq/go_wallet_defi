package dao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/model"
	"time"
)

type NFTDao struct{}

var NFT = &NFTDao{}

// Insert 插入NFT
func (d *NFTDao) Insert(ctx context.Context, nft *model.NFT) error {
	_, err := g.DB().Model("nft").Ctx(ctx).Data(nft).Insert()
	return err
}

// UpdateOwner 更新NFT所有者
func (d *NFTDao) UpdateOwner(ctx context.Context, id uint64, owner string) error {
	_, err := g.DB().Model("nft").Ctx(ctx).
		Where("id", id).
		Data(g.Map{
			"owner":      owner,
			"updated_at": g.NewVar(time.Now().Unix()),
		}).
		Update()
	return err
}

// GetById 根据ID获取NFT
func (d *NFTDao) GetById(ctx context.Context, id uint64) (*model.NFT, error) {
	var nft *model.NFT
	err := g.DB().Model("nft").Ctx(ctx).Where("id", id).Scan(&nft)
	return nft, err
}

// GetList 获取NFT列表
func (d *NFTDao) GetList(ctx context.Context, owner, creator string, contractId uint64, page, pageSize int) ([]*model.NFT, int, error) {
	m := g.DB().Model("nft").Where("status", 1)

	if owner != "" {
		m = m.Where("owner", owner)
	}
	if creator != "" {
		m = m.Where("creator", creator)
	}
	if contractId > 0 {
		m = m.Where("contract_id", contractId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []*model.NFT
	err = m.Page(page, pageSize).Order("id DESC").Scan(&list)

	return list, total, err
}

// InsertTransfer 插入转移记录
func (d *NFTDao) InsertTransfer(ctx context.Context, transfer *model.NFTTransfer) error {
	_, err := g.DB().Model("nft_transfer").Ctx(ctx).Data(transfer).Insert()
	return err
}

// InsertMarket 插入市场记录
func (d *NFTDao) InsertMarket(ctx context.Context, market *model.NFTMarket) error {
	_, err := g.DB().Model("nft_market").Ctx(ctx).Data(market).Insert()
	return err
}

// UpdateMarket 更新市场记录
func (d *NFTDao) UpdateMarket(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("nft_market").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// GetMarketList 获取市场列表
func (d *NFTDao) GetMarketList(ctx context.Context, seller string, status int, page, pageSize int) ([]*model.NFTMarket, int, error) {
	m := g.DB().Model("nft_market")

	if seller != "" {
		m = m.Where("seller", seller)
	}
	if status > 0 {
		m = m.Where("status", status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []*model.NFTMarket
	err = m.Page(page, pageSize).Order("id DESC").Scan(&list)

	return list, total, err
}

// CountByContract 获取合约NFT数量
func (d *NFTDao) CountByContract(ctx context.Context, contractId uint64) (int, error) {
	return g.DB().Model("nft").Where("contract_id", contractId).Count()
}
