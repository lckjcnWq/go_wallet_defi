package controller

import (
	"context"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/service"
)

type ProtocolController struct{}

// SwapTokens DEX交易
func (c *ProtocolController) SwapTokens(ctx context.Context, req *v1.SwapTokensReq) (res *v1.SwapTokensRes, err error) {
	params := &defi.SwapParams{
		FromToken:  req.FromToken,
		ToToken:    req.ToToken,
		FromAmount: req.FromAmount,
		Sender:     req.Sender,
		Receiver:   req.Receiver,
	}

	txHash, err := service.Protocol.SwapTokens(ctx, req.Protocol, req.ChainId, params)
	if err != nil {
		return nil, err
	}

	return &v1.SwapTokensRes{
		TxHash: txHash,
	}, nil
}

// GetQuotes 获取报价
func (c *ProtocolController) GetQuotes(ctx context.Context, req *v1.GetQuotesReq) (res *v1.GetQuotesRes, err error) {
	params := &defi.SwapParams{
		FromToken:  req.FromToken,
		ToToken:    req.ToToken,
		FromAmount: req.FromAmount,
	}

	quotes, err := service.Protocol.GetQuotes(ctx, req.ChainId, params)
	if err != nil {
		return nil, err
	}

	return &v1.GetQuotesRes{
		Quotes: quotes,
	}, nil
}

// LendingSupply 存款
func (c *ProtocolController) LendingSupply(ctx context.Context, req *v1.LendingSupplyReq) (res *v1.LendingSupplyRes, err error) {
	params := &defi.LendingParams{
		Token:   req.Token,
		Amount:  req.Amount,
		Address: req.Address,
	}

	txHash, err := service.Protocol.LendingSupply(ctx, req.Protocol, req.ChainId, params)
	if err != nil {
		return nil, err
	}

	return &v1.LendingSupplyRes{
		TxHash: txHash,
	}, nil
}

// LendingBorrow 借款
func (c *ProtocolController) LendingBorrow(ctx context.Context, req *v1.LendingBorrowReq) (res *v1.LendingBorrowRes, err error) {
	params := &defi.LendingParams{
		Token:   req.Token,
		Amount:  req.Amount,
		Address: req.Address,
	}

	txHash, err := service.Protocol.LendingBorrow(ctx, req.Protocol, req.ChainId, params)
	if err != nil {
		return nil, err
	}

	return &v1.LendingBorrowRes{
		TxHash: txHash,
	}, nil
}

// NFTBuy 购买NFT
func (c *ProtocolController) NFTBuy(ctx context.Context, req *v1.NFTBuyReq) (res *v1.NFTBuyRes, err error) {
	params := &defi.NFTParams{
		ContractAddress: req.ContractAddress,
		TokenId:         req.TokenId,
		Price:           req.Price,
		PayToken:        req.PayToken,
		Buyer:           req.Buyer,
	}

	txHash, err := service.Protocol.NFTBuy(ctx, req.Protocol, req.ChainId, params)
	if err != nil {
		return nil, err
	}

	return &v1.NFTBuyRes{
		TxHash: txHash,
	}, nil
}

// NFTMint 铸造NFT
func (c *ProtocolController) NFTMint(ctx context.Context, req *v1.NFTMintReq) (res *v1.NFTMintRes, err error) {
	params := &defi.NFTMintParams{
		ContractAddress: req.ContractAddress,
		TokenURI:        req.TokenURI,
		Receiver:        req.Receiver,
	}

	txHash, err := service.Protocol.NFTMint(ctx, req.ChainId, params)
	if err != nil {
		return nil, err
	}

	return &v1.NFTMintRes{
		TxHash: txHash,
	}, nil
}

// BridgeAsset 跨链资产
func (c *ProtocolController) BridgeAsset(ctx context.Context, req *v1.BridgeAssetReq) (res *v1.BridgeAssetRes, err error) {
	params := &defi.BridgeParams{
		FromChainId: req.FromChainId,
		ToChainId:   req.ToChainId,
		Token:       req.Token,
		Amount:      req.Amount,
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
	}

	txHash, err := service.Protocol.BridgeAsset(ctx, req.Protocol, params)
	if err != nil {
		return nil, err
	}

	return &v1.BridgeAssetRes{
		TxHash: txHash,
	}, nil
}

// AggregateSwap 聚合交易
func (c *ProtocolController) AggregateSwap(ctx context.Context, req *v1.AggregateSwapReq) (res *v1.AggregateSwapRes, err error) {
	params := &defi.SwapParams{
		FromToken:  req.FromToken,
		ToToken:    req.ToToken,
		FromAmount: req.FromAmount,
		Sender:     req.Sender,
		Receiver:   req.Receiver,
	}

	txHash, err := service.Protocol.AggregateSwap(ctx, req.Protocol, req.ChainId, params)
	if err != nil {
		return nil, err
	}

	return &v1.AggregateSwapRes{
		TxHash: txHash,
	}, nil
}
