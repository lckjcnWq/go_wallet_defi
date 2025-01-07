package dao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/model"
)

type DefiDao struct{}

var Defi = &DefiDao{}

// InsertDexTrade 插入DEX交易记录
func (d *DefiDao) InsertDexTrade(ctx context.Context, trade *model.DexTrade) error {
	_, err := g.DB().Model("dex_trade").Ctx(ctx).Data(trade).Insert()
	return err
}

// UpdateDexTrade 更新DEX交易记录
func (d *DefiDao) UpdateDexTrade(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("dex_trade").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// GetDexTradeByHash 根据交易哈希获取DEX交易记录
func (d *DefiDao) GetDexTradeByHash(ctx context.Context, hash string) (*model.DexTrade, error) {
	var trade *model.DexTrade
	err := g.DB().Model("dex_trade").Ctx(ctx).Where("hash", hash).Scan(&trade)
	return trade, err
}

// InsertLiquidity 插入流动性记录
func (d *DefiDao) InsertLiquidity(ctx context.Context, liquidity *model.Liquidity) error {
	_, err := g.DB().Model("liquidity").Ctx(ctx).Data(liquidity).Insert()
	return err
}

// UpdateLiquidity 更新流动性记录
func (d *DefiDao) UpdateLiquidity(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("liquidity").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// GetLiquidityByHash 根据交易哈希获取流动性记录
func (d *DefiDao) GetLiquidityByHash(ctx context.Context, hash string) (*model.Liquidity, error) {
	var liquidity *model.Liquidity
	err := g.DB().Model("liquidity").Ctx(ctx).Where("hash", hash).Scan(&liquidity)
	return liquidity, err
}

// InsertLending 插入借贷记录
func (d *DefiDao) InsertLending(ctx context.Context, lending *model.Lending) error {
	_, err := g.DB().Model("lending").Ctx(ctx).Data(lending).Insert()
	return err
}

// UpdateLending 更新借贷记录
func (d *DefiDao) UpdateLending(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("lending").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// GetLendingByHash 根据交易哈希获取借贷记录
func (d *DefiDao) GetLendingByHash(ctx context.Context, hash string) (*model.Lending, error) {
	var lending *model.Lending
	err := g.DB().Model("lending").Ctx(ctx).Where("hash", hash).Scan(&lending)
	return lending, err
}

// InsertYieldFarm 插入收益农场记录
func (d *DefiDao) InsertYieldFarm(ctx context.Context, farm *model.YieldFarm) error {
	_, err := g.DB().Model("yield_farm").Ctx(ctx).Data(farm).Insert()
	return err
}

// UpdateYieldFarm 更新收益农场记录
func (d *DefiDao) UpdateYieldFarm(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("yield_farm").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// GetYieldFarmByHash 根据交易哈希获取收益农场记录
func (d *DefiDao) GetYieldFarmByHash(ctx context.Context, hash string) (*model.YieldFarm, error) {
	var farm *model.YieldFarm
	err := g.DB().Model("yield_farm").Ctx(ctx).Where("hash", hash).Scan(&farm)
	return farm, err
}

// InsertVault 插入机枪池记录
func (d *DefiDao) InsertVault(ctx context.Context, vault *model.Vault) error {
	_, err := g.DB().Model("vault").Ctx(ctx).Data(vault).Insert()
	return err
}

// UpdateVault 更新机枪池记录
func (d *DefiDao) UpdateVault(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("vault").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// GetVaultByHash 根据交易哈希获取机枪池记录
func (d *DefiDao) GetVaultByHash(ctx context.Context, hash string) (*model.Vault, error) {
	var vault *model.Vault
	err := g.DB().Model("vault").Ctx(ctx).Where("hash", hash).Scan(&vault)
	return vault, err
}

// GetUserDexTrades 获取用户DEX交易记录
func (d *DefiDao) GetUserDexTrades(ctx context.Context, user string, page, pageSize int) ([]*model.DexTrade, int, error) {
	m := g.DB().Model("dex_trade").Where("user", user)

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []*model.DexTrade
	err = m.Page(page, pageSize).Order("id DESC").Scan(&list)

	return list, total, err
}

// GetUserLending 获取用户借贷记录
func (d *DefiDao) GetUserLending(ctx context.Context, user string, page, pageSize int) ([]*model.Lending, int, error) {
	m := g.DB().Model("lending").Where("user", user)

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []*model.Lending
	err = m.Page(page, pageSize).Order("id DESC").Scan(&list)

	return list, total, err
}

// GetUserYieldFarm 获取用户收益农场记录
func (d *DefiDao) GetUserYieldFarm(ctx context.Context, user string, page, pageSize int) ([]*model.YieldFarm, int, error) {
	m := g.DB().Model("yield_farm").Where("user", user)

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []*model.YieldFarm
	err = m.Page(page, pageSize).Order("id DESC").Scan(&list)

	return list, total, err
}

// GetUserVault 获取用户机枪池记录
func (d *DefiDao) GetUserVault(ctx context.Context, user string, page, pageSize int) ([]*model.Vault, int, error) {
	m := g.DB().Model("vault").Where("user", user)

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []*model.Vault
	err = m.Page(page, pageSize).Order("id DESC").Scan(&list)

	return list, total, err
}
