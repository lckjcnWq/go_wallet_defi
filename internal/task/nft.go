package task

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/contracts/nft"
	"go-wallet-defi/internal/pkg/ethclient"
	"math/big"
	"strings"
	"time"
)

// WatchNFTEvents 监听NFT事件
func WatchNFTEvents() {
	ctx := context.Background()
	client := ethclient.GetClient(ctx)

	// 解析ERC721 ABI
	erc721ABI, _ := abi.JSON(strings.NewReader(nft.ERC721ABI))

	// 解析ERC1155 ABI
	erc1155ABI, _ := abi.JSON(strings.NewReader(nft.ERC1155ABI))

	for {
		// 获取所有NFT合约
		var contracts []*model.Contract
		err := g.DB().Model("contract").
			Where("status", 1).
			WhereIn("standard", []string{"ERC721", "ERC1155"}).
			Scan(&contracts)

		if err != nil {
			g.Log().Error(ctx, err)
			time.Sleep(time.Second * 10)
			continue
		}

		for _, contract := range contracts {
			// 获取最新事件高度
			var lastTransfer *model.NFTTransfer
			err := g.DB().Model("nft_transfer").
				Where("contract_id", contract.Id).
				Order("block_number DESC").
				Scan(&lastTransfer)

			fromBlock := contract.DeployHeight
			if lastTransfer != nil {
				fromBlock = lastTransfer.BlockNumber + 1
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
				var (
					from    common.Address
					to      common.Address
					tokenId *big.Int
					amount  *big.Int
				)

				// 根据合约标准解析事件
				if contract.Standard == "ERC721" {
					if log.Topics[0] == erc721ABI.Events["Transfer"].ID {
						from = common.HexToAddress(log.Topics[1].Hex())
						to = common.HexToAddress(log.Topics[2].Hex())
						tokenId = new(big.Int).SetBytes(log.Topics[3].Bytes())
						amount = big.NewInt(1)
					}
				} else {
					if log.Topics[0] == erc1155ABI.Events["TransferSingle"].ID {
						from = common.HexToAddress(log.Topics[1].Hex())
						to = common.HexToAddress(log.Topics[2].Hex())
						tokenId = new(big.Int).SetBytes(log.Topics[3].Bytes())
						amount = new(big.Int).SetBytes(log.Data)
					}
				}

				if tokenId == nil {
					continue
				}

				// 查找NFT
				nft, err := dao.NFT.GetByTokenId(ctx, contract.Id, tokenId.String())
				if err != nil {
					g.Log().Error(ctx, err)
					continue
				}

				// 保存转移记录
				transfer := &model.NFTTransfer{
					NftId:       nft.Id,
					From:        from.Hex(),
					To:          to.Hex(),
					Amount:      amount.Uint64(),
					Type:        "transfer",
					Hash:        log.TxHash.Hex(),
					BlockNumber: int64(log.BlockNumber),
					BlockTime:   time.Now().Unix(),
					CreatedAt:   time.Now().Unix(),
				}

				err = dao.NFT.InsertTransfer(ctx, transfer)
				if err != nil {
					g.Log().Error(ctx, err)
					continue
				}

				// 更新NFT所有者
				err = dao.NFT.UpdateOwner(ctx, nft.Id, to.Hex())
				if err != nil {
					g.Log().Error(ctx, err)
				}
			}
		}

		time.Sleep(time.Second * 10)
	}
}
现在第4关的NFT模块已经完整
