package service

import (
	"context"
	"go-wallet-defi/internal/logic"
	"go-wallet-defi/internal/model"
)

type IBridge interface {
	// CrossTransfer 跨链转账
	CrossTransfer(ctx context.Context, fromChainId, toChainId uint64, fromAddress, toAddress, tokenAddress, amount string) (hash string, nonce uint64, err error)

	// GetCrossTransfers 获取跨链交易列表
	GetCrossTransfers(ctx context.Context, fromChainId, toChainId uint64, address string, status, page, pageSize int) ([]*model.CrossTransfer, int, error)

	// ProcessLockEvent 处理锁定事件
	ProcessLockEvent(ctx context.Context, chainId uint64, token, from string, amount string, toChainId uint64, toAddress string, nonce uint64, hash string) error

	// ProcessUnlockEvent 处理解锁事件
	ProcessUnlockEvent(ctx context.Context, chainId uint64, token, to string, amount string, fromChainId uint64, nonce uint64, hash string) error
}

// Bridge 获取跨链桥服务
func Bridge() IBridge {
	if localBridge == nil {
		localBridge = &logic.BridgeLogic{}
	}
	return localBridge
}

var localBridge IBridge
