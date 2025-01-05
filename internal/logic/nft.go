package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/ethclient"
	"go-wallet-defi/internal/pkg/ipfs"
	"math/big"
	"strings"
	"time"
)

type NFTLogic struct{}

// Mint 铸造NFT
func (l *NFTLogic) Mint(ctx context.Context, req *v1.MintNFTReq) (hash string, tokenId string, err error) {
	//1.获取合约信息
	contract, err := dao.Contract.GetById(ctx, req.ContractId)
	if err != nil {
		return "", "", err
	}
	//2.判断合约标准
	var standard string
	if strings.Contains(contract.ABI, "ERC721") {
		standard = "ERC721"
	} else if strings.Contains(contract.ABI, "ERC1155") {
		standard = "ERC1155"
	} else {
		return "", "", errors.New("unsupported contract standard")
	}
	//3.生成tokenId
	if req.TokenId == "" {
		//生成递增的tokenId
		count, err := dao.NFT.CountByContract(ctx, req.ContractId)
		if err != nil {
			return "", "", err
		}
		tokenId = big.NewInt(int64(count + 1)).String()
	} else {
		tokenId = req.TokenId
	}
	//4.上传图片到IPFS
	ipfsClient := ipfs.NewIPFSClient("http://localhost:5001", "http://localhost:8080")
	imageHash, err := ipfsClient.UploadFile(req.Image)

	//5.构建元数据
	metadata := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"image":       ipfsClient.GetGatewayURL(imageHash),
	}

	if req.Attributes != "" {
		var attrs []map[string]interface{}
		err = json.Unmarshal([]byte(req.Attributes), &attrs)
		if err != nil {
			return "", "", err
		}
		metadata["attributes"] = attrs
	}
	//6.上传元数据到IPFS
	metadataHash, err := ipfsClient.UploadJSON(metadata)
	if err != nil {
		return "", "", err
	}
	//7.调用合约mint方法
	client := ethclient.GetClient(ctx)
	parsed, err := abi.JSON(strings.NewReader(contract.ABI))
	if err != nil {
		return "", "", err
	}

	var data []byte
	if standard == "ERC721" {
		data, err = parsed.Pack("mint", common.HexToAddress(req.To), interface{}(new(big.Int).SetString(tokenId, 10)))
	} else {
		data, err = parsed.Pack("mint",
			common.HexToAddress(req.To),
			interface{}(new(big.Int).SetString(tokenId, 10)),
			new(big.Int).SetUint64(req.Amount),
			[]byte{},
		)
	}
	if err != nil {
		return "", "", err
	}

	//8.构建交易
	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(req.From))
	if err != nil {
		return "", "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", "", err
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: common.HexToAddress(req.From),
		To:   &common.HexToAddress(contract.Address),
		Data: data,
	})
	if err != nil {
		return "", "", err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(contract.Address), big.NewInt(0), gasLimit, gasPrice, data)

	//9.签名交易
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", "", err
	}

	wallet, err := dao.Wallet.GetByAddress(ctx, req.From)
	if err != nil {
		return "", "", err
	}

	privateKey, err := crypto.HexToECDSA(wallet.PrivateKey)
	if err != nil {
		return "", "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", "", err
	}

	// 发送交易
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", "", err
	}

	// 保存NFT信息
	nft := &model.NFT{
		ContractId:  req.ContractId,
		TokenId:     tokenId,
		Owner:       req.To,
		Creator:     req.From,
		URI:         ipfsClient.GetGatewayURL(metadataHash),
		Name:        req.Name,
		Description: req.Description,
		Image:       ipfsClient.GetGatewayURL(imageHash),
		Attributes:  req.Attributes,
		Standard:    standard,
		Amount:      req.Amount,
		Status:      1,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = dao.NFT.Insert(ctx, nft)
	if err != nil {
		return "", "", err
	}

	// 保存铸造记录
	transfer := &model.NFTTransfer{
		NftId:     nft.Id,
		From:      "0x0000000000000000000000000000000000000000",
		To:        req.To,
		Amount:    req.Amount,
		Type:      "mint",
		Hash:      signedTx.Hash().Hex(),
		CreatedAt: time.Now().Unix(),
	}

	err = dao.NFT.InsertTransfer(ctx, transfer)
	if err != nil {
		return "", "", err
	}

	return signedTx.Hash().Hex(), tokenId, nil
}

// Transfer 转移NFT
func (s *NFTLogic) Transfer(ctx context.Context, nftId uint64, from, to string, amount uint64) (hash string, err error) {
	// 获取NFT信息
	nft, err := dao.NFT.GetById(ctx, nftId)
	if err != nil {
		return "", err
	}

	// 获取合约信息
	contract, err := dao.Contract.GetById(ctx, nft.ContractId)
	if err != nil {
		return "", err
	}

	// 调用合约transfer方法
	client := ethclient.GetClient(ctx)
	parsed, err := abi.JSON(strings.NewReader(contract.ABI))
	if err != nil {
		return "", err
	}

	var data []byte
	if nft.Standard == "ERC721" {
		data, err = parsed.Pack("safeTransferFrom",
			common.HexToAddress(from),
			common.HexToAddress(to),
			interface{}(new(big.Int).SetString(nft.TokenId, 10)),
		)
	} else {
		data, err = parsed.Pack("safeTransferFrom",
			common.HexToAddress(from),
			common.HexToAddress(to),
			interface{}(new(big.Int).SetString(nft.TokenId, 10)),
			new(big.Int).SetUint64(amount),
			[]byte{},
		)
	}
	if err != nil {
		return "", err
	}

	// 构建交易
	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(from))
	if err != nil {
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: common.HexToAddress(from),
		To:   &common.HexToAddress(contract.Address),
		Data: data,
	})
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(contract.Address), big.NewInt(0), gasLimit, gasPrice, data)

	// 签名交易
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", err
	}

	wallet, err := dao.Wallet.GetByAddress(ctx, from)
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(wallet.PrivateKey)
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			return "", err
		}

		// 发送交易
		err = client.SendTransaction(ctx, signedTx)
		if err != nil {
			return "", err
		}

		// 更新NFT所有者
		err = dao.NFT.UpdateOwner(ctx, nftId, to)
		if err != nil {
			return "", err
		}

		// 保存转移记录
		transfer := &model.NFTTransfer{
			NftId:     nftId,
			From:      from,
			To:        to,
			Amount:    amount,
			Type:      "transfer",
			Hash:      signedTx.Hash().Hex(),
			CreatedAt: time.Now().Unix(),
		}

		err = dao.NFT.InsertTransfer(ctx, transfer)
		if err != nil {
			return "", err
		}

		return signedTx.Hash().Hex(), nil
	}
}

// List 获取NFT列表
func (s *NFTLogic) List(ctx context.Context, owner, creator string, contractId uint64, page, pageSize int) ([]*model.NFT, int, error) {
	return dao.NFT.GetList(ctx, owner, creator, contractId, page, pageSize)
}

// ListMarket 获取市场列表
func (s *NFTLogic) ListMarket(ctx context.Context, seller string, status int, page, pageSize int) ([]*model.NFTMarket, int, error) {
	return dao.NFT.GetMarketList(ctx, seller, status, page, pageSize)
}
