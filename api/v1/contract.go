package v1

import "github.com/gogf/gf/v2/frame/g"

// DeployContractReq 部署合约
type DeployContractReq struct {
	g.Meta   `path:"/contract/deploy" method:"post" tags:"合约管理" summary:"部署合约"`
	Name     string `v:"required" dc:"合约名称"`
	ABI      string `v:"required" dc:"合约ABI"`
	Bytecode string `v:"required" dc:"合约字节码"`
	From     string `v:"required" dc:"部署地址"`
	Args     string `dc:"构造参数(JSON)"`
	Network  string `dc:"网络,默认ethereum"`
}

type DeployContractRes struct {
	Hash    string `json:"hash" dc:"交易哈希"`
	Address string `json:"address" dc:"合约地址"`
}

// CallContractReq 调用合约
type CallContractReq struct {
	g.Meta  `path:"/contract/call" method:"post" tags:"合约管理" summary:"调用合约"`
	Address string `v:"required" dc:"合约地址"`
	Method  string `v:"required" dc:"方法名称"`
	Args    string `dc:"调用参数(JSON)"`
	From    string `v:"required" dc:"调用地址"`
	Value   string `dc:"调用金额(ETH)"`
}

type CallContractRes struct {
	Hash   string      `json:"hash,omitempty" dc:"交易哈希(写入方法)"`
	Result interface{} `json:"result,omitempty" dc:"返回结果(读取方法)"`
}

// GetContractEventsReq 获取合约事件
type GetContractEventsReq struct {
	g.Meta    `path:"/contract/events" method:"get" tags:"合约管理" summary:"获取合约事件"`
	Address   string `v:"required" dc:"合约地址"`
	EventName string `dc:"事件名称"`
	FromBlock int64  `dc:"起始区块"`
	ToBlock   int64  `dc:"结束区块"`
	Page      int    `d:"1" dc:"页码"`
	PageSize  int    `d:"10" dc:"每页数量"`
}

type GetContractEventsRes struct {
	List  []ContractEventInfo `json:"list" dc:"事件列表"`
	Total int                 `json:"total" dc:"总数"`
}

type ContractEventInfo struct {
	Name        string      `json:"name" dc:"事件名称"`
	Data        interface{} `json:"data" dc:"事件数据"`
	BlockNumber int64       `json:"blockNumber" dc:"区块高度"`
	TxHash      string      `json:"txHash" dc:"交易哈希"`
	CreatedAt   int64       `json:"createdAt" dc:"创建时间"`
}
