package defi

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

// UniswapV2Router UniswapV2路由合约
type UniswapV2Router struct {
	address common.Address
	abi     abi.ABI
	client  *ethclient.Client
}

// NewUniswapV2Router 创建UniswapV2路由实例
func NewUniswapV2Router(address common.Address, client *ethclient.Client) (*UniswapV2Router, error) {
	parsed, err := abi.JSON(strings.NewReader(UniswapV2RouterABI))
	if err != nil {
		return nil, err
	}

	return &UniswapV2Router{
		address: address,
		abi:     parsed,
		client:  client,
	}, nil
}

// SwapExactTokensForTokens 按精确输入兑换代币
func (r *UniswapV2Router) PackSwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) ([]byte, error) {
	return r.abi.Pack("swapExactTokensForTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapTokensForExactTokens 按精确输出兑换代币
func (r *UniswapV2Router) PackSwapTokensForExactTokens(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) ([]byte, error) {
	return r.abi.Pack("swapTokensForExactTokens", amountOut, amountInMax, path, to, deadline)
}

// AddLiquidity 添加流动性
func (r *UniswapV2Router) PackAddLiquidity(tokenA, tokenB common.Address, amountA, amountB, amountAMin, amountBMin *big.Int, to common.Address, deadline *big.Int) ([]byte, error) {
	return r.abi.Pack("addLiquidity", tokenA, tokenB, amountA, amountB, amountAMin, amountBMin, to, deadline)
}

// RemoveLiquidity 移除流动性
func (r *UniswapV2Router) PackRemoveLiquidity(tokenA, tokenB common.Address, liquidity *big.Int, amountAMin, amountBMin *big.Int, to common.Address, deadline *big.Int) ([]byte, error) {
	return r.abi.Pack("removeLiquidity", tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// AavePool Aave借贷池合约
type AavePool struct {
	address common.Address
	abi     abi.ABI
	client  *ethclient.Client
}

// NewAavePool 创建Aave借贷池实例
func NewAavePool(address common.Address, client *ethclient.Client) (*AavePool, error) {
	parsed, err := abi.JSON(strings.NewReader(AavePoolABI))
	if err != nil {
		return nil, err
	}

	return &AavePool{
		address: address,
		abi:     parsed,
		client:  client,
	}, nil
}

// Supply 存款
func (p *AavePool) PackSupply(asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) ([]byte, error) {
	return p.abi.Pack("supply", asset, amount, onBehalfOf, referralCode)
}

// Withdraw 提款
func (p *AavePool) PackWithdraw(asset common.Address, amount *big.Int, to common.Address) ([]byte, error) {
	return p.abi.Pack("withdraw", asset, amount, to)
}

// Borrow 借款
func (p *AavePool) PackBorrow(asset common.Address, amount *big.Int, interestRateMode uint8, referralCode uint16, onBehalfOf common.Address) ([]byte, error) {
	return p.abi.Pack("borrow", asset, amount, interestRateMode, referralCode, onBehalfOf)
}

// Repay 还款
func (p *AavePool) PackRepay(asset common.Address, amount *big.Int, interestRateMode uint8, onBehalfOf common.Address) ([]byte, error) {
	return p.abi.Pack("repay", asset, amount, interestRateMode, onBehalfOf)
}

// YearnVault Yearn机枪池合约
type YearnVault struct {
	address common.Address
	abi     abi.ABI
	client  *ethclient.Client
}

// NewYearnVault 创建Yearn机枪池实例
func NewYearnVault(address common.Address, client *ethclient.Client) (*YearnVault, error) {
	parsed, err := abi.JSON(strings.NewReader(YearnVaultABI))
	if err != nil {
		return nil, err
	}

	return &YearnVault{
		address: address,
		abi:     parsed,
		client:  client,
	}, nil
}

// Deposit 存入资金
func (v *YearnVault) PackDeposit(amount *big.Int) ([]byte, error) {
	return v.abi.Pack("deposit", amount)
}

// Withdraw 提取资金
func (v *YearnVault) PackWithdraw(shares *big.Int, recipient common.Address) ([]byte, error) {
	return v.abi.Pack("withdraw", shares, recipient)
}

// PricePerShare 获取每份额价格
func (v *YearnVault) PricePerShare(ctx context.Context) (*big.Int, error) {
	data, err := v.abi.Pack("pricePerShare")
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &v.address,
		Data: data,
	}

	output, err := v.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	var result *big.Int
	err = v.abi.UnpackIntoInterface(&result, "pricePerShare", output)
	return result, err
}

// Farm 收益农场合约
type Farm struct {
	address common.Address
	abi     abi.ABI
	client  *ethclient.Client
}

// NewFarm 创建收益农场实例
func NewFarm(address common.Address, client *ethclient.Client) (*Farm, error) {
	parsed, err := abi.JSON(strings.NewReader(FarmABI))
	if err != nil {
		return nil, err
	}

	return &Farm{
		address: address,
		abi:     parsed,
		client:  client,
	}, nil
}

// Stake 质押代币
func (f *Farm) PackStake(amount *big.Int) ([]byte, error) {
	return f.abi.Pack("stake", amount)
}

// Withdraw 提取代币
func (f *Farm) PackWithdraw(amount *big.Int) ([]byte, error) {
	return f.abi.Pack("withdraw", amount)
}

// GetReward 领取奖励
func (f *Farm) PackGetReward() ([]byte, error) {
	return f.abi.Pack("getReward")
}

// Earned 获取待领取奖励
func (f *Farm) Earned(ctx context.Context, account common.Address) (*big.Int, error) {
	data, err := f.abi.Pack("earned", account)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &f.address,
		Data: data,
	}

	output, err := f.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	var result *big.Int
	err = f.abi.UnpackIntoInterface(&result, "earned", output)
	return result, err
}
