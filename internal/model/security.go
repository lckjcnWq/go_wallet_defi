package model

import "github.com/gogf/gf/v2/os/gtime"

// MultiSigWallet 多签钱包
type MultiSigWallet struct {
	Id        uint64      `json:"id"           description:"ID"`
	ChainId   uint64      `json:"chain_id"     description:"链ID"`
	Address   string      `json:"address"      description:"合约地址"`
	Name      string      `json:"name"         description:"钱包名称"`
	Owners    string      `json:"owners"       description:"所有者地址列表JSON"`
	Threshold int         `json:"threshold"    description:"签名阈值"`
	CreatedBy string      `json:"created_by"   description:"创建者"`
	Status    int         `json:"status"       description:"状态 1:正常 2:禁用"`
	CreatedAt *gtime.Time `json:"created_at"   description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updated_at"   description:"更新时间"`
}

// MultiSigTransaction 多签交易
type MultiSigTransaction struct {
	Id            uint64      `json:"id"               description:"ID"`
	WalletId      uint64      `json:"wallet_id"        description:"多签钱包ID"`
	ChainId       uint64      `json:"chain_id"         description:"链ID"`
	TxHash        string      `json:"tx_hash"          description:"交易哈希"`
	To            string      `json:"to"               description:"接收地址"`
	Value         string      `json:"value"            description:"金额"`
	Data          string      `json:"data"             description:"数据"`
	Description   string      `json:"description"      description:"交易描述"`
	Confirmations string      `json:"confirmations"    description:"确认签名JSON"`
	ExecutedAt    *gtime.Time `json:"executed_at"      description:"执行时间"`
	Status        int         `json:"status"           description:"状态 0:待确认 1:已确认 2:已执行 3:已拒绝"`
	CreatedAt     *gtime.Time `json:"created_at"       description:"创建时间"`
	UpdatedAt     *gtime.Time `json:"updated_at"       description:"更新时间"`
}

// Whitelist 白名单
type Whitelist struct {
	Id        uint64      `json:"id"           description:"ID"`
	UserId    uint64      `json:"user_id"      description:"用户ID"`
	Address   string      `json:"address"      description:"地址"`
	Name      string      `json:"name"         description:"名称"`
	Type      int         `json:"type"         description:"类型 1:转账 2:合约调用"`
	ExpireAt  *gtime.Time `json:"expire_at"    description:"过期时间"`
	Status    int         `json:"status"       description:"状态 1:正常 2:禁用"`
	CreatedAt *gtime.Time `json:"created_at"   description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updated_at"   description:"更新时间"`
}

// TransactionLimit 交易限额
type TransactionLimit struct {
	Id           uint64      `json:"id"               description:"ID"`
	UserId       uint64      `json:"user_id"          description:"用户ID"`
	TokenAddress string      `json:"token_address"    description:"代币地址"`
	SingleLimit  string      `json:"single_limit"     description:"单笔限额"`
	DailyLimit   string      `json:"daily_limit"      description:"日限额"`
	WeeklyLimit  string      `json:"weekly_limit"     description:"周限额"`
	MonthlyLimit string      `json:"monthly_limit"    description:"月限额"`
	Status       int         `json:"status"           description:"状态 1:正常 2:禁用"`
	CreatedAt    *gtime.Time `json:"created_at"       description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updated_at"       description:"更新时间"`
}

// RiskRule 风控规则
type RiskRule struct {
	Id        uint64      `json:"id"           description:"ID"`
	Name      string      `json:"name"         description:"规则名称"`
	Type      int         `json:"type"         description:"类型 1:地址黑名单 2:金额预警 3:频率限制"`
	Content   string      `json:"content"      description:"规则内容JSON"`
	Action    int         `json:"action"       description:"动作 1:拒绝 2:预警 3:需人工审核"`
	Status    int         `json:"status"       description:"状态 1:正常 2:禁用"`
	CreatedAt *gtime.Time `json:"created_at"   description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updated_at"   description:"更新时间"`
}

// RiskLog 风控日志
type RiskLog struct {
	Id        uint64      `json:"id"           description:"ID"`
	UserId    uint64      `json:"user_id"      description:"用户ID"`
	RuleId    uint64      `json:"rule_id"      description:"规则ID"`
	Type      int         `json:"type"         description:"类型"`
	Content   string      `json:"content"      description:"内容"`
	Result    int         `json:"result"       description:"结果 1:通过 2:拒绝 3:预警"`
	CreatedAt *gtime.Time `json:"created_at"   description:"创建时间"`
}
