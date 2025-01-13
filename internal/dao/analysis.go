package dao

import (
	"context"
	"go-wallet-defi/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type AnalysisDao struct{}

var Analysis = &AnalysisDao{}

// UpdateTxAnalysis 更新交易分析
func (d *AnalysisDao) UpdateTxAnalysis(ctx context.Context, analysis *model.TransactionAnalysis) error {
	_, err := g.DB().Model("transaction_analysis").
		Where("address", analysis.Address).
		Data(analysis).
		Save()
	return err
}

// GetTxAnalysis 获取交易分析
func (d *AnalysisDao) GetTxAnalysis(ctx context.Context, address string) (*model.TransactionAnalysis, error) {
	var analysis *model.TransactionAnalysis
	err := g.DB().Model("transaction_analysis").
		Where("address", address).
		Scan(&analysis)
	return analysis, err
}

// UpdateAssetAnalysis 更新资产分析
func (d *AnalysisDao) UpdateAssetAnalysis(ctx context.Context, analysis *model.AssetAnalysis) error {
	_, err := g.DB().Model("asset_analysis").
		Where("address", analysis.Address).
		Where("token_address", analysis.TokenAddress).
		Data(analysis).
		Save()
	return err
}

// GetAssetAnalysis 获取资产分析
func (d *AnalysisDao) GetAssetAnalysis(ctx context.Context, address string) ([]*model.AssetAnalysis, error) {
	var analysis []*model.AssetAnalysis
	err := g.DB().Model("asset_analysis").
		Where("address", address).
		Scan(&analysis)
	return analysis, err
}

// UpdateProfitAnalysis 更新收益分析
func (d *AnalysisDao) UpdateProfitAnalysis(ctx context.Context, analysis *model.ProfitAnalysis) error {
	_, err := g.DB().Model("profit_analysis").
		Where("address", analysis.Address).
		Where("type", analysis.Type).
		Where("platform", analysis.Platform).
		Data(analysis).
		Save()
	return err
}

// GetProfitAnalysis 获取收益分析
func (d *AnalysisDao) GetProfitAnalysis(ctx context.Context, address string) ([]*model.ProfitAnalysis, error) {
	var analysis []*model.ProfitAnalysis
	err := g.DB().Model("profit_analysis").
		Where("address", address).
		Scan(&analysis)
	return analysis, err
}

// CreateRiskAnalysis 创建风险分析
func (d *AnalysisDao) CreateRiskAnalysis(ctx context.Context, analysis *model.RiskAnalysis) error {
	_, err := g.DB().Model("risk_analysis").
		Data(analysis).
		Insert()
	return err
}

// GetRiskAnalysis 获取风险分析
func (d *AnalysisDao) GetRiskAnalysis(ctx context.Context, address string) ([]*model.RiskAnalysis, error) {
	var analysis []*model.RiskAnalysis
	err := g.DB().Model("risk_analysis").
		Where("address", address).
		Order("created_at DESC").
		Scan(&analysis)
	return analysis, err
}

// UpdateTxPattern 更新交易模式
func (d *AnalysisDao) UpdateTxPattern(ctx context.Context, pattern *model.TxPattern) error {
	_, err := g.DB().Model("tx_pattern").
		Where("address", pattern.Address).
		Where("pattern", pattern.Pattern).
		Data(pattern).
		Save()
	return err
}

// GetTxPatterns 获取交易模式
func (d *AnalysisDao) GetTxPatterns(ctx context.Context, address string) ([]*model.TxPattern, error) {
	var patterns []*model.TxPattern
	err := g.DB().Model("tx_pattern").
		Where("address", address).
		Order("frequency DESC").
		Scan(&patterns)
	return patterns, err
}

// UpdateMarketIndicator 更新市场指标
func (d *AnalysisDao) UpdateMarketIndicator(ctx context.Context, indicator *model.MarketIndicator) error {
	_, err := g.DB().Model("market_indicator").
		Where("symbol", indicator.Symbol).
		Data(indicator).
		Save()
	return err
}

// GetMarketIndicators 获取市场指标
func (d *AnalysisDao) GetMarketIndicators(ctx context.Context) ([]*model.MarketIndicator, error) {
	var indicators []*model.MarketIndicator
	err := g.DB().Model("market_indicator").
		Order("volume_24h DESC").
		Scan(&indicators)
	return indicators, err
}

// UpdateAssetDistribution 更新资产分布
func (d *AnalysisDao) UpdateAssetDistribution(ctx context.Context, distribution *model.AssetDistribution) error {
	_, err := g.DB().Model("asset_distribution").
		Where("address", distribution.Address).
		Where("asset_type", distribution.AssetType).
		Data(distribution).
		Save()
	return err
}

// GetAssetDistribution 获取资产分布
func (d *AnalysisDao) GetAssetDistribution(ctx context.Context, address string) ([]*model.AssetDistribution, error) {
	var distribution []*model.AssetDistribution
	err := g.DB().Model("asset_distribution").
		Where("address", address).
		Order("percentage DESC").
		Scan(&distribution)
	return distribution, err
}
