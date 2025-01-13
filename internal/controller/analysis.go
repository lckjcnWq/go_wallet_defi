package controller

import (
	"context"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/service"
)

type AnalysisController struct{}

// AnalyzeAddress 分析地址
func (c *AnalysisController) AnalyzeAddress(ctx context.Context, req *v1.AnalyzeAddressReq) (res *v1.AnalyzeAddressRes, err error) {
	// 分析交易
	err = service.Analysis.AnalyzeTransaction(ctx, req.ChainId, req.Address)
	if err != nil {
		return nil, err
	}

	// 分析资产
	err = service.Analysis.AnalyzeAssets(ctx, req.ChainId, req.Address)
	if err != nil {
		return nil, err
	}

	// 分析收益
	err = service.Analysis.AnalyzeProfit(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	// 分析风险
	err = service.Analysis.AnalyzeRisk(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	return &v1.AnalyzeAddressRes{}, nil
}

// GetAnalysisReport 获取分析报告
func (c *AnalysisController) GetAnalysisReport(ctx context.Context, req *v1.GetAnalysisReportReq) (res *v1.GetAnalysisReportRes, err error) {
	report, err := service.Analysis.GetAnalysisReport(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	return &v1.GetAnalysisReportRes{
		Report: report,
	}, nil
}

// UpdateMarketData 更新市场数据
func (c *AnalysisController) UpdateMarketData(ctx context.Context, req *v1.UpdateMarketDataReq) (res *v1.UpdateMarketDataRes, err error) {
	err = service.Analysis.UpdateMarketData(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateMarketDataRes{}, nil
}
