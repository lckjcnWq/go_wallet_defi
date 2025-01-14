package model

import "github.com/gogf/gf/v2/os/gtime"

// SwapTransaction DEX交易
type SwapTransaction struct {
	Id         uint64      `json:"id"               description:"ID"`
	ChainId    uint64      `json:"chain_id"         description:"链ID"`
	Protocol   string      `json:"protocol"         description:"协议 UNISWAP/SUSHISWAP"`
	TxHash     string      `json:"tx_hash"          description:"交易哈希"`
	FromToken  string      `json:"from_token"       description:"支付代币"`
	ToToken    string      `json:"to_token"         description:"接收代币"`
	FromAmount string      `json:"from_amount"      description:"支付金额"`
	ToAmount   string      `json:"to_amount"        description:"接收金额"`
	Sender     string      `json:"sender"           description:"发送地址"`
	Receiver   string      `json:"receiver"         description:"接收地址"`
	Status     int         `json:"status"           description:"状态 0:待执行 1:已完成 2:失败"`
	CreatedAt  *gtime.Time `json:"created_at"       description:"创建时间"`
	ExecutedAt *gtime.Time `json:"executed_at"      description:"执行时间"`
}

// LendingPosition 借贷仓位
type LendingPosition struct {
	Id               uint64      `json:"id"               description:"ID"`
	ChainId          uint64      `json:"chain_id"         description:"链ID"`
	Protocol         string      `json:"protocol"         description:"协议 AAVE/COMPOUND"`
	Address          string      `json:"address"          description:"地址"`
	Token            string      `json:"token"            description:"代币地址"`
	SupplyAmount     string      `json:"supply_amount"    description:"存款金额"`
	BorrowAmount     string      `json:"borrow_amount"    description:"借款金额"`
	CollateralFactor string      `json:"collateral_factor" description:"抵押率"`
	HealthFactor     string      `json:"health_factor"    description:"健康因子"`
	SupplyRate       string      `json:"supply_rate"      description:"存款利率"`
	BorrowRate       string      `json:"borrow_rate"      description:"借款利率"`
	UpdatedAt        *gtime.Time `json:"updated_at"       description:"更新时间"`
}

// NFTTransaction NFT交易
type NFTTransaction struct {
	Id              uint64      `json:"id"               description:"ID"`
	ChainId         uint64      `json:"chain_id"         description:"链ID"`
	Protocol        string      `json:"protocol"         description:"协议 OPENSEA/LOOKSRARE"`
	TxHash          string      `json:"tx_hash"          description:"交易哈希"`
	ContractAddress string      `json:"contract_address" description:"NFT合约地址"`
	TokenId         string      `json:"token_id"         description:"Token ID"`
	FromAddress     string      `json:"from_address"     description:"卖方地址"`
	ToAddress       string      `json:"to_address"       description:"买方地址"`
	Price           string      `json:"price"            description:"价格"`
	PayToken        string      `json:"pay_token"        description:"支付代币"`
	Type            string      `json:"type"             description:"类型 SALE/MINT"`
	Status          int         `json:"status"           description:"状态 0:待执行 1:已完成 2:失败"`
	CreatedAt       *gtime.Time `json:"created_at"       description:"创建时间"`
	ExecutedAt      *gtime.Time `json:"executed_at"      description:"执行时间"`
}

// BridgeTransaction 跨链交易
type BridgeTransaction struct {
	Id          uint64      `json:"id"               description:"ID"`
	Protocol    string      `json:"protocol"         description:"协议 MULTICHAIN/LAYERBRIDGE"`
	FromChainId uint64      `json:"from_chain_id"    description:"源链ID"`
	ToChainId   uint64      `json:"to_chain_id"      description:"目标链ID"`
	FromTxHash  string      `json:"from_tx_hash"     description:"源交易哈希"`
	ToTxHash    string      `json:"to_tx_hash"       description:"目标交易哈希"`
	Token       string      `json:"token"            description:"代币地址"`
	Amount      string      `json:"amount"           description:"金额"`
	FromAddress string      `json:"from_address"     description:"发送地址"`
	ToAddress   string      `json:"to_address"       description:"接收地址"`
	Status      int         `json:"status"           description:"状态 0:待确认 1:已完成 2:失败"`
	CreatedAt   *gtime.Time `json:"created_at"       description:"创建时间"`
	CompletedAt *gtime.Time `json:"completed_at"     description:"完成时间"`
}

// AggregatorTransaction 聚合交易
type AggregatorTransaction struct {
	Id         uint64      `json:"id"               description:"ID"`
	Protocol   string      `json:"protocol"         description:"协议 1INCH/0X/PARASWAP"`
	ChainId    uint64      `json:"chain_id"         description:"链ID"`
	TxHash     string      `json:"tx_hash"          description:"交易哈希"`
	FromToken  string      `json:"from_token"       description:"支付代币"`
	ToToken    string      `json:"to_token"         description:"接收代币"`
	FromAmount string      `json:"from_amount"      description:"支付金额"`
	ToAmount   string      `json:"to_amount"        description:"接收金额"`
	Sender     string      `json:"sender"           description:"发送地址"`
	Receiver   string      `json:"receiver"         description:"接收地址"`
	Route      string      `json:"route"            description:"路由JSON"`
	Status     int         `json:"status"           description:"状态 0:待执行 1:已完成 2:失败"`
	CreatedAt  *gtime.Time `json:"created_at"       description:"创建时间"`
	ExecutedAt *gtime.Time `json:"executed_at"      description:"执行时间"`
}

// QuoteResult 报价结果
type QuoteResult struct {
	Protocol    string `json:"protocol"`
	FromToken   string `json:"from_token"`
	ToToken     string `json:"to_token"`
	FromAmount  string `json:"from_amount"`
	ToAmount    string `json:"to_amount"`
	Price       string `json:"price"`
	PriceImpact string `json:"price_impact"`
	Gas         uint64 `json:"gas"`
	Route       string `json:"route"`
}
