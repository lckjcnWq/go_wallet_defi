package model

// NFT NFT信息
type NFT struct {
	Id          uint64 `json:"id"`          // ID
	ContractId  uint64 `json:"contractId"`  // 合约ID
	TokenId     string `json:"tokenId"`     // 代币ID
	Owner       string `json:"owner"`       // 持有者地址
	Creator     string `json:"creator"`     // 创建者地址
	URI         string `json:"uri"`         // 元数据URI
	Name        string `json:"name"`        // NFT名称
	Description string `json:"description"` // NFT描述
	Image       string `json:"image"`       // 图片URI
	ExternalUrl string `json:"externalUrl"` // 外部链接
	Attributes  string `json:"attributes"`  // 属性JSON
	Standard    string `json:"standard"`    // 标准(ERC721/ERC1155)
	Amount      uint64 `json:"amount"`      // 数量(ERC1155)
	Status      int    `json:"status"`      // 状态
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
	UpdatedAt   int64  `json:"updatedAt"`   // 更新时间
}

// NFTTransfer NFT交易记录
type NFTTransfer struct {
	Id          uint64 `json:"id"`          // ID
	NftId       uint64 `json:"nftId"`       // NFT ID
	From        string `json:"from"`        // 发送地址
	To          string `json:"to"`          // 接收地址
	Amount      uint64 `json:"amount"`      // 数量
	Type        string `json:"type"`        // 类型
	Hash        string `json:"hash"`        // 交易哈希
	BlockNumber int64  `json:"blockNumber"` // 区块高度
	BlockTime   int64  `json:"blockTime"`   // 区块时间
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
}

// NFTMarket NFT市场信息
type NFTMarket struct {
	Id        uint64 `json:"id"`        // ID
	NftId     uint64 `json:"nftId"`     // NFT ID
	Seller    string `json:"seller"`    // 卖家地址
	Price     string `json:"price"`     // 价格
	Currency  string `json:"currency"`  // 支付代币
	Status    int    `json:"status"`    // 状态
	ExpiresAt int64  `json:"expiresAt"` // 过期时间
	CreatedAt int64  `json:"createdAt"` // 创建时间
	UpdatedAt int64  `json:"updatedAt"` // 更新时间
}

// NFTApproval NFT授权信息
type NFTApproval struct {
	Id         uint64 `json:"id"`         // ID
	ContractId uint64 `json:"contractId"` // 合约ID
	Owner      string `json:"owner"`      // 持有者地址
	Operator   string `json:"operator"`   // 授权地址
	Approved   bool   `json:"approved"`   // 是否授权
	CreatedAt  int64  `json:"createdAt"`  // 创建时间
	UpdatedAt  int64  `json:"updatedAt"`  // 更新时间
}
