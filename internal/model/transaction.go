package model

// Transaction 交易记录
type Transaction struct {
	Id           uint64 `json:"id" dc:"交易ID"`
	UserId       uint64 `json:"userId" dc:"用户ID"`
	FromAddress  string `json:"fromAddress" dc:"发送地址"`
	ToAddress    string `json:"toAddress" dc:"接收地址"`
	Amount       string `json:"amount" dc:"转账金额"`
	TokenAddress string `json:"tokenAddress" dc:"代币合约地址"`
	Hash         string `json:"hash" dc:"交易哈希"`
	Nonce        uint64 `json:"nonce" dc:"交易nonce"`
	GasPrice     string `json:"gasPrice" dc:"gas价格(wei)"`
	GasLimit     uint64 `json:"gasLimit" dc:"gas限制"`
	Data         string `json:"data" dc:"交易数据"`
	Status       int    `json:"status" dc:"状态 0:待处理 1:已确认 2:失败"`
	BlockNumber  int64  `json:"blockNumber" dc:"区块高度"`
	BlockTime    int64  `json:"blockTime" dc:"区块时间"`
	CreatedAt    int64  `json:"createdAt" dc:"创建时间"`
	UpdatedAt    int64  `json:"updatedAt" dc:"更新时间"`
}

// TransactionStatus 交易状态
const (
	TransactionStatusPending = 0 // 待处理
	TransactionStatusSuccess = 1 // 已确认
	TransactionStatusFailed  = 2 // 失败
)

// GetStatusDesc 获取状态描述
func (t *Transaction) GetStatusDesc() string {
	switch t.Status {
	case TransactionStatusPending:
		return "待处理"
	case TransactionStatusSuccess:
		return "已确认"
	case TransactionStatusFailed:
		return "失败"
	default:
		return "未知"
	}
}
