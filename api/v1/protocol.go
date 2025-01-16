package v1

import (
	"go-wallet-defi/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type SwapTokensReq struct {
	g.Meta     `path:"/protocol/swap" method:"post"`
	Protocol   string `json:"protocol"    v:"required"`
	ChainId    uint64 `json:"chain_id"    v:"required"`
	FromToken  string `json:"from_token"  v:"required"`
	ToToken    string `json:"to_token"    v:"required"`
	FromAmount string `json:"from_amount" v:"required"`
	Sender     string `json:"sender"      v:"required"`
	Receiver   string `json:"receiver"    v:"required"`
}

type SwapTokensRes struct {
	TxHash string `json:"tx_hash"`
}

type GetQuotesReq struct {
	g.Meta     `path:"/protocol/quotes" method:"get"`
	ChainId    uint64 `json:"chain_id"    v:"required"`
	FromToken  string `json:"from_token"  v:"required"`
	ToToken    string `json:"to_token"    v:"required"`
	FromAmount string `json:"from_amount" v:"required"`
}

type GetQuotesRes struct {
	Quotes []*model.QuoteResult `json:"quotes"`
}

type LendingSupplyReq struct {
	g.Meta   `path:"/protocol/lending/supply" method:"post"`
	Protocol string `json:"protocol" v:"required"`
	ChainId  uint64 `json:"chain_id" v:"required"`
	Token    string `json:"token"    v:"required"`
	Amount   string `json:"amount"   v:"required"`
	Address  string `json:"address"  v:"required"`
}

type LendingSupplyRes struct {
	TxHash string `json:"tx_hash"`
}

type LendingBorrowReq struct {
	g.Meta   `path:"/protocol/lending/borrow" method:"post"`
	Protocol string `json:"protocol" v:"required"`
	ChainId  uint64 `json:"chain_id" v:"required"`
	Token    string `json:"token"    v:"required"`
	Amount   string `json:"amount"   v:"required"`
	Address  string `json:"address"  v:"required"`
}

type LendingBorrowRes struct {
	TxHash string `json:"tx_hash"`
}

type NFTBuyReq struct {
	g.Meta          `path:"/protocol/nft/buy" method:"post"`
	Protocol        string `json:"protocol"         v:"required"`
	ChainId         uint64 `json:"chain_id"         v:"required"`
	ContractAddress string `json:"contract_address" v:"required"`
	TokenId         string `json:"token_id"         v:"required"`
	Price           string `json:"price"            v:"required"`
	PayToken        string `json:"pay_token"        v:"required"`
	Buyer           string `json:"buyer"            v:"required"`
}

type NFTBuyRes struct {
	TxHash string `json:"tx_hash"`
}

type NFTMintReq struct {
	g.Meta          `path:"/protocol/nft/mint" method:"post"`
	ChainId         uint64 `json:"chain_id"         v:"required"`
	ContractAddress string `json:"contract_address" v:"required"`
	TokenURI        string `json:"token_uri"        v:"required"`
	Receiver        string `json:"receiver"         v:"required"`
}

type NFTMintRes struct {
	TxHash string `json:"tx_hash"`
}

type BridgeAssetReq struct {
	g.Meta      `path:"/protocol/bridge" method:"post"`
	Protocol    string `json:"protocol"     v:"required"`
	FromChainId uint64 `json:"from_chain_id" v:"required"`
	ToChainId   uint64 `json:"to_chain_id"   v:"required"`
	Token       string `json:"token"        v:"required"`
	Amount      string `json:"amount"       v:"required"`
	FromAddress string `json:"from_address"  v:"required"`
	ToAddress   string `json:"to_address"    v:"required"`
}

type BridgeAssetRes struct {
	TxHash string `json:"tx_hash"`
}

type AggregateSwapReq struct {
	g.Meta     `path:"/protocol/aggregate/swap" method:"post"`
	Protocol   string `json:"protocol"    v:"required"`
	ChainId    uint64 `json:"chain_id"    v:"required"`
	FromToken  string `json:"from_token"  v:"required"`
	ToToken    string `json:"to_token"    v:"required"`
	FromAmount string `json:"from_amount" v:"required"`
	Sender     string `json:"sender"      v:"required"`
	Receiver   string `json:"receiver"    v:"required"`
}

type AggregateSwapRes struct {
	TxHash string `json:"tx_hash"`
}
