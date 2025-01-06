package v1

import "github.com/gogf/gf/v2/frame/g"

// CrossTransferReq 跨链转账请求
type CrossTransferReq struct {
	g.Meta       `path:"/bridge/transfer" method:"post" tags:"跨链桥" summary:"跨链转账"`
	FromChainId  uint64 `v:"required" dc:"来源链ID"`
	ToChainId    uint64 `v:"required" dc:"目标链ID"`
	FromAddress  string `v:"required" dc:"来源地址"`
	ToAddress    string `v:"required" dc:"目标地址"`
	TokenAddress string `dc:"代币地址(空表示原生代币)"`
	Amount       string `v:"required" dc:"金额"`
}

type CrossTransferRes struct {
	Hash  string `json:"hash" dc:"交易哈希"`
	Nonce uint64 `json:"nonce" dc:"交易序号"`
}

// GetCrossTransferReq 获取跨链交易请求
type GetCrossTransferReq struct {
	g.Meta      `path:"/bridge/transfer" method:"get" tags:"跨链桥" summary:"获取跨链交易"`
	FromChainId uint64 `dc:"来源链ID"`
	ToChainId   uint64 `dc:"目标链ID"`
	Address     string `dc:"用户地址"`
	Status      int    `dc:"状态"`
	Page        int    `d:"1" dc:"页码"`
	PageSize    int    `d:"10" dc:"每页数量"`
}

type GetCrossTransferRes struct {
	List  []CrossTransferInfo `json:"list" dc:"交易列表"`
	Total int                 `json:"total" dc:"总数"`
}

type CrossTransferInfo struct {
	Id           uint64 `json:"id" dc:"ID"`
	FromChainId  uint64 `json:"fromChainId" dc:"来源链ID"`
	ToChainId    uint64 `json:"toChainId" dc:"目标链ID"`
	FromAddress  string `json:"fromAddress" dc:"来源地址"`
	ToAddress    string `json:"toAddress" dc:"目标地址"`
	TokenAddress string `json:"tokenAddress" dc:"代币地址"`
	Amount       string `json:"amount" dc:"金额"`
	Fee          string `json:"fee" dc:"手续费"`
	FromHash     string `json:"fromHash" dc:"来源链交易哈希"`
	ToHash       string `json:"toHash" dc:"目标链交易哈希"`
	Status       int    `json:"status" dc:"状态"`
	CreatedAt    int64  `json:"createdAt" dc:"创建时间"`
}
