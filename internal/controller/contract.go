package controller

import (
	"context"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/service"
)

type ContractController struct{}

// Deploy 部署合约
func (c *ContractController) Deploy(ctx context.Context, req *v1.DeployContractReq) (res *v1.DeployContractRes, err error) {
	hash, address, err := service.Contract().Deploy(ctx, req.Name, req.ABI, req.Bytecode, req.From, req.Args, req.Network)
	if err != nil {
		return nil, err
	}

	return &v1.DeployContractRes{
		Hash:    hash,
		Address: address,
	}, nil
}

// Call 调用合约
func (c *ContractController) Call(ctx context.Context, req *v1.CallContractReq) (res *v1.CallContractRes, err error) {
	result, err := service.Contract().Call(ctx, req.Address, req.Method, req.Args, req.From, req.Value)
	if err != nil {
		return nil, err
	}

	return &v1.CallContractRes{
		Hash:   result.(string),
		Result: result,
	}, nil
}

// GetEvents 获取合约事件
func (c *ContractController) GetEvents(ctx context.Context, req *v1.GetContractEventsReq) (res *v1.GetContractEventsRes, err error) {
	events, total, err := service.Contract().GetEvents(ctx, req.Address, req.EventName, req.FromBlock, req.ToBlock, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]v1.ContractEventInfo, 0, len(events))
	for _, event := range events {
		list = append(list, v1.ContractEventInfo{
			Name:        event.Name,
			Data:        event.Data,
			BlockNumber: event.BlockNumber,
			TxHash:      event.TxHash,
			CreatedAt:   event.CreatedAt,
		})
	}

	return &v1.GetContractEventsRes{
		List:  list,
		Total: total,
	}, nil
}
