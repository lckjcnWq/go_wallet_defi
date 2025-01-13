package model

import "github.com/gogf/gf/v2/os/gtime"

// TransactionAnalysis 交易分析
type TransactionAnalysis struct {
	Id           uint64      `json:"id"               description:"ID"`
	Address      string      `json:"address"          description:"地址"`
	TotalTxCount int64       `json:"total_tx_count"   description:"总交易次数"`
	TotalVolume  string      `json:"total_volume"     description:"总交易量"`
	AvgGasUsed   int64       `json:"avg_gas_used"     description:"平均gas消耗"`
	TotalGasFee  string      `json:"total_gas_fee"    description:"总gas费用"`
	FirstTxTime  *gtime.Time `json:"first_tx_time"    description:"首次交易时间"`
	LastTxTime   *gtime.Time `json:"last_tx_time"     description:"最后交易时间"`
	UpdatedAt    *gtime.Time `json:"updated_at"       description:"更新时间"`
}

// AssetAnalysis 资产分析
type AssetAnalysis struct {
	Id           uint64      `json:"id"               description:"ID"`
	Address      string      `json:"address"          description:"地址"`
	TokenAddress string      `json:"token_address"    description:"代币地址"`
	TokenType    string      `json:"token_type"       description:"代币类型 ERC20/ERC721"`
	Balance      string      `json:"balance"          description:"余额"`
	Value        string      `json:"value"            description:"估值"`
	CostBasis    string      `json:"cost_basis"       description:"成本"`
	UpdatedAt    *gtime.Time `json:"updated_at"       description:"更新时间"`
}

// ProfitAnalysis 收益分析
type ProfitAnalysis struct {
	Id        uint64      `json:"id"               description:"ID"`
	Address   string      `json:"address"          description:"地址"`
	Type      string      `json:"type"             description:"类型 DEFI/STAKING/NFT"`
	Platform  string      `json:"platform"         description:"平台"`
	Principal string      `json:"principal"        description:"本金"`
	Profit    string      `json:"profit"           description:"收益"`
	Apy       float64     `json:"apy"              description:"年化收益率"`
	StartTime *gtime.Time `json:"start_time"       description:"开始时间"`
	UpdatedAt *gtime.Time `json:"updated_at"       description:"更新时间"`
}

// RiskAnalysis 风险分析
type RiskAnalysis struct {
	Id          uint64      `json:"id"               description:"ID"`
	Address     string      `json:"address"          description:"地址"`
	Type        string      `json:"type"             description:"类型 TX/ADDRESS/ASSET/MARKET"`
	RiskLevel   int         `json:"risk_level"       description:"风险等级 1-5"`
	RiskFactor  string      `json:"risk_factor"      description:"风险因子"`
	Description string      `json:"description"      description:"描述"`
	CreatedAt   *gtime.Time `json:"created_at"       description:"创建时间"`
}

// TxPattern 交易模式
type TxPattern struct {
	Id         uint64      `json:"id"               description:"ID"`
	Address    string      `json:"address"          description:"地址"`
	Pattern    string      `json:"pattern"          description:"模式类型"`
	Frequency  int64       `json:"frequency"        description:"频率"`
	AvgAmount  string      `json:"avg_amount"       description:"平均金额"`
	TimeRange  string      `json:"time_range"       description:"时间范围"`
	Confidence float64     `json:"confidence"       description:"置信度"`
	UpdatedAt  *gtime.Time `json:"updated_at"       description:"更新时间"`
}

// MarketIndicator 市场指标
type MarketIndicator struct {
	Id             uint64      `json:"id"               description:"ID"`
	Symbol         string      `json:"symbol"           description:"交易对"`
	Price          string      `json:"price"            description:"价格"`
	Volume24h      string      `json:"volume_24h"       description:"24h成交量"`
	PriceChange24h float64     `json:"price_change_24h" description:"24h价格变化"`
	Volatility     float64     `json:"volatility"       description:"波动率"`
	MarketCap      string      `json:"market_cap"       description:"市值"`
	UpdatedAt      *gtime.Time `json:"updated_at"       description:"更新时间"`
}

// AssetDistribution 资产分布
type AssetDistribution struct {
	Id         uint64      `json:"id"               description:"ID"`
	Address    string      `json:"address"          description:"地址"`
	AssetType  string      `json:"asset_type"       description:"资产类型"`
	Percentage float64     `json:"percentage"       description:"占比"`
	Value      string      `json:"value"            description:"价值"`
	UpdatedAt  *gtime.Time `json:"updated_at"       description:"更新时间"`
}
