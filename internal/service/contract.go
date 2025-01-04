package service

import (
	"context"
	"go-wallet-defi/internal/logic"
	"go-wallet-defi/internal/model"
)

type IContract interface {
	// Deploy 部署合约
	Deploy(ctx context.Context, name, abi, bytecode, from, args, network string) (hash, address string, err error)

	// Call 调用合约
	Call(ctx context.Context, address, method, args, from, value string) (result interface{}, err error)

	// GetEvents 获取合约事件
	GetEvents(ctx context.Context, address, eventName string, fromBlock, toBlock int64, page, pageSize int) ([]*model.ContractEvent, int, error)
}

// Contract 获取合约服务
func Contract() IContract {
	if localContract == nil {
		localContract = &logic.ContractLogic{}
	}
	return localContract
}

// SetContract 设置合约服务
func SetContract(s IContract) {
	localContract = s
}

var localContract IContract
