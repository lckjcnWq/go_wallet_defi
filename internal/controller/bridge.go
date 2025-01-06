package controller

import (
	"context"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/service"
)

type BridgeController struct{}

// CrossTransfer 跨链转账
func (c *BridgeController) CrossTransfer(ctx context.Context, req *v1.CrossTransferReq) (res *v1.CrossTransferRes, err error) {
	hash, nonce, err := service.Bridge().CrossTransfer(ctx,
		req.FromChainId,
		req.ToChainId,
		req.FromAddress,
		req.ToAddress,
		req.TokenAddress,
		req.Amount,
	)
	if err != nil {
		return nil, err
	}

	return &v1.CrossTransferRes{
		Hash:  hash,
		Nonce: nonce,
	}, nil
}

// GetCrossTransfer 获取跨链交易列表
func (c *BridgeController) GetCrossTransfer(ctx context.Context, req *v1.GetCrossTransferReq) (res *v1.GetCrossTransferRes, err error) {
	transfers, total, err := service.Bridge().GetCrossTransfers(ctx,
		req.FromChainId,
		req.ToChainId,
		req.Address,
		req.Status,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		return nil, err
	}

	list := make([]v1.CrossTransferInfo, 0, len(transfers))
	for _, transfer := range transfers {
		list = append(list, v1.CrossTransferInfo{
			Id:           transfer.Id,
			FromChainId:  transfer.FromChainId,
			ToChainId:    transfer.ToChainId,
			FromAddress:  transfer.FromAddress,
			ToAddress:    transfer.ToAddress,
			TokenAddress: transfer.TokenAddress,
			Amount:       transfer.Amount,
			Fee:          transfer.Fee,
			FromHash:     transfer.FromHash,
			ToHash:       transfer.ToHash,
			Status:       transfer.Status,
			CreatedAt:    transfer.CreatedAt,
		})
	}

	return &v1.GetCrossTransferRes{
		List:  list,
		Total: total,
	}, nil
}
