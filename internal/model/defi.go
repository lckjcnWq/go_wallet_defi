package model

// DexTrade DEX交易记录
type DexTrade struct {
	Id         uint64 `json:"id"`         // ID
	ChainId    uint64 `json:"chainId"`    // 链ID
	FromToken  string `json:"fromToken"`  // 支付代币
	ToToken    string `json:"toToken"`    // 获得代币
	FromAmount string `json:"fromAmount"` // 支付金额
	ToAmount   string `json:"toAmount"`   // 获得金额
	User       string `json:"user"`       // 用户地址
	Router     string `json:"router"`     // 路由合约
	Path       string `json:"path"`       // 兑换路径
	Type       string `json:"type"`       // 类型
	Hash       string `json:"hash"`       // 交易哈希
	Status     int    `json:"status"`     // 状态
	Error      string `json:"error"`      // 错误信息
	CreatedAt  int64  `json:"createdAt"`  // 创建时间
	UpdatedAt  int64  `json:"updatedAt"`  // 更新时间
}

// Liquidity 流动性记录
type Liquidity struct {
	Id        uint64 `json:"id"`        // ID
	ChainId   uint64 `json:"chainId"`   // 链ID
	Pair      string `json:"pair"`      // 交易对地址
	Token0    string `json:"token0"`    // 代币0
	Token1    string `json:"token1"`    // 代币1
	Amount0   string `json:"amount0"`   // 数量0
	Amount1   string `json:"amount1"`   // 数量1
	Liquidity string `json:"liquidity"` // LP数量
	User      string `json:"user"`      // 用户地址
	Type      string `json:"type"`      // 类型
	Hash      string `json:"hash"`      // 交易哈希
	Status    int    `json:"status"`    // 状态
	Error     string `json:"error"`     // 错误信息
	CreatedAt int64  `json:"createdAt"` // 创建时间
	UpdatedAt int64  `json:"updatedAt"` // 更新时间
}

// Lending 借贷记录
type Lending struct {
	Id        uint64 `json:"id"`        // ID
	ChainId   uint64 `json:"chainId"`   // 链ID
	Pool      string `json:"pool"`      // 借贷池
	Token     string `json:"token"`     // 代币地址
	Amount    string `json:"amount"`    // 数量
	User      string `json:"user"`      // 用户地址
	Type      string `json:"type"`      // 类型
	RateMode  int    `json:"rateMode"`  // 利率模式
	Hash      string `json:"hash"`      // 交易哈希
	Status    int    `json:"status"`    // 状态
	Error     string `json:"error"`     // 错误信息
	CreatedAt int64  `json:"createdAt"` // 创建时间
	UpdatedAt int64  `json:"updatedAt"` // 更新时间
}

// YieldFarm 收益农场
type YieldFarm struct {
	Id          uint64 `json:"id"`          // ID
	ChainId     uint64 `json:"chainId"`     // 链ID
	Pool        string `json:"pool"`        // 矿池地址
	StakeToken  string `json:"stakeToken"`  // 质押代币
	RewardToken string `json:"rewardToken"` // 奖励代币
	Amount      string `json:"amount"`      // 质押数量
	Reward      string `json:"reward"`      // 奖励数量
	User        string `json:"user"`        // 用户地址
	Type        string `json:"type"`        // 类型
	Hash        string `json:"hash"`        // 交易哈希
	Status      int    `json:"status"`      // 状态
	Error       string `json:"error"`       // 错误信息
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
	UpdatedAt   int64  `json:"updatedAt"`   // 更新时间
}

// Vault 收益聚合
type Vault struct {
	Id        uint64 `json:"id"`        // ID
	ChainId   uint64 `json:"chainId"`   // 链ID
	Vault     string `json:"vault"`     // 机枪池地址
	Token     string `json:"token"`     // 存入代币
	Amount    string `json:"amount"`    // 存入数量
	Shares    string `json:"shares"`    // 份额数量
	User      string `json:"user"`      // 用户地址
	Type      string `json:"type"`      // 类型
	Hash      string `json:"hash"`      // 交易哈希
	Status    int    `json:"status"`    // 状态
	Error     string `json:"error"`     // 错误信息
	CreatedAt int64  `json:"createdAt"` // 创建时间
	UpdatedAt int64  `json:"updatedAt"` // 更新时间
}
