package service

import (
	"context"
	"go-wallet-defi/api/v1"
	"go-wallet-defi/internal/logic"
	"go-wallet-defi/internal/model"
)

type INFT interface {
	// Mint 铸造NFT
	Mint(ctx context.Context, req *v1.MintNFTReq) (hash string, tokenId string, err error)

	// Transfer 转移NFT
	Transfer(ctx context.Context, nftId uint64, from, to string, amount uint64) (hash string, err error)

	// List 获取NFT列表
	List(ctx context.Context, owner, creator string, contractId uint64, page, pageSize int) ([]*model.NFT, int, error)

	// ListMarket 获取市场列表
	ListMarket(ctx context.Context, seller string, status int, page, pageSize int) ([]*model.NFTMarket, int, error)
}

// NFT 获取NFT服务
func NFT() INFT {
	if localNFT == nil {
		localNFT = &logic.NFTLogic{}
	}
	return localNFT
}

var localNFT INFT
