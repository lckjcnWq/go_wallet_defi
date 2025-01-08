package logic

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/contracts/defi"
	"go-wallet-defi/internal/pkg/contracts/token"
	"go-wallet-defi/internal/pkg/ethclientx"
	"math/big"
	"time"
)

type DefiLogic struct{}

// Swap 代币兑换
func (s *DefiLogic) Swap(ctx context.Context, chainId uint64, fromToken, toToken string, amount string, fromAddress string, slippageBps int, swapType string) (hash string, amountOut string, err error) {
	//1.获取客户端
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", "", err
	}
	//2.创建路由合约实例
	router, err := defi.NewUniswapV2Router(common.HexToAddress(""), client)
	if err != nil {
		return "", "", err
	}
	//3.计算滑点
	amountBig, _ := new(big.Int).SetString(amount, 10)
	slippage := new(big.Int).Mul(amountBig, big.NewInt(int64(slippageBps)))
	slippage = new(big.Int).Div(slippage, big.NewInt(10000))
	//4.设置交易路径
	path := []common.Address{
		common.HexToAddress(fromToken),
		common.HexToAddress(toToken),
	}
	//5.获取交易deadline
	deadline := big.NewInt(time.Now().Unix() + 1200) //20分钟超时
	var data []byte
	if swapType == "EXACT_INPUT" {
		// 精确输入兑换
		data, err = router.PackSwapExactTokensForTokens(
			amountBig,
			new(big.Int).Sub(amountBig, slippage),
			path,
			common.HexToAddress(fromAddress),
			deadline,
		)
	} else {
		// 精确输出兑换
		data, err = router.PackSwapTokensForExactTokens(
			amountBig,
			new(big.Int).Add(amountBig, slippage),
			path,
			common.HexToAddress(fromAddress),
			deadline,
		)
	}
	if err != nil {
		return "", "", err
	}
	// 如果fromToken不是ETH,需要先approve
	if fromToken != "0x0000000000000000000000000000000000000000" {
		erc20, err := token.NewERC20(common.HexToAddress(fromToken), client)
		if err != nil {
			return "", "", err
		}

		approveData, err := erc20.PackApprove(common.HexToAddress("router address"), amountBig)
		if err != nil {
			return "", "", err
		}

		approveHash, err := s.sendTransaction(ctx, client, fromAddress, fromToken, big.NewInt(0), approveData)
		if err != nil {
			return "", "", err
		}

		// 等待approve交易确认
		_, err = s.waitTransaction(ctx, client, approveHash)
		if err != nil {
			return "", "", err
		}
	}

	// 发送交易
	hash, err = s.sendTransaction(ctx, client, fromAddress, "router address", big.NewInt(0), data)
	if err != nil {
		return "", "", err
	}

	// 保存交易记录
	trade := &model.DexTrade{
		ChainId:    chainId,
		FromToken:  fromToken,
		ToToken:    toToken,
		FromAmount: amount,
		ToAmount:   "0", // 等待交易完成后更新
		User:       fromAddress,
		Router:     "router address",
		Path:       "[\"" + fromToken + "\",\"" + toToken + "\"]",
		Type:       swapType,
		Hash:       hash,
		Status:     0,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	err = dao.Defi.InsertDexTrade(ctx, trade)
	if err != nil {
		return "", "", err
	}

	return hash, "0", nil
}

// AddLiquidity 添加流动性
func (s *DefiLogic) AddLiquidity(ctx context.Context, chainId uint64, tokenA, tokenB string, amountA, amountB string, fromAddress string, slippageBps int) (hash string, liquidity string, err error) {
	// 获取客户端
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", "", err
	}

	// 创建路由合约实例
	router, err := defi.NewUniswapV2Router(common.HexToAddress("router address"), client)
	if err != nil {
		return "", "", err
	}

	// 计算滑点
	amountABig, _ := new(big.Int).SetString(amountA, 10)
	amountBBig, _ := new(big.Int).SetString(amountB, 10)

	slippageA := new(big.Int).Mul(amountABig, big.NewInt(int64(slippageBps)))
	slippageA = new(big.Int).Div(slippageA, big.NewInt(10000))

	slippageB := new(big.Int).Mul(amountBBig, big.NewInt(int64(slippageBps)))
	slippageB = new(big.Int).Div(slippageB, big.NewInt(10000))

	// 获取交易deadline
	deadline := big.NewInt(time.Now().Unix() + 1200)

	// 构造交易数据
	data, err := router.PackAddLiquidity(
		common.HexToAddress(tokenA),
		common.HexToAddress(tokenB),
		amountABig,
		amountBBig,
		new(big.Int).Sub(amountABig, slippageA),
		new(big.Int).Sub(amountBBig, slippageB),
		common.HexToAddress(fromAddress),
		deadline,
	)
	if err != nil {
		return "", "", err
	}

	// approve tokenA
	if tokenA != "0x0000000000000000000000000000000000000000" {
		erc20A, err := token.NewERC20(common.HexToAddress(tokenA), client)
		if err != nil {
			return "", "", err
		}

		approveData, err := erc20A.PackApprove(common.HexToAddress("router address"), amountABig)
		if err != nil {
			return "", "", err
		}

		approveHash, err := s.sendTransaction(ctx, client, fromAddress, tokenA, big.NewInt(0), approveData)
		if err != nil {
			return "", "", err
		}

		_, err = s.waitTransaction(ctx, client, approveHash)
		if err != nil {
			return "", "", err
		}
	}

	// approve tokenB
	if tokenB != "0x0000000000000000000000000000000000000000" {
		erc20B, err := token.NewERC20(common.HexToAddress(tokenB), client)
		if err != nil {
			return "", "", err
		}

		approveData, err := erc20B.PackApprove(common.HexToAddress("router address"), amountBBig)
		if err != nil {
			return "", "", err
		}

		approveHash, err := s.sendTransaction(ctx, client, fromAddress, tokenB, big.NewInt(0), approveData)
		if err != nil {
			return "", "", err
		}

		_, err = s.waitTransaction(ctx, client, approveHash)
		if err != nil {
			return "", "", err
		}
	}

	// 发送交易
	hash, err = s.sendTransaction(ctx, client, fromAddress, "router address", big.NewInt(0), data)
	if err != nil {
		return "", "", err
	}

	// 保存记录
	liquidity := &model.Liquidity{
		ChainId:   chainId,
		Token0:    tokenA,
		Token1:    tokenB,
		Amount0:   amountA,
		Amount1:   amountB,
		Liquidity: "0", // 等待交易完成后更新
		User:      fromAddress,
		Type:      "ADD",
		Hash:      hash,
		Status:    0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = dao.Defi.InsertLiquidity(ctx, liquidity)
	if err != nil {
		return "", "", err
	}

	return hash, "0", nil
}

// RemoveLiquidity 移除流动性
func (s *DefiLogic) RemoveLiquidity(ctx context.Context, chainId uint64, pair string, liquidity string, fromAddress string, slippageBps int) (hash string, amount0, amount1 string, err error) {
	// 获取客户端
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", "", "", err
	}

	// 获取交易对中的代币
	pairContract, err := defi.NewUniswapV2Pair(common.HexToAddress(pair), client)
	if err != nil {
		return "", "", "", err
	}

	token0, err := pairContract.Token0()
	if err != nil {
		return "", "", "", err
	}

	token1, err := pairContract.Token1()
	if err != nil {
		return "", "", "", err
	}

	// 创建路由合约实例
	router, err := defi.NewUniswapV2Router(common.HexToAddress("router address"), client)
	if err != nil {
		return "", "", "", err
	}

	liquidityBig, _ := new(big.Int).SetString(liquidity, 10)

	// 获取当前储备量
	reserve0, reserve1, _, err := pairContract.GetReserves()
	if err != nil {
		return "", "", "", err
	}

	// 计算预期获得的代币数量
	amount0Big := new(big.Int).Mul(liquidityBig, reserve0)
	amount0Big = new(big.Int).Div(amount0Big, new(big.Int).SetUint64(1e18))

	amount1Big := new(big.Int).Mul(liquidityBig, reserve1)
	amount1Big = new(big.Int).Div(amount1Big, new(big.Int).SetUint64(1e18))

	// 计算滑点
	slippage0 := new(big.Int).Mul(amount0Big, big.NewInt(int64(slippageBps)))
	slippage0 = new(big.Int).Div(slippage0, big.NewInt(10000))

	slippage1 := new(big.Int).Mul(amount1Big, big.NewInt(int64(slippageBps)))
	slippage1 = new(big.Int).Div(slippage1, big.NewInt(10000))

	// 获取交易deadline
	deadline := big.NewInt(time.Now().Unix() + 1200)

	// 构造交易数据
	data, err := router.PackRemoveLiquidity(
		token0,
		token1,
		liquidityBig,
		new(big.Int).Sub(amount0Big, slippage0),
		new(big.Int).Sub(amount1Big, slippage1),
		common.HexToAddress(fromAddress),
		deadline,
	)
	if err != nil {
		return "", "", "", err
	}

	// approve LP token
	erc20, err := token.NewERC20(common.HexToAddress(pair), client)
	if err != nil {
		return "", "", "", err
	}

	approveData, err := erc20.PackApprove(common.HexToAddress("router address"), liquidityBig)
	if err != nil {
		return "", "", "", err
	}

	approveHash, err := s.sendTransaction(ctx, client, fromAddress, pair, big.NewInt(0), approveData)
	if err != nil {
		return "", "", "", err
	}

	_, err = s.waitTransaction(ctx, client, approveHash)
	if err != nil {
		return "", "", "", err
	}

	// 发送交易
	hash, err = s.sendTransaction(ctx, client, fromAddress, "router address", big.NewInt(0), data)
	if err != nil {
		return "", "", "", err
	}

	// 保存记录
	liquidityRecord := &model.Liquidity{
		ChainId:   chainId,
		Pair:      pair,
		Token0:    token0.Hex(),
		Token1:    token1.Hex(),
		Amount0:   amount0Big.String(),
		Amount1:   amount1Big.String(),
		Liquidity: liquidity,
		User:      fromAddress,
		Type:      "REMOVE",
		Hash:      hash,
		Status:    0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = dao.Defi.InsertLiquidity(ctx, liquidityRecord)
	if err != nil {
		return "", "", "", err
	}

	return hash, amount0Big.String(), amount1Big.String(), nil
}

// Supply 存款到Aave
func (s *DefiLogic) Supply(ctx context.Context, chainId uint64, pool, token string, amount string, fromAddress string) (hash string, err error) {
	// 获取客户端
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", err
	}

	// 创建借贷池合约实例
	aavePool, err := defi.NewAavePool(common.HexToAddress(pool), client)
	if err != nil {
		return "", err
	}

	amountBig, _ := new(big.Int).SetString(amount, 10)

	// 构造交易数据
	data, err := aavePool.PackSupply(
		common.HexToAddress(token),
		amountBig,
		common.HexToAddress(fromAddress),
		0,
	)
	if err != nil {
		return "", err
	}

	// approve
	if token != "0x0000000000000000000000000000000000000000" {
		erc20, err := token.NewERC20(common.HexToAddress(token), client)
		if err != nil {
			return "", err
		}

		approveData, err := erc20.PackApprove(common.HexToAddress(pool), amountBig)
		if err != nil {
			return "", err
		}

		approveHash, err := s.sendTransaction(ctx, client, fromAddress, token, big.NewInt(0), approveData)
		if err != nil {
			return "", err
		}

		_, err = s.waitTransaction(ctx, client, approveHash)
		if err != nil {
			return "", err
		}
	}

	// 发送交易
	hash, err = s.sendTransaction(ctx, client, fromAddress, pool, big.NewInt(0), data)
	if err != nil {
		return "", err
	}
	// 接上一段Supply方法的代码...

	// 保存交易记录
	lending := &model.Lending{
		ChainId:   chainId,
		Pool:      pool,
		Token:     token,
		Amount:    amount,
		User:      fromAddress,
		Type:      "SUPPLY",
		Hash:      hash,
		Status:    0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = dao.Defi.InsertLending(ctx, lending)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// Withdraw 从Aave提取
func (s *DefiLogic) Withdraw(ctx context.Context, chainId uint64, pool, token string, amount string, fromAddress string) (hash string, err error) {
	// 获取客户端
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", err
	}

	// 创建借贷池合约实例
	aavePool, err := defi.NewAavePool(common.HexToAddress(pool), client)
	if err != nil {
		return "", err
	}

	amountBig, _ := new(big.Int).SetString(amount, 10)

	// 构造交易数据
	data, err := aavePool.PackWithdraw(
		common.HexToAddress(token),
		amountBig,
		common.HexToAddress(fromAddress),
	)
	if err != nil {
		return "", err
	}

	// 发送交易
	hash, err = s.sendTransaction(ctx, client, fromAddress, pool, big.NewInt(0), data)
	if err != nil {
		return "", err
	}

	// 保存交易记录
	lending := &model.Lending{
		ChainId:   chainId,
		Pool:      pool,
		Token:     token,
		Amount:    amount,
		User:      fromAddress,
		Type:      "WITHDRAW",
		Hash:      hash,
		Status:    0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = dao.Defi.InsertLending(ctx, lending)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// Borrow 从Aave借款
func (s *DefiLogic) Borrow(ctx context.Context, chainId uint64, pool, token string, amount string, rateMode int, fromAddress string) (hash string, err error) {
	client, err := ethclient.GetClient(chainId)
	if err != nil {
		return "", err
	}

	aavePool, err := defi.NewAavePool(common.HexToAddress(pool), client)
	if err != nil {
		return "", err
	}

	amountBig, _ := new(big.Int).SetString(amount, 10)

	data, err := aavePool.PackBorrow(
		common.HexToAddress(token),
		amountBig,
		uint8(rateMode),
		0,
		common.HexToAddress(fromAddress),
	)
	if err != nil {
		return "", err
	}

	hash, err = s.sendTransaction(ctx, client, fromAddress, pool, big.NewInt(0), data)
	if err != nil {
		return "", err
	}

	lending := &model.Lending{
		ChainId:   chainId,
		Pool:      pool,
		Token:     token,
		Amount:    amount,
		User:      fromAddress,
		Type:      "BORROW",
		RateMode:  rateMode,
		Hash:      hash,
		Status:    0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = dao.Defi.InsertLending(ctx, lending)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// Repay 向Aave还款
func (s *DefiLogic) Repay(ctx context.Context, chainId uint64, pool, token string, amount string, rateMode int, fromAddress string) (hash string, err error) {
	client, err := ethclient.GetClient(chainId)
	if err != nil {
		return "", err
	}

	aavePool, err := defi.NewAavePool(common.HexToAddress(pool), client)
	if err != nil {
		return "", err
	}

	amountBig, _ := new(big.Int).SetString(amount, 10)

	data, err := aavePool.PackRepay(
		common.HexToAddress(token),
		amountBig,
		uint8(rateMode),
		common.HexToAddress(fromAddress),
	)
	if err != nil {
		return "", err
	}

	// approve
	if token != "0x0000000000000000000000000000000000000000" {
		erc20, err := token.NewERC20(common.HexToAddress(token), client)
		if err != nil {
			return "", err
		}

		approveData, err := erc20.PackApprove(common.HexToAddress(pool), amountBig)
		if err != nil {
			return "", err
		}

		approveHash, err := s.sendTransaction(ctx, client, fromAddress, token, big.NewInt(0), approveData)
		if err != nil {
			return "", err
		}

		_, err = s.waitTransaction(ctx, client, approveHash)
		if err != nil {
			return "", err
		}
	}

	hash, err = s.sendTransaction(ctx, client, fromAddress, pool, big.NewInt(0), data)
	if err != nil {
		return "", err
	}

	lending := &model.Lending{
		ChainId:   chainId,
		Pool:      pool,
		Token:     token,
		Amount:    amount,
		User:      fromAddress,
		Type:      "REPAY",
		RateMode:  rateMode,
		Hash:      hash,
		Status:    0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = dao.Defi.InsertLending(ctx, lending)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// Stake 质押到收益农场
func (s *DefiLogic) Stake(ctx context.Context, chainId uint64, pool string, amount string, fromAddress string) (hash string, err error) {
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", err
	}

	farm, err := defi.NewFarm(common.HexToAddress(pool), client)
	if err != nil {
		return "", err
	}

	amountBig, _ := new(big.Int).SetString(amount, 10)

	data, err := farm.PackStake(amountBig)
	if err != nil {
		return "", err
	}

	// approve stake token
	stakeToken := "" // 从合约获取质押代币地址
	if stakeToken != "0x0000000000000000000000000000000000000000" {
		erc20, err := token.NewERC20(common.HexToAddress(stakeToken), client)
		if err != nil {
			return "", err
		}

		approveData, err := erc20.PackApprove(common.HexToAddress(pool), amountBig)
		if err != nil {
			return "", err
		}

		approveHash, err := s.sendTransaction(ctx, client, fromAddress, stakeToken, big.NewInt(0), approveData)
		if err != nil {
			return "", err
		}

		_, err = s.waitTransaction(ctx, client, approveHash)
		if err != nil {
			return "", err
		}
	}

	hash, err = s.sendTransaction(ctx, client, fromAddress, pool, big.NewInt(0), data)
	if err != nil {
		return "", err
	}

	farm := &model.YieldFarm{
		ChainId:     chainId,
		Pool:        pool,
		StakeToken:  stakeToken,
		RewardToken: "", // 从合约获取奖励代币地址
		Amount:      amount,
		User:        fromAddress,
		Type:        "STAKE",
		Hash:        hash,
		Status:      0,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = dao.Defi.InsertYieldFarm(ctx, farm)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// Unstake 从收益农场解除质押
func (s *DefiLogic) Unstake(ctx context.Context, chainId uint64, pool string, amount string, fromAddress string) (hash string, err error) {
	client, err := ethclient.GetClient(chainId)
	if err != nil {
		return "", err
	}

	farm, err := defi.NewFarm(common.HexToAddress(pool), client)
	if err != nil {
		return "", err
	}

	amountBig, _ := new(big.Int).SetString(amount, 10)

	data, err := farm.PackWithdraw(amountBig)
	if err != nil {
		return "", err
	}

	hash, err = s.sendTransaction(ctx, client, fromAddress, pool, big.NewInt(0), data)
	if err != nil {
		return "", err
	}

	farmRecord := &model.YieldFarm{
		ChainId:     chainId,
		Pool:        pool,
		StakeToken:  "", // 从合约获取
		RewardToken: "", // 从合约获取
		Amount:      amount,
		User:        fromAddress,
		Type:        "UNSTAKE",
		Hash:        hash,
		Status:      0,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = dao.Defi.InsertYieldFarm(ctx, farmRecord)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// ClaimReward 领取农场奖励
func (s *DefiLogic) ClaimReward(ctx context.Context, chainId uint64, pool string, fromAddress string) (hash string, reward string, err error) {
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", "", err
	}

	farm, err := defi.NewFarm(common.HexToAddress(pool), client)
	if err != nil {
		return "", "", err
	}

	// 获取待领取奖励
	rewardBig, err := farm.Earned(ctx, common.HexToAddress(fromAddress))
	if err != nil {
		return "", "", err
	}

	data, err := farm.PackGetReward()
	if err != nil {
		return "", "", err
	}

	hash, err = s.sendTransaction(ctx, client, fromAddress, pool, big.NewInt(0), data)
	if err != nil {
		return "", "", err
	}

	farmRecord := &model.YieldFarm{
		ChainId:     chainId,
		Pool:        pool,
		StakeToken:  "", // 从合约获取
		RewardToken: "", // 从合约获取
		Amount:      "0",
		Reward:      rewardBig.String(),
		User:        fromAddress,
		Type:        "CLAIM",
		Hash:        hash,
		Status:      0,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = dao.Defi.InsertYieldFarm(ctx, farmRecord)
	if err != nil {
		return "", "", err
	}

	return hash, rewardBig.String(), nil
}

// DepositVault 存入机枪池
func (s *DefiLogic) DepositVault(ctx context.Context, chainId uint64, vault string, amount string, fromAddress string) (hash string, shares string, err error) {
	client, err := ethclientx.GetClientByChainId(ctx, chainId)
	if err != nil {
		return "", "", err
	}

	yearnVault, err := defi.NewYearnVault(common.HexToAddress(vault), client)
	if err != nil {
		return "", "", err
	}

	amountBig, _ := new(big.Int).SetString(amount, 10)

	// 获取当前价格比率计算份额
	pricePerShare, err := yearnVault.PricePerShare(ctx)
	if err != nil {
		return "", "", err
	}

	sharesBig := new(big.Int).Mul(amountBig, big.NewInt(1e18))
	sharesBig = sharesBig.Div(sharesBig, pricePerShare)

	data, err := yearnVault.PackDeposit(amountBig)
	if err != nil {
		return "", "", err
	}

	// approve token
	token := "" // 从合约获取存款代币地址
	if token != "0x0000000000000000000000000000000000000000" {
		erc20, err := token.NewERC20(common.HexToAddress(token), client)
		if err != nil {
			return "", "", err
		}

		approveData, err := erc20.PackApprove(common.HexToAddress(vault), amountBig)
		if err != nil {
			return "", "", err
		}

		approveHash, err := s.sendTransaction(ctx, client, fromAddress, token, big.NewInt(0), approveData)
		if err != nil {
			return "", "", err
		}

		_, err = s.waitTransaction(ctx, client, approveHash)
		if err != nil {
			return "", "", err
		}
	}

	hash, err = s.sendTransaction(ctx, client, fromAddress, vault, big.NewInt(0), data)
	if err != nil {
		return "", "", err
	}

	vaultRecord := &model.Vault{
		ChainId:   chainId,
		Vault:     vault,
		Token:     token,
		Amount:    amount,
		Shares:    sharesBig.String(),
		User:      fromAddress,
		Type:      "DEPOSIT",
		Hash:      hash,
		Status:    0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = dao.Defi.InsertVault(ctx, vaultRecord)
	if err != nil {
		return "", "", err
	}

	return hash, sharesBig.String(), nil
}
