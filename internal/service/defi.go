package service

import (
	"context"
	"go-wallet-defi/internal/logic"
)

type IDefi interface {
	// Swap 代币兑换
	Swap(ctx context.Context, chainId uint64, fromToken, toToken string, amount string, fromAddress string, slippageBps int, swapType string) (hash string, amountOut string, err error)

	// AddLiquidity 添加流动性
	AddLiquidity(ctx context.Context, chainId uint64, tokenA, tokenB string, amountA, amountB string, fromAddress string, slippageBps int) (hash string, liquidity string, err error)

	// RemoveLiquidity 移除流动性
	RemoveLiquidity(ctx context.Context, chainId uint64, pair string, liquidity string, fromAddress string, slippageBps int) (hash string, amount0, amount1 string, err error)

	// Supply 存款
	Supply(ctx context.Context, chainId uint64, pool, token string, amount string, fromAddress string) (hash string, err error)

	// Withdraw 提款
	Withdraw(ctx context.Context, chainId uint64, pool, token string, amount string, fromAddress string) (hash string, err error)

	// Borrow 借款
	Borrow(ctx context.Context, chainId uint64, pool, token string, amount string, rateMode int, fromAddress string) (hash string, err error)

	// Repay 还款
	Repay(ctx context.Context, chainId uint64, pool, token string, amount string, rateMode int, fromAddress string) (hash string, err error)

	// Stake 质押
	Stake(ctx context.Context, chainId uint64, pool string, amount string, fromAddress string) (hash string, err error)

	// Unstake 解质押
	Unstake(ctx context.Context, chainId uint64, pool string, amount string, fromAddress string) (hash string, err error)

	// ClaimReward 领取奖励
	ClaimReward(ctx context.Context, chainId uint64, pool string, fromAddress string) (hash string, reward string, err error)

	// DepositVault 存入机枪池
	DepositVault(ctx context.Context, chainId uint64, vault string, amount string, fromAddress string) (hash string, shares string, err error)

	// WithdrawVault 提取机枪池
	WithdrawVault(ctx context.Context, chainId uint64, vault string, shares string, fromAddress string) (hash string, amount string, err error)
}

// Defi 获取DeFi服务
func Defi() IDefi {
	if localDefi == nil {
		localDefi = &logic.DefiLogic{}
	}
	return localDefi
}

var localDefi IDefi
