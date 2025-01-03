package service

import (
	"context"
	"go-wallet-defi/internal/model"
)

type ITransaction interface {
	// TransferEth ETH转账
	TransferEth(ctx context.Context, from, to, amount string, gasPrice string, gasLimit uint64) (hash string, err error)

	// TransferToken 代币转账
	TransferToken(ctx context.Context, from, to, amount, token string, gasPrice string, gasLimit uint64) (hash string, err error)

	// GetTransactions 获取交易记录
	GetTransactions(ctx context.Context, address string, page, pageSize int) ([]*model.Transaction, int, error)

	// UpdateTransactionStatus 更新交易状态
	UpdateTransactionStatus(ctx context.Context, hash string) error
}
