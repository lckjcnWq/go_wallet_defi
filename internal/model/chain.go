package model

type Chain struct {
	Id            uint64 `json:"id"`            // 链ID
	Name          string `json:"name"`          // 链名称
	ChainId       uint64 `json:"chainId"`       // 链ID
	Symbol        string `json:"symbol"`        // 原生代币符号
	Decimals      int    `json:"decimals"`      // 原生代币精度
	ExplorerUrl   string `json:"explorerUrl"`   // 浏览器地址
	RpcUrls       string `json:"rpcUrls"`       // RPC节点地址
	BridgeAddress string `json:"bridgeAddress"` // 跨链桥合约地址
	Status        int    `json:"status"`        // 状态
	CreatedAt     int64  `json:"createdAt"`     // 创建时间
	UpdatedAt     int64  `json:"updatedAt"`     // 更新时间
}

// ContractMapping 合约地址映射
type ContractMapping struct {
	Id          uint64 `json:"id"`          // ID
	Name        string `json:"name"`        // 合约名称
	FromChainId uint64 `json:"fromChainId"` // 来源链ID
	FromAddress string `json:"fromAddress"` // 来源合约地址
	ToChainId   uint64 `json:"toChainId"`   // 目标链ID
	ToAddress   string `json:"toAddress"`   // 目标合约地址
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
	UpdatedAt   int64  `json:"updatedAt"`   // 更新时间
}

// CrossTransfer 跨链交易
type CrossTransfer struct {
	Id           uint64 `json:"id"`           // ID
	FromChainId  uint64 `json:"fromChainId"`  // 来源链ID
	ToChainId    uint64 `json:"toChainId"`    // 目标链ID
	FromAddress  string `json:"fromAddress"`  // 来源地址
	ToAddress    string `json:"toAddress"`    // 目标地址
	TokenAddress string `json:"tokenAddress"` // 代币地址
	Amount       string `json:"amount"`       // 金额
	Fee          string `json:"fee"`          // 手续费
	Nonce        uint64 `json:"nonce"`        // 交易序号
	FromHash     string `json:"fromHash"`     // 来源链交易哈希
	ToHash       string `json:"toHash"`       // 目标链交易哈希
	Status       int    `json:"status"`       // 状态
	Error        string `json:"error"`        // 错误信息
	CreatedAt    int64  `json:"createdAt"`    // 创建时间
	UpdatedAt    int64  `json:"updatedAt"`    // 更新时间
}
