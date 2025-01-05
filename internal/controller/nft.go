package controller

import (
	"context"
	"encoding/json"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/service"
)

type NFTController struct{}

// Mint 铸造NFT
func (c *NFTController) Mint(ctx context.Context, req *v1.MintNFTReq) (res *v1.MintNFTRes, err error) {
	hash, tokenId, err := service.NFT().Mint(ctx, req)
	if err != nil {
		return nil, err
	}

	return &v1.MintNFTRes{
		Hash:    hash,
		TokenId: tokenId,
	}, nil
}

// Transfer 转移NFT
func (c *NFTController) Transfer(ctx context.Context, req *v1.TransferNFTReq) (res *v1.TransferNFTRes, err error) {
	hash, err := service.NFT().Transfer(ctx, req.NftId, req.From, req.To, req.Amount)
	if err != nil {
		return nil, err
	}

	return &v1.TransferNFTRes{
		Hash: hash,
	}, nil
}

// List 获取NFT列表
func (c *NFTController) List(ctx context.Context, req *v1.ListNFTReq) (res *v1.ListNFTRes, err error) {
	nfts, total, err := service.NFT().List(ctx, req.Owner, req.Creator, req.ContractId, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]v1.NFTInfo, 0, len(nfts))
	for _, nft := range nfts {
		var attrs []map[string]interface{}
		if nft.Attributes != "" {
			_ = json.Unmarshal([]byte(nft.Attributes), &attrs)
		}

		list = append(list, v1.NFTInfo{
			Id:          nft.Id,
			ContractId:  nft.ContractId,
			TokenId:     nft.TokenId,
			Owner:       nft.Owner,
			Creator:     nft.Creator,
			Name:        nft.Name,
			Description: nft.Description,
			Image:       nft.Image,
			Attributes:  attrs,
			Standard:    nft.Standard,
			Amount:      nft.Amount,
			CreatedAt:   nft.CreatedAt,
		})
	}

	return &v1.ListNFTRes{
		List:  list,
		Total: total,
	}, nil
}

// ListMarket 获取市场列表
func (c *NFTController) ListMarket(ctx context.Context, req *v1.ListNFTMarketReq) (res *v1.ListNFTMarketRes, err error) {
	markets, total, err := service.NFT().ListMarket(ctx, req.Seller, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]v1.NFTMarketInfo, 0, len(markets))
	for _, market := range markets {
		nft, _ := dao.NFT.GetById(ctx, market.NftId)
		if nft == nil {
			continue
		}

		var attrs []map[string]interface{}
		if nft.Attributes != "" {
			_ = json.Unmarshal([]byte(nft.Attributes), &attrs)
		}

		list = append(list, v1.NFTMarketInfo{
			Id: market.Id,
			NFT: v1.NFTInfo{
				Id:          nft.Id,
				ContractId:  nft.ContractId,
				TokenId:     nft.TokenId,
				Owner:       nft.Owner,
				Creator:     nft.Creator,
				Name:        nft.Name,
				Description: nft.Description,
				Image:       nft.Image,
				Attributes:  attrs,
				Standard:    nft.Standard,
				Amount:      nft.Amount,
				CreatedAt:   nft.CreatedAt,
			},
			Seller:    market.Seller,
			Price:     market.Price,
			Currency:  market.Currency,
			Status:    market.Status,
			ExpiresAt: market.ExpiresAt,
			CreatedAt: market.CreatedAt,
		})
	}

	return &v1.ListNFTMarketRes{
		List:  list,
		Total: total,
	}, nil
}
