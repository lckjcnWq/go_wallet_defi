package task

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/ethclientx"
	"math/big"
	"strings"
	"time"
)

// WatchContractEvents 监听合约事件
func WatchContractEvents() {
	ctx := context.Background()
	client := ethclientx.GetClient(ctx)

	for {
		// 获取所有合约
		var contracts []*model.Contract
		err := g.DB().Model("contract").
			Where("status", 1).
			Scan(&contracts)

		if err != nil {
			g.Log().Error(ctx, err)
			time.Sleep(time.Second * 10)
			continue
		}

		for _, contract := range contracts {
			// 获取最新事件高度
			var lastEvent *model.ContractEvent
			err := g.DB().Model("contract_event").
				Where("contract_id", contract.Id).
				Order("block_number DESC").
				Scan(&lastEvent)

			fromBlock := contract.DeployHeight
			if lastEvent != nil {
				fromBlock = lastEvent.BlockNumber + 1
			}

			// 获取当前区块高度
			header, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}

			// 解析ABI
			parsed, err := abi.JSON(strings.NewReader(contract.ABI))
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}

			// 构建事件过滤器
			query := ethereum.FilterQuery{
				FromBlock: big.NewInt(fromBlock),
				ToBlock:   header.Number,
				Addresses: []common.Address{common.HexToAddress(contract.Address)},
			}

			// 获取日志
			logs, err := client.FilterLogs(ctx, query)
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}

			// 解析事件日志
			for _, log := range logs {
				// 查找事件
				var event *abi.Event
				for _, e := range parsed.Events {
					if e.ID() == log.Topics[0] {
						event = &e
						break
					}
				}
				if event == nil {
					continue
				}

				// 解析事件数据
				data, err := event.Inputs.UnpackValues(log.Data)
				if err != nil {
					g.Log().Error(ctx, err)
					continue
				}

				// 转换为JSON
				dataJson, err := json.Marshal(data)
				if err != nil {
					g.Log().Error(ctx, err)
					continue
				}

				// 保存事件
				contractEvent := &model.ContractEvent{
					ContractId:  contract.Id,
					Name:        event.Name,
					Signature:   event.Sig(),
					Topics:      common.Bytes2Hex(log.Topics[0].Bytes()),
					Data:        string(dataJson),
					BlockNumber: int64(log.BlockNumber),
					BlockHash:   log.BlockHash.Hex(),
					TxHash:      log.TxHash.Hex(),
					TxIndex:     int(log.TxIndex),
					LogIndex:    int(log.Index),
					CreatedAt:   time.Now().Unix(),
				}

				err = dao.Contract.InsertEvent(ctx, contractEvent)
				if err != nil {
					g.Log().Error(ctx, err)
				}
			}
		}

		time.Sleep(time.Second * 10)
	}
}
