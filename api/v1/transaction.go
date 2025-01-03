package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 转账ETH
type TransferEthReq struct {
	g.Meta   `path:"/transaction/transfer-eth" method:"post" tags:"交易管理" summary:"ETH转账"`
	From     string `v:"required" dc:"发送地址"`
	To       string `v:"required" dc:"接收地址"`
	Amount   string `v:"required" dc:"转账金额(ETH)"`
	GasPrice string `dc:"gas价格(gwei),默认建议价格"`
	GasLimit uint64 `dc:"gas限制,默认21000"`
}

type TransferEthRes struct {
	Hash string `json:"hash" dc:"交易哈希"`
}

// 转账代币
type TransferTokenReq struct {
	g.Meta   `path:"/transaction/transfer-token" method:"post" tags:"交易管理" summary:"代币转账"`
	From     string `v:"required" dc:"发送地址"`
	To       string `v:"required" dc:"接收地址"`
	Amount   string `v:"required" dc:"转账金额"`
	Token    string `v:"required" dc:"代币合约地址"`
	GasPrice string `dc:"gas价格(gwei)"`
	GasLimit uint64 `dc:"gas限制"`
}

type TransferTokenRes struct {
	Hash string `json:"hash" dc:"交易哈希"`
}

// 获取交易记录
type GetTransactionsReq struct {
	g.Meta   `path:"/transaction/list" method:"get" tags:"交易管理" summary:"交易记录"`
	Address  string `v:"required" dc:"钱包地址"`
	Page     int    `d:"1" dc:"页码"`
	PageSize int    `d:"10" dc:"每页数量"`
}

type GetTransactionsRes struct {
	List  []TransactionInfo `json:"list" dc:"交易列表"`
	Total int               `json:"total" dc:"总数"`
}

type TransactionInfo struct {
	Hash        string `json:"hash" dc:"交易哈希"`
	From        string `json:"from" dc:"发送地址"`
	To          string `json:"to" dc:"接收地址"`
	Amount      string `json:"amount" dc:"转账金额"`
	Token       string `json:"token" dc:"代币地址"`
	Status      int    `json:"status" dc:"状态 0:待处理 1:已确认 2:失败"`
	BlockNumber int64  `json:"blockNumber" dc:"区块高度"`
	BlockTime   int64  `json:"blockTime" dc:"区块时间"`
	CreatedAt   int64  `json:"createdAt" dc:"创建时间"`
}
