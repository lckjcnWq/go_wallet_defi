package model

// Contract 合约信息
type Contract struct {
	Id           uint64 `json:"id"`           // 合约ID
	Name         string `json:"name"`         // 合约名称
	Address      string `json:"address"`      // 合约地址
	ABI          string `json:"abi"`          // 合约ABI
	Bytecode     string `json:"bytecode"`     // 合约字节码
	Network      string `json:"network"`      // 网络
	DeployHash   string `json:"deployHash"`   // 部署交易哈希
	DeployHeight int64  `json:"deployHeight"` // 部署区块高度
	Creator      string `json:"creator"`      // 创建者地址
	Status       int    `json:"status"`       // 状态 1:正常 0:禁用
	CreatedAt    int64  `json:"createdAt"`    // 创建时间
	UpdatedAt    int64  `json:"updatedAt"`    // 更新时间
}

// ContractEvent 合约事件
type ContractEvent struct {
	Id          uint64 `json:"id"`          // 事件ID
	ContractId  uint64 `json:"contractId"`  // 合约ID
	Name        string `json:"name"`        // 事件名称
	Signature   string `json:"signature"`   // 事件签名
	Topics      string `json:"topics"`      // 事件topics
	Data        string `json:"data"`        // 事件数据
	BlockNumber int64  `json:"blockNumber"` // 区块高度
	BlockHash   string `json:"blockHash"`   // 区块哈希
	TxHash      string `json:"txHash"`      // 交易哈希
	TxIndex     int    `json:"txIndex"`     // 交易索引
	LogIndex    int    `json:"logIndex"`    // 日志索引
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
}

// ContractCall 合约调用记录
type ContractCall struct {
	Id         uint64 `json:"id"`         // 调用ID
	ContractId uint64 `json:"contractId"` // 合约ID
	Method     string `json:"method"`     // 方法名称
	Params     string `json:"params"`     // 调用参数
	From       string `json:"from"`       // 调用地址
	Value      string `json:"value"`      // 调用金额
	Hash       string `json:"hash"`       // 交易哈希
	Status     int    `json:"status"`     // 状态
	Error      string `json:"error"`      // 错误信息
	CreatedAt  int64  `json:"createdAt"`  // 创建时间
	UpdatedAt  int64  `json:"updatedAt"`  // 更新时间
}
