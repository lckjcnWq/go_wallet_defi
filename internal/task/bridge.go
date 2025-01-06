package task

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/contracts/bridge"
	"go-wallet-defi/internal/service"
	"math/big"
	"strings"
	"time"
)

// WatchBridgeEvents 监听跨链桥事件
func WatchBridgeEvents() {
	ctx := context.Background()

	// 解析ABI
	parsed, _ := abi.JSON(strings.NewReader(bridge.BridgeABI))

	for {
		// 获取所有链信息
		var chains []*model.Chain
		err := g.DB().Model("chain").
			Where("status", 1).
			WhereNot("bridge_address", "").
			Scan(&chains)

		if err != nil {
			g.Log().Error(ctx, err)
			time.Sleep(time.Second * 10)
			continue
		}

		for _, chain := range chains {
			// 获取最新事件高度
			var lastTransfer *model.CrossTransfer
			err := g.DB().Model("cross_transfer").
				Where("from_chain_id = ? OR to_chain_id = ?", chain.ChainId, chain.ChainId).
				Order("block_number DESC").
				Scan(&lastTransfer)

			fromBlock := chain.DeployHeight
			if lastTransfer != nil && lastTransfer.BlockNumber > fromBlock {
				fromBlock = lastTransfer.BlockNumber + 1
			}

			// 获取客户端
			client, err := ethclient.GetClient(chain.ChainId)
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}

			// 获取当前区块高度
			header, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}

			// 构建事件过滤器
			query := ethereum.FilterQuery{
				FromBlock: big.NewInt(fromBlock),
				ToBlock:   header.Number,
				Addresses: []common.Address{common.HexToAddress(chain.BridgeAddress)},
			}

			// 获取日志
			logs, err := client.FilterLogs(ctx, query)
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}

			// 处理事件日志
			for _, log := range logs {
				switch log.Topics[0] {
				case parsed.Events["Lock"].ID:
					// 解析Lock事件
					var lockEvent struct {
						Token     common.Address
						From      common.Address
						Amount    *big.Int
						ToChainId *big.Int
						ToAddress common.Address
						Nonce     *big.Int
					}

					err = parsed.UnpackIntoInterface(&lockEvent, "Lock", log.Data)
					if err != nil {
						g.Log().Error(ctx, err)
						continue
					}

					// 处理锁定事件
					err = service.Bridge().ProcessLockEvent(ctx,
						chain.ChainId,
						lockEvent.Token.Hex(),
						lockEvent.From.Hex(),
						lockEvent.Amount.String(),
						lockEvent.ToChainId.Uint64(),
						lockEvent.ToAddress.Hex(),
						lockEvent.Nonce.Uint64(),
						log.TxHash.Hex(),
					)
					if err != nil {
						g.Log().Error(ctx, err)
					}

				case parsed.Events["Unlock"].ID:
					// 解析Unlock事件
					var unlockEvent struct {
						Token       common.Address
						To          common.Address
						Amount      *big.Int
						FromChainId *big.Int
						Nonce       *big.Int
					}

					err = parsed.UnpackIntoInterface(&unlockEvent, "Unlock", log.Data)
					if err != nil {
						g.Log().Error(ctx, err)
						continue
					}

					// 处理解锁事件
					err = service.Bridge().ProcessUnlockEvent(ctx,
						chain.ChainId,
						unlockEvent.Token.Hex(),
						unlockEvent.To.Hex(),
						unlockEvent.Amount.String(),
						unlockEvent.FromChainId.Uint64(),
						unlockEvent.Nonce.Uint64(),
						log.TxHash.Hex(),
					)
					if err != nil {
						g.Log().Error(ctx, err)
					}
				}
			}
		}

		time.Sleep(time.Second * 10)
	}
}
