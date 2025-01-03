package controller

import (
	"context"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/service"
)

type TransactionController struct{}

// TransferEth ETH转账
func (c *TransactionController) TransferEth(ctx context.Context, req *v1.TransferEthReq) (res *v1.TransferEthRes, err error) {
	hash, err := service.Transaction().TransferEth(ctx, req.From, req.To, req.Amount, req.GasPrice, req.GasLimit)
	if err != nil {
		return nil, err
	}

	return &v1.TransferEthRes{
		Hash: hash,
	}, nil
}

// TransferToken 代币转账
func (c *TransactionController) TransferToken(ctx context.Context, req *v1.TransferTokenReq) (res *v1.TransferTokenRes, err error) {
	hash, err := service.Transaction().TransferToken(ctx, req.From, req.To, req.Amount, req.Token, req.GasPrice, req.GasLimit)
	if err != nil {
		return nil, err
	}

	return &v1.TransferTokenRes{
		Hash: hash,
	}, nil
}

// GetTransactions 获取交易记录
func (c *TransactionController) GetTransactions(ctx context.Context, req *v1.GetTransactionsReq) (res *v1.GetTransactionsRes, err error) {
	transactions, total, err := service.Transaction().GetTransactions(ctx, req.Address, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]v1.TransactionInfo, 0, len(transactions))
	for _, tx := range transactions {
		list = append(list, v1.TransactionInfo{
			Hash:        tx.Hash,
			From:        tx.FromAddress,
			To:          tx.ToAddress,
			Amount:      tx.Amount,
			Token:       tx.TokenAddress,
			Status:      tx.Status,
			BlockNumber: tx.BlockNumber,
			BlockTime:   tx.BlockTime,
			CreatedAt:   tx.CreatedAt,
		})
	}

	return &v1.GetTransactionsRes{
		List:  list,
		Total: total,
	}, nil
}
