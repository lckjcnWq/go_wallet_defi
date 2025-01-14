package controller

import (
	"context"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/service"
)

type GasController struct{}

// BatchTransactions 批量执行交易
func (c *GasController) BatchTransactions(ctx context.Context, req *v1.BatchTransactionsReq) (res *v1.BatchTransactionsRes, err error) {
	txHash, err := service.Gas.BatchTransactions(ctx, req.ChainId, req.From, req.Calls)
	if err != nil {
		return nil, err
	}

	return &v1.BatchTransactionsRes{
		TxHash: txHash,
	}, nil
}

// EstimateGasPrice 估算Gas价格
func (c *GasController) EstimateGasPrice(ctx context.Context, req *v1.EstimateGasPriceReq) (res *v1.EstimateGasPriceRes, err error) {
	gasPrice, err := service.Gas.EstimateGasPrice(ctx, req.ChainId, req.Strategy)
	if err != nil {
		return nil, err
	}

	return &v1.EstimateGasPriceRes{
		GasPrice: gasPrice.String(),
	}, nil
}

// AccelerateTransaction 加速交易
func (c *GasController) AccelerateTransaction(ctx context.Context, req *v1.AccelerateTransactionReq) (res *v1.AccelerateTransactionRes, err error) {
	newTxHash, err := service.Gas.AccelerateTransaction(ctx, req.ChainId, req.TxHash)
	if err != nil {
		return nil, err
	}

	return &v1.AccelerateTransactionRes{
		NewTxHash: newTxHash,
	}, nil
}

// CancelTransaction 取消交易
func (c *GasController) CancelTransaction(ctx context.Context, req *v1.CancelTransactionReq) (res *v1.CancelTransactionRes, err error) {
	newTxHash, err := service.Gas.CancelTransaction(ctx, req.ChainId, req.TxHash)
	if err != nil {
		return nil, err
	}

	return &v1.CancelTransactionRes{
		NewTxHash: newTxHash,
	}, nil
}

// GetGasInfo 获取Gas信息
func (c *GasController) GetGasInfo(ctx context.Context, req *v1.GetGasInfoReq) (res *v1.GetGasInfoRes, err error) {
	info, err := service.Gas.GetGasInfo(ctx, req.ChainId)
	if err != nil {
		return nil, err
	}

	return &v1.GetGasInfoRes{
		Info: info,
	}, nil
}
