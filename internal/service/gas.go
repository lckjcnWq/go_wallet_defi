package service

import (
	"context"
	"encoding/json"
	"errors"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/pkg/ethclientx"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type GasService struct{}

var Gas = &GasService{}

// BatchTransactions 批量执行交易
func (s *GasService) BatchTransactions(ctx context.Context, chainId uint64, calls []model.Call) (string, error) {
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", err
	}
	//1.编码调整数据
	callsData, err := contracts.PackMulticallBatch(calls)
	if err != nil {
		return "", err
	}
	//2.估算gas
	gasPrice, err := s.EstimateGas(ctx, chainId, callsData)
	if err != nil {
		return "", err
	}
	//3.执行交易
	tx, err := client.SendTransaction(ctx, contracts.MulticallAddress, big.NewInt(0), callsData)
	if err != nil {
		return "", err
	}

	// 保存批量交易记录
	callsJson, _ := json.Marshal(calls)
	batchTx := &model.BatchTransaction{
		BatchId:  generateBatchId(),
		ChainId:  chainId,
		TxHash:   tx.Hash().Hex(),
		From:     from,
		Calls:    string(callsJson),
		Status:   0,
		GasPrice: gasPrice.String(),
	}

	err = dao.Gas.CreateBatchTransaction(ctx, batchTx)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil

}

// EstimateGasPrice 估算Gas价格
func (s *GasService) EstimateGasPrice(ctx context.Context, chainId uint64, strategyType string) (*big.Int, error) {
	// 获取最近的Gas价格记录
	prices, err := dao.Gas.GetRecentGasPrices(ctx, chainId, 10)
	if err != nil {
		return nil, err
	}

	// 获取当前的Gas策略
	strategy, err := dao.Gas.GetGasStrategy(ctx, chainId, strategyType)
	if err != nil {
		return nil, err
	}

	// 分析历史Gas价格趋势
	trend := s.analyzeGasTrend(prices)

	// 根据策略和趋势计算Gas价格
	baseFee, _ := new(big.Int).SetString(strategy.BaseFee, 10)
	priority, _ := new(big.Int).SetString(strategy.Priority, 10)

	gasPrice := new(big.Int).Add(baseFee, priority)

	// 根据趋势调整
	if trend > 0 {
		gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(11))
		gasPrice = gasPrice.Div(gasPrice, big.NewInt(10))
	} else if trend < 0 {
		gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(9))
		gasPrice = gasPrice.Div(gasPrice, big.NewInt(10))
	}

	return gasPrice, nil
}

// UpdateGasPrices 更新Gas价格
func (s *GasService) UpdateGasPrices(ctx context.Context, chainId uint64) error {
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return err
	}

	// 获取最新区块
	block, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		return err
	}

	// 获取基础费用
	baseFee := block.BaseFee()

	// 计算优先级费用
	maxPriorityFeePerGas, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}

	// 保存Gas价格记录
	price := &model.GasPrice{
		ChainId:     chainId,
		BlockNumber: block.NumberU64(),
		BaseFee:     baseFee.String(),
		Priority:    maxPriorityFeePerGas.String(),
		Average:     new(big.Int).Add(baseFee, maxPriorityFeePerGas).String(),
	}

	err = dao.Gas.CreateGasPrice(ctx, price)
	if err != nil {
		return err
	}

	// 更新Gas策略
	err = s.updateGasStrategies(ctx, chainId, baseFee, maxPriorityFeePerGas)
	if err != nil {
		return err
	}

	return nil
}

func (s *GasService) updateGasStrategies(ctx context.Context, chainId uint64, baseFee, priority *big.Int) error {
	//定义策略倍数
	strategies := map[string]struct {
		baseMul     int64
		priorityMul int64
		wait        int
	}{
		"FASTEST":  {12, 20, 15},
		"FAST":     {11, 15, 30},
		"STANDARD": {10, 10, 60},
		"SLOW":     {9, 5, 180},
	}

	for sType, s := range strategies {
		baseFeeMul := new(big.Int).Mul(baseFee, big.NewInt(s.baseMul))
		baseFeeFinal := baseFeeMul.Div(baseFeeMul, big.NewInt(10))

		priorityMul := new(big.Int).Mul(priority, big.NewInt(s.priorityMul))
		priorityFinal := priorityMul.Div(priorityMul, big.NewInt(10))

		strategy := &model.GasStrategy{
			ChainId:      chainId,
			Type:         sType,
			BaseFee:      baseFeeFinal.String(),
			Priority:     priorityFinal.String(),
			EstimateWait: s.wait,
		}

		err := dao.Gas.UpdateGasStrategy(ctx, strategy)
		if err != nil {
			return err
		}
	}

	return nil
}

// 检查三明治攻击
func (s *GasService) checkSandwichAttack(ctx context.Context, chainId uint64, tx *ethclient.Transaction) error {
	client, err := ethclient.GetClient(chainId)
	if err != nil {
		return err
	}

	// 获取当前区块的交易
	block, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		return err
	}

	txs := block.Transactions()
	for i, blockTx := range txs {
		// 如果发现可疑的前后交易对
		if i > 0 && i < len(txs)-1 {
			prevTx := txs[i-1]
			nextTx := txs[i+1]

			// 检查是否存在同一代币的买入卖出
			if s.isSandwichPattern(prevTx, blockTx, nextTx) {
				// 记录MEV防护
				protection := &model.MEVProtection{
					TxHash:      tx.Hash().Hex(),
					Type:        "SANDWICH",
					From:        prevTx.From().Hex(),
					To:          blockTx.To().Hex(),
					BlockNumber: block.NumberU64(),
					GasPrice:    blockTx.GasPrice().String(),
				}

				err = dao.Gas.CreateMEVProtection(ctx, protection)
				if err != nil {
					return err
				}

				return errors.New("potential sandwich attack detected")
			}
		}
	}

	return nil
}

// AccelerateTransaction 加速交易
func (s *GasService) AccelerateTransaction(ctx context.Context, chainId uint64, txHash string) (string, error) {
	client, err := ethclient.GetClient(chainId)
	if err != nil {
		return "", err
	}

	// 获取原始交易
	tx, isPending, err := client.TransactionByHash(ctx, common.HexToHash(txHash))
	if err != nil {
		return "", err
	}

	if !isPending {
		return "", errors.New("transaction is already confirmed")
	}

	// 计算新的gas价格(原价格的1.5倍)
	oldGasPrice := tx.GasPrice()
	newGasPrice := new(big.Int).Mul(oldGasPrice, big.NewInt(15))
	newGasPrice = newGasPrice.Div(newGasPrice, big.NewInt(10))

	// 创建新交易
	nonce, err := client.PendingNonceAt(ctx, tx.From())
	if err != nil {
		return "", err
	}

	newTx := types.NewTransaction(
		nonce,
		*tx.To(),
		tx.Value(),
		tx.Gas(),
		newGasPrice,
		tx.Data(),
	)

	// 发送新交易
	signedTx, err := client.SignAndSendTransaction(ctx, newTx)
	if err != nil {
		return "", err
	}

	// 记录加速记录
	acceleration := &model.TxAcceleration{
		TxHash:      txHash,
		NewTxHash:   signedTx.Hash().Hex(),
		ChainId:     chainId,
		From:        tx.From().Hex(),
		OldGasPrice: oldGasPrice.String(),
		NewGasPrice: newGasPrice.String(),
		Status:      0,
	}

	err = dao.Gas.CreateTxAcceleration(ctx, acceleration)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// CancelTransaction 取消交易
func (s *GasService) CancelTransaction(ctx context.Context, chainId uint64, txHash string) (string, error) {
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", err
	}

	// 获取原始交易
	tx, isPending, err := client.TransactionByHash(ctx, common.HexToHash(txHash))
	if err != nil {
		return "", err
	}

	if !isPending {
		return "", errors.New("transaction is already confirmed")
	}

	// 创建0值交易发送到自己的地址(使用更高的gas价格)
	nonce := tx.Nonce()
	gasPrice := new(big.Int).Mul(tx.GasPrice(), big.NewInt(2)) // 使用2倍gas价格

	newTx := types.NewTransaction(
		nonce,
		tx.From(),
		big.NewInt(0),
		21000, // 基本转账gas限制
		gasPrice,
		nil,
	)

	// 发送取消交易
	signedTx, err := client.SignAndSendTransaction(ctx, newTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// GetGasInfo 获取Gas信息
func (s *GasService) GetGasInfo(ctx context.Context, chainId uint64) (map[string]interface{}, error) {
	// 获取当前Gas价格
	prices, err := dao.Gas.GetRecentGasPrices(ctx, chainId, 1)
	if err != nil {
		return nil, err
	}

	// 获取所有Gas策略
	strategies := make(map[string]*model.GasStrategy)
	for _, sType := range []string{"FASTEST", "FAST", "STANDARD", "SLOW"} {
		strategy, err := dao.Gas.GetGasStrategy(ctx, chainId, sType)
		if err != nil {
			return nil, err
		}
		strategies[sType] = strategy
	}

	return map[string]interface{}{
		"current_price": prices[0],
		"strategies":    strategies,
	}, nil
}
