package service

import (
	"context"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/ethclientx"
	"math/big"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

type AnalysisService struct{}

var Analysis = &AnalysisService{}

// AnalyzeTransaction 分析交易
func (s *AnalysisService) AnalyzeTransaction(ctx context.Context, chainId uint64, address string) error {
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return err
	}
	//1.获取最新区块
	block, err := client.BlockNumber(ctx)
	if err != nil {
		return err
	}
	//2.获取历史交易
	txs, err := dao.Transaction.GetAddressTransactions(ctx, address, block)
	if err != nil {
		return err
	}
	//3.分析交易
	var totalVolume = new(big.Int)
	var totalGas = new(big.Int)

	var firstTxTime *gtime.Time
	var lastTxTime *gtime.Time

	for _, tx := range txs {
		// 计算总交易量
		amount, _ := new(big.Int).SetString(tx.Amount, 10)
		totalVolume = totalVolume.Add(totalVolume, amount)

		// 计算总gas费用
		gasUsed, _ := new(big.Int).SetString(tx.GasUsed, 10)
		gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
		gasFee := new(big.Int).Mul(gasUsed, gasPrice)
		totalGas = totalGas.Add(totalGas, gasFee)

		// 记录首次和最后交易时间
		txTime := tx.CreatedAt.Time
		if firstTxTime == nil || txTime.Before(*firstTxTime) {
			firstTxTime = &txTime
		}
		if lastTxTime == nil || txTime.After(*lastTxTime) {
			lastTxTime = &txTime
		}
	}

	// 更新交易分析
	analysis := &model.TransactionAnalysis{
		Address:      address,
		TotalTxCount: int64(len(txs)),
		TotalVolume:  totalVolume.String(),
		AvgGasUsed:   totalGas.Div(totalGas, big.NewInt(int64(len(txs)))).Int64(),
		TotalGasFee:  totalGas.String(),
		FirstTxTime:  gtime.New(firstTxTime),
		LastTxTime:   gtime.New(lastTxTime),
	}

	err = dao.Analysis.UpdateTxAnalysis(ctx, analysis)
	if err != nil {
		return err
	}

	// 分析交易模式
	err = s.analyzeTxPattern(ctx, address, txs)
	if err != nil {
		return err
	}

	return nil

}

// 分析交易模式
func (s *AnalysisService) analyzeTxPattern(ctx context.Context, address string, txs []*model.Transaction) error {
	// 分析周期性交易
	periodicPattern, err := s.analyzePeriodicPattern(txs)
	if err != nil {
		return err
	}
	if periodicPattern != nil {
		err = dao.Analysis.UpdateTxPattern(ctx, periodicPattern)
		if err != nil {
			return err
		}
	}

	// 分析大额交易
	largePattern, err := s.analyzeLargePattern(txs)
	if err != nil {
		return err
	}
	if largePattern != nil {
		err = dao.Analysis.UpdateTxPattern(ctx, largePattern)
		if err != nil {
			return err
		}
	}

	return nil
}

// AnalyzeAssets 分析资产
func (s *AnalysisService) AnalyzeAssets(ctx context.Context, chainId uint64, address string) error {
	// 获取ERC20代币资产
	tokens, err := dao.Token.GetUserTokens(ctx, address)
	if err != nil {
		return err
	}

	// 获取NFT资产
	nfts, err := dao.NFT.GetUserNFTs(ctx, address)
	if err != nil {
		return err
	}

	// 获取市场价格
	prices, err := s.getTokenPrices(ctx)
	if err != nil {
		return err
	}

	var totalValue = new(big.Int)

	// 分析ERC20资产
	for _, token := range tokens {
		balance, _ := new(big.Int).SetString(token.Balance, 10)
		price, exists := prices[token.TokenAddress]
		if !exists {
			continue
		}

		value := new(big.Int).Mul(balance, price)

		analysis := &model.AssetAnalysis{
			Address:      address,
			TokenAddress: token.TokenAddress,
			TokenType:    "ERC20",
			Balance:      token.Balance,
			Value:        value.String(),
			CostBasis:    token.CostBasis,
		}

		err = dao.Analysis.UpdateAssetAnalysis(ctx, analysis)
		if err != nil {
			return err
		}

		totalValue = totalValue.Add(totalValue, value)
	}

	// 分析NFT资产
	for _, nft := range nfts {
		floorPrice, exists := prices[nft.ContractAddress]
		if !exists {
			continue
		}

		value := new(big.Int).Mul(big.NewInt(1), floorPrice) // Using floor price

		analysis := &model.AssetAnalysis{
			Address:      address,
			TokenAddress: nft.ContractAddress,
			TokenType:    "ERC721",
			Balance:      "1",
			Value:        value.String(),
			CostBasis:    nft.CostBasis,
		}

		err = dao.Analysis.UpdateAssetAnalysis(ctx, analysis)
		if err != nil {
			return err
		}

		totalValue = totalValue.Add(totalValue, value)
	}

	// 更新资产分布
	err = s.updateAssetDistribution(ctx, address, totalValue, tokens, nfts)
	if err != nil {
		return err
	}

	return nil
}

// AnalyzeProfit 分析收益
func (s *AnalysisService) AnalyzeProfit(ctx context.Context, address string) error {
	// 分析DeFi收益
	err := s.analyzeDeFiProfit(ctx, address)
	if err != nil {
		return err
	}

	// 分析Staking收益
	err = s.analyzeStakingProfit(ctx, address)
	if err != nil {
		return err
	}

	// 分析NFT收益
	err = s.analyzeNFTProfit(ctx, address)
	if err != nil {
		return err
	}

	return nil
}

// 分析DeFi收益
func (s *AnalysisService) analyzeDeFiProfit(ctx context.Context, address string) error {
	// 获取用户的DeFi仓位
	positions, err := dao.DeFi.GetUserPositions(ctx, address)
	if err != nil {
		return err
	}

	for _, pos := range positions {
		// 计算收益
		principal, _ := new(big.Int).SetString(pos.Principal, 10)
		current, _ := new(big.Int).SetString(pos.Current, 10)
		profit := new(big.Int).Sub(current, principal)

		// 计算APY
		duration := time.Since(pos.StartTime.Time)
		apy := float64(profit.Int64()) / float64(principal.Int64()) * 365 * 24 * 3600 / duration.Seconds()

		analysis := &model.ProfitAnalysis{
			Address:   address,
			Type:      "DEFI",
			Platform:  pos.Platform,
			Principal: pos.Principal,
			Profit:    profit.String(),
			Apy:       apy,
			StartTime: pos.StartTime,
		}

		err = dao.Analysis.UpdateProfitAnalysis(ctx, analysis)
		if err != nil {
			return err
		}
	}

	return nil
}

// AnalyzeRisk 分析风险
func (s *AnalysisService) AnalyzeRisk(ctx context.Context, address string) error {
	// 分析交易风险
	err := s.analyzeTxRisk(ctx, address)
	if err != nil {
		return err
	}

	// 分析地址风险
	err = s.analyzeAddressRisk(ctx, address)
	if err != nil {
		return err
	}

	// 分析资产风险
	err = s.analyzeAssetRisk(ctx, address)
	if err != nil {
		return err
	}

	// 分析市场风险
	err = s.analyzeMarketRisk(ctx, address)
	if err != nil {
		return err
	}

	return nil
}

// 分析交易风险
func (s *AnalysisService) analyzeTxRisk(ctx context.Context, address string) error {
	// 获取近期交易
	txs, err := dao.Transaction.GetRecentTransactions(ctx, address, 100)
	if err != nil {
		return err
	}

	// 检查大额交易
	largeThreshold := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100)) // 100 ETH
	for _, tx := range txs {
		amount, _ := new(big.Int).SetString(tx.Amount, 10)
		if amount.Cmp(largeThreshold) > 0 {
			risk := &model.RiskAnalysis{
				Address:     address,
				Type:        "TX",
				RiskLevel:   3,
				RiskFactor:  "LARGE_AMOUNT",
				Description: "Large amount transaction detected",
			}
			err = dao.Analysis.CreateRiskAnalysis(ctx, risk)
			if err != nil {
				return err
			}
		}
	}

	// 检查频繁交易
	if len(txs) > 50 { // 如果24小时内交易超过50笔
		risk := &model.RiskAnalysis{
			Address:     address,
			Type:        "TX",
			RiskLevel:   2,
			RiskFactor:  "HIGH_FREQUENCY",
			Description: "High frequency trading detected",
		}
		err = dao.Analysis.CreateRiskAnalysis(ctx, risk)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateMarketData 更新市场数据
func (s *AnalysisService) UpdateMarketData(ctx context.Context) error {
	// 获取市场数据
	indicators, err := s.getMarketData(ctx)
	if err != nil {
		return err
	}

	for _, indicator := range indicators {
		err = dao.Analysis.UpdateMarketIndicator(ctx, indicator)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAnalysisReport 获取分析报告
func (s *AnalysisService) GetAnalysisReport(ctx context.Context, address string) (map[string]interface{}, error) {
	// 获取交易分析
	txAnalysis, err := dao.Analysis.GetTxAnalysis(ctx, address)
	if err != nil {
		return nil, err
	}

	// 获取资产分析
	assetAnalysis, err := dao.Analysis.GetAssetAnalysis(ctx, address)
	if err != nil {
		return nil, err
	}

	// 获取收益分析
	profitAnalysis, err := dao.Analysis.GetProfitAnalysis(ctx, address)
	if err != nil {
		return nil, err
	}

	// 获取风险分析
	riskAnalysis, err := dao.Analysis.GetRiskAnalysis(ctx, address)
	if err != nil {
		return nil, err
	}

	// 获取交易模式
	txPatterns, err := dao.Analysis.GetTxPatterns(ctx, address)
	if err != nil {
		return nil, err
	}

	// 获取资产分布
	distribution, err := dao.Analysis.GetAssetDistribution(ctx, address)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"transaction_analysis": txAnalysis,
		"asset_analysis":       assetAnalysis,
		"profit_analysis":      profitAnalysis,
		"risk_analysis":        riskAnalysis,
		"tx_patterns":          txPatterns,
		"asset_distribution":   distribution,
	}, nil
}
