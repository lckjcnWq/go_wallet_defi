package service

import (
	"context"
	"encoding/json"
	"errors"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"math/big"
	"sort"
)

type ProtocolService struct{}

var Protocol = &ProtocolService{}

func (s *ProtocolService) SwapTokens(ctx context.Context, protocol string, chainId uint64, params *defi.SwapParams) (string, error) {
	dex, err := defi.GetDEX(Protocol)

	if err != nil {
		return "", err
	}
	// 查询兑换路径
	route, err := dex.GetBestRoute(ctx, chainId, params)
	if err != nil {
		return "", err
	}

	// 执行兑换
	txHash, err := dex.Swap(ctx, chainId, params, route)
	if err != nil {
		return "", err
	}

	// 保存交易记录
	tx := &model.SwapTransaction{
		ChainId:    chainId,
		Protocol:   protocol,
		TxHash:     txHash,
		FromToken:  params.FromToken,
		ToToken:    params.ToToken,
		FromAmount: params.FromAmount,
		ToAmount:   route.ToAmount,
		Sender:     params.Sender,
		Receiver:   params.Receiver,
		Status:     0,
	}

	err = dao.Protocol.CreateSwapTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// GetQuotes 获取兑换报价
func (s *ProtocolService) GetQuotes(ctx context.Context, chainId uint64, params *defi.SwapParams) ([]*model.QuoteResult, error) {
	var quotes []*model.QuoteResult

	// 获取所有支持的DEX
	dexes := defi.GetSupportedDEXes()

	// 并行获取报价
	ch := make(chan *model.QuoteResult, len(dexes))
	for _, protocol := range dexes {
		go func(protocol string) {
			dex, err := defi.GetDEX(protocol)
			if err != nil {
				ch <- nil
				return
			}

			route, err := dex.GetBestRoute(ctx, chainId, params)
			if err != nil {
				ch <- nil
				return
			}

			quote := &model.QuoteResult{
				Protocol:    protocol,
				FromToken:   params.FromToken,
				ToToken:     params.ToToken,
				FromAmount:  params.FromAmount,
				ToAmount:    route.ToAmount,
				Price:       route.Price,
				PriceImpact: route.PriceImpact,
				Gas:         route.Gas,
				Route:       route.Path,
			}
			ch <- quote
		}(protocol)
	}

	// 收集报价结果
	for range dexes {
		if quote := <-ch; quote != nil {
			quotes = append(quotes, quote)
		}
	}

	// 按照收益排序
	sort.Slice(quotes, func(i, j int) bool {
		amountI, _ := new(big.Int).SetString(quotes[i].ToAmount, 10)
		amountJ, _ := new(big.Int).SetString(quotes[j].ToAmount, 10)
		return amountI.Cmp(amountJ) > 0
	})

	return quotes, nil
}

// LendingSupply 存款
func (s *ProtocolService) LendingSupply(ctx context.Context, protocol string, chainId uint64, params *defi.LendingParams) (string, error) {
	// 获取借贷协议实例
	lending, err := defi.GetLendingProtocol(protocol)
	if err != nil {
		return "", err
	}

	// 执行存款
	txHash, err := lending.Supply(ctx, chainId, params)
	if err != nil {
		return "", err
	}

	// 更新借贷仓位
	position, err := lending.GetPosition(ctx, chainId, params.Address, params.Token)
	if err != nil {
		return "", err
	}

	err = dao.Protocol.UpdateLendingPosition(ctx, position)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// LendingBorrow 借款
func (s *ProtocolService) LendingBorrow(ctx context.Context, protocol string, chainId uint64, params *defi.LendingParams) (string, error) {
	lending, err := defi.GetLendingProtocol(protocol)
	if err != nil {
		return "", err
	}

	// 检查健康因子
	healthFactor, err := lending.GetHealthFactor(ctx, chainId, params.Address)
	if err != nil {
		return "", err
	}

	if healthFactor.Cmp(big.NewInt(120)) < 0 { // 健康因子小于1.2
		return "", errors.New("health factor too low")
	}

	// 执行借款
	txHash, err := lending.Borrow(ctx, chainId, params)
	if err != nil {
		return "", err
	}

	// 更新仓位
	position, err := lending.GetPosition(ctx, chainId, params.Address, params.Token)
	if err != nil {
		return "", err
	}

	err = dao.Protocol.UpdateLendingPosition(ctx, position)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// NFTBuy 购买NFT
func (s *ProtocolService) NFTBuy(ctx context.Context, protocol string, chainId uint64, params *defi.NFTParams) (string, error) {
	// 获取NFT交易所实例
	marketplace, err := defi.GetNFTMarketplace(protocol)
	if err != nil {
		return "", err
	}

	// 检查价格和可用性
	listing, err := marketplace.GetListing(ctx, chainId, params.ContractAddress, params.TokenId)
	if err != nil {
		return "", err
	}

	if listing.Price.Cmp(new(big.Int).SetString(params.Price, 10)) != 0 {
		return "", errors.New("price mismatch")
	}

	// 执行购买
	txHash, err := marketplace.Buy(ctx, chainId, params)
	if err != nil {
		return "", err
	}

	// 保存交易记录
	tx := &model.NFTTransaction{
		ChainId:         chainId,
		Protocol:        protocol,
		TxHash:          txHash,
		ContractAddress: params.ContractAddress,
		TokenId:         params.TokenId,
		FromAddress:     listing.Seller,
		ToAddress:       params.Buyer,
		Price:           params.Price,
		PayToken:        params.PayToken,
		Type:            "SALE",
		Status:          0,
	}

	err = dao.Protocol.CreateNFTTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// NFTMint 铸造NFT
func (s *ProtocolService) NFTMint(ctx context.Context, chainId uint64, params *defi.NFTMintParams) (string, error) {
	client, err := ethclient.GetClient(chainId)
	if err != nil {
		return "", err
	}

	// 铸造NFT
	txHash, err := client.MintNFT(ctx, params.ContractAddress, params.TokenURI, params.Receiver)
	if err != nil {
		return "", err
	}

	// 保存交易记录
	tx := &model.NFTTransaction{
		ChainId:         chainId,
		Protocol:        "MINT",
		TxHash:          txHash,
		ContractAddress: params.ContractAddress,
		TokenId:         "0", // Will be updated after confirmation
		FromAddress:     "0x0000000000000000000000000000000000000000",
		ToAddress:       params.Receiver,
		Price:           "0",
		Type:            "MINT",
		Status:          0,
	}

	err = dao.Protocol.CreateNFTTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// BridgeAsset 跨链资产转移
func (s *ProtocolService) BridgeAsset(ctx context.Context, protocol string, params *defi.BridgeParams) (string, error) {
	// 获取跨链桥实例
	bridge, err := defi.GetBridge(protocol)
	if err != nil {
		return "", err
	}

	// 获取跨链费用
	fee, err := bridge.EstimateFee(ctx, params)
	if err != nil {
		return "", err
	}

	// 执行跨链
	txHash, err := bridge.Bridge(ctx, params)
	if err != nil {
		return "", err
	}

	// 保存交易记录
	tx := &model.BridgeTransaction{
		Protocol:    protocol,
		FromChainId: params.FromChainId,
		ToChainId:   params.ToChainId,
		FromTxHash:  txHash,
		Token:       params.Token,
		Amount:      params.Amount,
		FromAddress: params.FromAddress,
		ToAddress:   params.ToAddress,
		Status:      0,
	}

	err = dao.Protocol.CreateBridgeTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

// AggregateSwap 聚合交易
func (s *ProtocolService) AggregateSwap(ctx context.Context, protocol string, chainId uint64, params *defi.SwapParams) (string, error) {
	// 获取聚合器实例
	aggregator, err := defi.GetAggregator(protocol)
	if err != nil {
		return "", err
	}

	// 获取最优路由
	route, err := aggregator.GetBestRoute(ctx, chainId, params)
	if err != nil {
		return "", err
	}

	// 执行交易
	txHash, err := aggregator.Swap(ctx, chainId, params, route)
	if err != nil {
		return "", err
	}

	// 保存交易记录
	routeJson, _ := json.Marshal(route)
	tx := &model.AggregatorTransaction{
		Protocol:   protocol,
		ChainId:    chainId,
		TxHash:     txHash,
		FromToken:  params.FromToken,
		ToToken:    params.ToToken,
		FromAmount: params.FromAmount,
		ToAmount:   route.ToAmount,
		Sender:     params.Sender,
		Receiver:   params.Receiver,
		Route:      string(routeJson),
		Status:     0,
	}

	err = dao.Protocol.CreateAggregatorTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
