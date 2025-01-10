package dao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/model"
)

type SecurityDao struct{}

var Security = &SecurityDao{}

// CreateMultiSigWallet 创建多签钱包
func (d *SecurityDao) CreateMultiSigWallet(ctx context.Context, wallet *model.MultiSigWallet) error {
	_, err := g.DB().Model("multi_sig_wallet").Ctx(ctx).Data(wallet).Insert()
	return err
}

// GetMultiSigWallet 获取多签钱包
func (d *SecurityDao) GetMultiSigWallet(ctx context.Context, id uint64) (*model.MultiSigWallet, error) {
	var wallet *model.MultiSigWallet
	err := g.DB().Model("multi_sig_wallet").Ctx(ctx).Where("id", id).Scan(&wallet)
	return wallet, err
}

// UpdateMultiSigWallet 更新多签钱包
func (d *SecurityDao) UpdateMultiSigWallet(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("multi_sig_wallet").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// CreateMultiSigTx 创建多签交易
func (d *SecurityDao) CreateMultiSigTx(ctx context.Context, tx *model.MultiSigTransaction) error {
	_, err := g.DB().Model("multi_sig_transaction").Ctx(ctx).Data(tx).Insert()
	return err
}

// GetMultiSigTx 获取多签交易
func (d *SecurityDao) GetMultiSigTx(ctx context.Context, id uint64) (*model.MultiSigTransaction, error) {
	var tx *model.MultiSigTransaction
	err := g.DB().Model("multi_sig_transaction").Ctx(ctx).Where("id", id).Scan(&tx)
	return tx, err
}

// UpdateMultiSigTx 更新多签交易
func (d *SecurityDao) UpdateMultiSigTx(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("multi_sig_transaction").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// CreateWhitelist 创建白名单
func (d *SecurityDao) CreateWhitelist(ctx context.Context, whitelist *model.Whitelist) error {
	_, err := g.DB().Model("whitelist").Ctx(ctx).Data(whitelist).Insert()
	return err
}

// GetWhitelist 获取白名单
func (d *SecurityDao) GetWhitelist(ctx context.Context, id uint64) (*model.Whitelist, error) {
	var whitelist *model.Whitelist
	err := g.DB().Model("whitelist").Ctx(ctx).Where("id", id).Scan(&whitelist)
	return whitelist, err
}

// UpdateWhitelist 更新白名单
func (d *SecurityDao) UpdateWhitelist(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("whitelist").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// CreateTransactionLimit 创建交易限额
func (d *SecurityDao) CreateTransactionLimit(ctx context.Context, limit *model.TransactionLimit) error {
	_, err := g.DB().Model("transaction_limit").Ctx(ctx).Data(limit).Insert()
	return err
}

// GetTransactionLimit 获取交易限额
func (d *SecurityDao) GetTransactionLimit(ctx context.Context, id uint64) (*model.TransactionLimit, error) {
	var limit *model.TransactionLimit
	err := g.DB().Model("transaction_limit").Ctx(ctx).Where("id", id).Scan(&limit)
	return limit, err
}

// UpdateTransactionLimit 更新交易限额
func (d *SecurityDao) UpdateTransactionLimit(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("transaction_limit").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// CreateRiskRule 创建风控规则
func (d *SecurityDao) CreateRiskRule(ctx context.Context, rule *model.RiskRule) error {
	_, err := g.DB().Model("risk_rule").Ctx(ctx).Data(rule).Insert()
	return err
}

// GetRiskRule 获取风控规则
func (d *SecurityDao) GetRiskRule(ctx context.Context, id uint64) (*model.RiskRule, error) {
	var rule *model.RiskRule
	err := g.DB().Model("risk_rule").Ctx(ctx).Where("id", id).Scan(&rule)
	return rule, err
}

// UpdateRiskRule 更新风控规则
func (d *SecurityDao) UpdateRiskRule(ctx context.Context, id uint64, data g.Map) error {
	_, err := g.DB().Model("risk_rule").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

// CreateRiskLog 创建风控日志
func (d *SecurityDao) CreateRiskLog(ctx context.Context, log *model.RiskLog) error {
	_, err := g.DB().Model("risk_log").Ctx(ctx).Data(log).Insert()
	return err
}

// GetUserTxAmount 获取用户交易金额统计
func (d *SecurityDao) GetUserTxAmount(ctx context.Context, userId uint64, tokenAddress string, startTime, endTime int64) (string, error) {
	var amount struct{ Total string }
	err := g.DB().Model("transaction").
		Ctx(ctx).
		Where("user_id", userId).
		Where("token_address", tokenAddress).
		Where("created_at between ? and ?", startTime, endTime).
		Fields("sum(amount) as total").
		Scan(&amount)
	return amount.Total, err
}

// CheckWhitelist 检查白名单
func (d *SecurityDao) CheckWhitelist(ctx context.Context, userId uint64, address string, txType int) (bool, error) {
	count, err := g.DB().Model("whitelist").
		Ctx(ctx).
		Where("user_id", userId).
		Where("address", address).
		Where("type", txType).
		Where("status", 1).
		Where("expire_at > now() or expire_at is null").
		Count()
	return count > 0, err
}
