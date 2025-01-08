package controller

import (
	"context"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/service"
)

type DefiController struct{}

// SwapToken 代币兑换
func (c *DefiController) SwapToken(ctx context.Context, req *v1.SwapTokenReq) (res *v1.SwapTokenRes, err error) {
	hash, amountOut, err := service.Defi().Swap(ctx,
		req.ChainId,
		req.FromToken,
		req.ToToken,
		req.Amount,
		req.FromAddress,
		req.SlippageBps,
		req.Type,
	)
	if err != nil {
		return nil, err
	}

	return &v1.SwapTokenRes{
		Hash:      hash,
		AmountOut: amountOut,
	}, nil
}

// AddLiquidity 添加流动性
func (c *DefiController) AddLiquidity(ctx context.Context, req *v1.AddLiquidityReq) (res *v1.AddLiquidityRes, err error) {
	hash, liquidity, err := service.Defi().AddLiquidity(ctx,
		req.ChainId,
		req.TokenA,
		req.TokenB,
		req.AmountA,
		req.AmountB,
		req.FromAddress,
		req.SlippageBps,
	)
	if err != nil {
		return nil, err
	}

	return &v1.AddLiquidityRes{
		Hash:      hash,
		Liquidity: liquidity,
	}, nil
}

// RemoveLiquidity 移除流动性
func (c *DefiController) RemoveLiquidity(ctx context.Context, req *v1.RemoveLiquidityReq) (res *v1.RemoveLiquidityRes, err error) {
	hash, amount0, amount1, err := service.Defi().RemoveLiquidity(ctx,
		req.ChainId,
		req.Pair,
		req.Liquidity,
		req.FromAddress,
		req.SlippageBps,
	)
	if err != nil {
		return nil, err
	}

	return &v1.RemoveLiquidityRes{
		Hash:    hash,
		Amount0: amount0,
		Amount1: amount1,
	}, nil
}

// Supply 存款
func (c *DefiController) Supply(ctx context.Context, req *v1.SupplyReq) (res *v1.SupplyRes, err error) {
	hash, err := service.Defi().Supply(ctx,
		req.ChainId,
		req.Pool,
		req.Token,
		req.Amount,
		req.FromAddress,
	)
	if err != nil {
		return nil, err
	}

	return &v1.SupplyRes{Hash: hash}, nil
}

// Borrow 借款
func (c *DefiController) Borrow(ctx context.Context, req *v1.BorrowReq) (res *v1.BorrowRes, err error) {
	hash, err := service.Defi().Borrow(ctx,
		req.ChainId,
		req.Pool,
		req.Token,
		req.Amount,
		req.RateMode,
		req.FromAddress,
	)
	if err != nil {
		return nil, err
	}

	return &v1.BorrowRes{Hash: hash}, nil
}

// Repay 还款
func (c *DefiController) Repay(ctx context.Context, req *v1.RepayReq) (res *v1.RepayRes, err error) {
	hash, err := service.Defi().Repay(ctx,
		req.ChainId,
		req.Pool,
		req.Token,
		req.Amount,
		req.RateMode,
		req.FromAddress,
	)
	if err != nil {
		return nil, err
	}

	return &v1.RepayRes{Hash: hash}, nil
}

// Stake 质押
func (c *DefiController) Stake(ctx context.Context, req *v1.StakeReq) (res *v1.StakeRes, err error) {
	hash, err := service.Defi().Stake(ctx,
		req.ChainId,
		req.Pool,
		req.Amount,
		req.FromAddress,
	)
	if err != nil {
		return nil, err
	}

	return &v1.StakeRes{Hash: hash}, nil
}

// Unstake 解质押
func (c *DefiController) Unstake(ctx context.Context, req *v1.UnstakeReq) (res *v1.UnstakeRes, err error) {
	hash, err := service.Defi().Unstake(ctx,
		req.ChainId,
		req.Pool,
		req.Amount,
		req.FromAddress,
	)
	if err != nil {
		return nil, err
	}

	return &v1.UnstakeRes{Hash: hash}, nil
}

// ClaimReward 领取奖励
func (c *DefiController) ClaimReward(ctx context.Context, req *v1.ClaimRewardReq) (res *v1.ClaimRewardRes, err error) {
	hash, reward, err := service.Defi().ClaimReward(ctx,
		req.ChainId,
		req.Pool,
		req.FromAddress,
	)
	if err != nil {
		return nil, err
	}

	return &v1.ClaimRewardRes{
		Hash:   hash,
		Reward: reward,
	}, nil
}

// DepositVault 存入机枪池
func (c *DefiController) DepositVault(ctx context.Context, req *v1.DepositVaultReq) (res *v1.DepositVaultRes, err error) {
	hash, shares, err := service.Defi().DepositVault(ctx,
		req.ChainId,
		req.Vault,
		req.Amount,
		req.FromAddress,
	)
	if err != nil {
		return nil, err
	}

	return &v1.DepositVaultRes{
		Hash:   hash,
		Shares: shares,
	}, nil
}

// WithdrawVault 提取机枪池
func (c *DefiController) WithdrawVault(ctx context.Context, req *v1.WithdrawVaultReq) (res *v1.WithdrawVaultRes, err error) {
	hash, amount, err := service.Defi().WithdrawVault(ctx,
		req.ChainId,
		req.Vault,
		req.Shares,
		req.FromAddress,
	)
	if err != nil {
		return nil, err
	}

	return &v1.WithdrawVaultRes{
		Hash:   hash,
		Amount: amount,
	}, nil
}
