package v1

import "github.com/gogf/gf/v2/frame/g"

// MintNFTReq 铸造NFT
type MintNFTReq struct {
	g.Meta      `path:"/nft/mint" method:"post" summary:"铸造NFT" tags:"NFT管理"`
	ContractId  uint64 `json:"contractId" dc:"合约ID"`
	To          string `json:"to" dc:"接收地址"`
	TokenId     string `dc:"代币ID(可选)"`
	Amount      uint64 `d:"1" dc:"数量(ERC1155)"`
	Name        string `v:"required" dc:"NFT名称"`
	Description string `dc:"NFT描述"`
	Image       string `v:"required" dc:"图片"`
	Attributes  string `dc:"属性JSON"`
}

type MintNFTRes struct {
	Hash    string `json:"hash" dc:"交易哈希"`
	TokenId string `json:"tokenId" dc:"代币ID"`
}

// TransferNFTReq 转移NFT
type TransferNFTReq struct {
	g.Meta `path:"/nft/transfer" method:"post" tags:"NFT管理" summary:"转移NFT"`
	NftId  uint64 `v:"required" dc:"NFT ID"`
	From   string `v:"required" dc:"发送地址"`
	To     string `v:"required" dc:"接收地址"`
	Amount uint64 `d:"1" dc:"数量(ERC1155)"`
}

type TransferNFTRes struct {
	Hash string `json:"hash" dc:"交易哈希"`
}

// ListNFTReq 获取NFT列表
type ListNFTReq struct {
	g.Meta     `path:"/nft/list" method:"get" tags:"NFT管理" summary:"NFT列表"`
	Owner      string `dc:"持有者地址"`
	Creator    string `dc:"创建者地址"`
	ContractId uint64 `dc:"合约ID"`
	Page       int    `d:"1" dc:"页码"`
	PageSize   int    `d:"10" dc:"每页数量"`
}

type ListNFTRes struct {
	List  []NFTInfo `json:"list" dc:"NFT列表"`
	Total int       `json:"total" dc:"总数"`
}

type NFTInfo struct {
	Id          uint64                   `json:"id" dc:"NFT ID"`
	ContractId  uint64                   `json:"contractId" dc:"合约ID"`
	TokenId     string                   `json:"tokenId" dc:"代币ID"`
	Owner       string                   `json:"owner" dc:"持有者地址"`
	Creator     string                   `json:"creator" dc:"创建者地址"`
	Name        string                   `json:"name" dc:"NFT名称"`
	Description string                   `json:"description" dc:"NFT描述"`
	Image       string                   `json:"image" dc:"图片"`
	Attributes  []map[string]interface{} `json:"attributes" dc:"属性"`
	Standard    string                   `json:"standard" dc:"标准"`
	Amount      uint64                   `json:"amount" dc:"数量"`
	CreatedAt   int64                    `json:"createdAt" dc:"创建时间"`
}

// ListNFTMarketReq 获取NFT市场列表
type ListNFTMarketReq struct {
	g.Meta   `path:"/nft/market/list" method:"get" tags:"NFT管理" summary:"市场列表"`
	Seller   string `dc:"卖家地址"`
	Status   int    `d:"1" dc:"状态 1:在售"`
	Page     int    `d:"1" dc:"页码"`
	PageSize int    `d:"10" dc:"每页数量"`
}

type ListNFTMarketRes struct {
	List  []NFTMarketInfo `json:"list" dc:"市场列表"`
	Total int             `json:"total" dc:"总数"`
}

type NFTMarketInfo struct {
	Id        uint64  `json:"id" dc:"ID"`
	NFT       NFTInfo `json:"nft" dc:"NFT信息"`
	Seller    string  `json:"seller" dc:"卖家地址"`
	Price     string  `json:"price" dc:"价格"`
	Currency  string  `json:"currency" dc:"支付代币"`
	Status    int     `json:"status" dc:"状态"`
	ExpiresAt int64   `json:"expiresAt" dc:"过期时间"`
	CreatedAt int64   `json:"createdAt" dc:"创建时间"`
}
