package v1

import "github.com/gogf/gf/v2/frame/g"

// SwapTokenReq 代币兑换请求
type SwapTokenReq struct {
	g.Meta      `path:"/defi/swap" method:"post" tags:"DeFi" summary:"代币兑换"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	FromToken   string `v:"required" dc:"支付代币地址"`
	ToToken     string `v:"required" dc:"获得代币地址"`
	Amount      string `v:"required" dc:"兑换数量"`
	FromAddress string `v:"required" dc:"支付地址"`
	SlippageBps int    `d:"30" dc:"滑点(万分之)"`
	Type        string `d:"EXACT_INPUT" dc:"类型(EXACT_INPUT/EXACT_OUTPUT)"`
}

type SwapTokenRes struct {
	Hash      string `json:"hash" dc:"交易哈希"`
	AmountOut string `json:"amountOut" dc:"获得数量"`
}

// AddLiquidityReq 添加流动性请求
type AddLiquidityReq struct {
	g.Meta      `path:"/defi/liquidity/add" method:"post" tags:"DeFi" summary:"添加流动性"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	TokenA      string `v:"required" dc:"代币A地址"`
	TokenB      string `v:"required" dc:"代币B地址"`
	AmountA     string `v:"required" dc:"代币A数量"`
	AmountB     string `v:"required" dc:"代币B数量"`
	FromAddress string `v:"required" dc:"支付地址"`
	SlippageBps int    `d:"30" dc:"滑点(万分之)"`
}

type AddLiquidityRes struct {
	Hash      string `json:"hash" dc:"交易哈希"`
	Liquidity string `json:"liquidity" dc:"LP代币数量"`
}

// RemoveLiquidityReq 移除流动性请求
type RemoveLiquidityReq struct {
	g.Meta      `path:"/defi/liquidity/remove" method:"post" tags:"DeFi" summary:"移除流动性"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	Pair        string `v:"required" dc:"交易对地址"`
	Liquidity   string `v:"required" dc:"LP代币数量"`
	FromAddress string `v:"required" dc:"地址"`
	SlippageBps int    `d:"30" dc:"滑点(万分之)"`
}

type RemoveLiquidityRes struct {
	Hash    string `json:"hash" dc:"交易哈希"`
	Amount0 string `json:"amount0" dc:"代币0数量"`
	Amount1 string `json:"amount1" dc:"代币1数量"`
}

// SupplyReq 存款请求
type SupplyReq struct {
	g.Meta      `path:"/defi/lending/supply" method:"post" tags:"DeFi" summary:"存款"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	Pool        string `v:"required" dc:"借贷池地址"`
	Token       string `v:"required" dc:"代币地址"`
	Amount      string `v:"required" dc:"数量"`
	FromAddress string `v:"required" dc:"地址"`
}

type SupplyRes struct {
	Hash string `json:"hash" dc:"交易哈希"`
}

// BorrowReq 借款请求
type BorrowReq struct {
	g.Meta      `path:"/defi/lending/borrow" method:"post" tags:"DeFi" summary:"借款"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	Pool        string `v:"required" dc:"借贷池地址"`
	Token       string `v:"required" dc:"代币地址"`
	Amount      string `v:"required" dc:"数量"`
	RateMode    int    `d:"2" dc:"利率模式 1:稳定 2:浮动"`
	FromAddress string `v:"required" dc:"地址"`
}

type BorrowRes struct {
	Hash string `json:"hash" dc:"交易哈希"`
}

// RepayReq 还款请求
type RepayReq struct {
	g.Meta      `path:"/defi/lending/repay" method:"post" tags:"DeFi" summary:"还款"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	Pool        string `v:"required" dc:"借贷池地址"`
	Token       string `v:"required" dc:"代币地址"`
	Amount      string `v:"required" dc:"数量"`
	RateMode    int    `d:"2" dc:"利率模式 1:稳定 2:浮动"`
	FromAddress string `v:"required" dc:"地址"`
}

type RepayRes struct {
	Hash string `json:"hash" dc:"交易哈希"`
}

// StakeReq 质押请求
type StakeReq struct {
	g.Meta      `path:"/defi/farm/stake" method:"post" tags:"DeFi" summary:"质押"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	Pool        string `v:"required" dc:"矿池地址"`
	Amount      string `v:"required" dc:"数量"`
	FromAddress string `v:"required" dc:"地址"`
}

type StakeRes struct {
	Hash string `json:"hash" dc:"交易哈希"`
}

// UnstakeReq 解质押请求
type UnstakeReq struct {
	g.Meta      `path:"/defi/farm/unstake" method:"post" tags:"DeFi" summary:"解质押"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	Pool        string `v:"required" dc:"矿池地址"`
	Amount      string `v:"required" dc:"数量"`
	FromAddress string `v:"required" dc:"地址"`
}

type UnstakeRes struct {
	Hash string `json:"hash" dc:"交易哈希"`
}

// ClaimRewardReq 领取奖励请求
type ClaimRewardReq struct {
	g.Meta      `path:"/defi/farm/claim" method:"post" tags:"DeFi" summary:"领取奖励"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	Pool        string `v:"required" dc:"矿池地址"`
	FromAddress string `v:"required" dc:"地址"`
}

type ClaimRewardRes struct {
	Hash   string `json:"hash" dc:"交易哈希"`
	Reward string `json:"reward" dc:"奖励数量"`
}

// DepositVaultReq 存入机枪池请求
type DepositVaultReq struct {
	g.Meta      `path:"/defi/vault/deposit" method:"post" tags:"DeFi" summary:"存入机枪池"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	Vault       string `v:"required" dc:"机枪池地址"`
	Amount      string `v:"required" dc:"数量"`
	FromAddress string `v:"required" dc:"地址"`
}

type DepositVaultRes struct {
	Hash   string `json:"hash" dc:"交易哈希"`
	Shares string `json:"shares" dc:"份额数量"`
}

// WithdrawVaultReq 提取机枪池请求
type WithdrawVaultReq struct {
	g.Meta      `path:"/defi/vault/withdraw" method:"post" tags:"DeFi" summary:"提取机枪池"`
	ChainId     uint64 `v:"required" dc:"链ID"`
	Vault       string `v:"required" dc:"机枪池地址"`
	Shares      string `v:"required" dc:"份额数量"`
	FromAddress string `v:"required" dc:"地址"`
}

type WithdrawVaultRes struct {
	Hash   string `json:"hash" dc:"交易哈希"`
	Amount string `json:"amount" dc:"提取数量"`
}
