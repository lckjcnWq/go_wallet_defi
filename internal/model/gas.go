package model

import "github.com/gogf/gf/v2/os/gtime"

// BatchTransaction 批量交易
type BatchTransaction struct {
	Id         uint64      `json:"id"           description:"ID"`
	BatchId    string      `json:"batch_id"     description:"批次ID"`
	ChainId    uint64      `json:"chain_id"     description:"链ID"`
	TxHash     string      `json:"tx_hash"      description:"交易哈希"`
	From       string      `json:"from"         description:"发送地址"`
	Calls      string      `json:"calls"        description:"调用数据JSON"`
	Status     int         `json:"status"       description:"状态 0:待执行 1:已执行 2:失败"`
	GasUsed    string      `json:"gas_used"     description:"Gas使用量"`
	GasPrice   string      `json:"gas_price"    description:"Gas价格"`
	CreatedAt  *gtime.Time `json:"created_at"   description:"创建时间"`
	ExecutedAt *gtime.Time `json:"executed_at"  description:"执行时间"`
}

// GasPrice Gas价格记录
type GasPrice struct {
	Id          uint64      `json:"id"           description:"ID"`
	ChainId     uint64      `json:"chain_id"     description:"链ID"`
	BlockNumber uint64      `json:"block_number" description:"区块号"`
	BaseFee     string      `json:"base_fee"     description:"基础费用"`
	Priority    string      `json:"priority"     description:"优先级费用"`
	Average     string      `json:"average"      description:"平均Gas价格"`
	CreatedAt   *gtime.Time `json:"created_at"   description:"创建时间"`
}

// GasStrategy Gas策略
type GasStrategy struct {
	Id           uint64      `json:"id"           description:"ID"`
	ChainId      uint64      `json:"chain_id"     description:"链ID"`
	Type         string      `json:"type"         description:"类型 FASTEST/FAST/STANDARD/SLOW"`
	BaseFee      string      `json:"base_fee"     description:"基础费用"`
	Priority     string      `json:"priority"     description:"优先级费用"`
	EstimateWait int         `json:"estimate_wait" description:"预估等待时间(秒)"`
	UpdatedAt    *gtime.Time `json:"updated_at"   description:"更新时间"`
}

// MEVProtection MEV防护记录
type MEVProtection struct {
	Id          uint64      `json:"id"           description:"ID"`
	TxHash      string      `json:"tx_hash"      description:"交易哈希"`
	Type        string      `json:"type"         description:"类型 SANDWICH/FRONTRUN"`
	From        string      `json:"from"         description:"攻击地址"`
	To          string      `json:"to"           description:"目标地址"`
	BlockNumber uint64      `json:"block_number" description:"区块号"`
	GasPrice    string      `json:"gas_price"    description:"Gas价格"`
	CreatedAt   *gtime.Time `json:"created_at"   description:"创建时间"`
}

// TxAcceleration 交易加速
type TxAcceleration struct {
	Id          uint64      `json:"id"           description:"ID"`
	TxHash      string      `json:"tx_hash"      description:"原始交易哈希"`
	NewTxHash   string      `json:"new_tx_hash"  description:"新交易哈希"`
	ChainId     uint64      `json:"chain_id"     description:"链ID"`
	From        string      `json:"from"         description:"发送地址"`
	OldGasPrice string      `json:"old_gas_price" description:"原Gas价格"`
	NewGasPrice string      `json:"new_gas_price" description:"新Gas价格"`
	Status      int         `json:"status"       description:"状态 0:待执行 1:已完成 2:失败"`
	CreatedAt   *gtime.Time `json:"created_at"   description:"创建时间"`
	ExecutedAt  *gtime.Time `json:"executed_at"  description:"执行时间"`
}
