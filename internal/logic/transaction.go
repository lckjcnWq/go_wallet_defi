package logic

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/contracts/ercx20"
	"go-wallet-defi/internal/pkg/ethclientx"
	"math/big"
	"time"
)

type TransactionLogic struct{}

// TransferEth 转账ETH
func (s *TransactionLogic) TransferEth(ctx context.Context, from, to, amount string, gasPrice string, gasLimit uint64) (string, error) {
	client := ethclientx.GetClient(ctx)
	//1.解析地址
	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	//2.获取nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", err
	}
	//3.解析金额
	value := new(big.Int)
	value.SetString(amount, 10)

	//4.获取gas价格
	var gp *big.Int
	if gasPrice == "" {
		gp, err = client.SuggestGasPrice(ctx)
		if err != nil {
			return "", err
		}
	} else {
		gp = new(big.Int)
		gp.SetString(gasPrice, 10)
	}

	//5.使用默认gas限制
	if gasLimit == 0 {
		gasLimit = 21000
	}

	//6.构建交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gp, nil)
	//7.获取chainID
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return "", err
	}
	//8.获取私钥
	wallet, err := dao.Wallet.GetByAddress(ctx, from)
	privateKey, err := crypto.HexToECDSA(wallet.PrivateKey)
	if err != nil {
		return "", err
	}
	//9.签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}
	//10.发送交易
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	//11.保存交易记录
	transaction := &model.Transaction{
		UserId:      wallet.UserId,
		FromAddress: from,
		ToAddress:   to,
		Amount:      amount,
		Hash:        signedTx.Hash().Hex(),
		Nonce:       nonce,
		GasPrice:    gp.String(),
		GasLimit:    gasLimit,
		Status:      0,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = dao.Transaction.Insert(ctx, transaction)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// TransferToken 转账代币
func (s *TransactionLogic) TransferToken(ctx context.Context, from, to, amount, token string, gasPrice string, gasLimit uint64) (string, error) {
	client := ethclientx.GetClient(ctx)

	// 解析合约地址
	tokenAddress := common.HexToAddress(token)

	// 创建合约实例
	instance, err := ercx20.NewERC20(tokenAddress, client)
	if err != nil {
		return "", err
	}

	// 解析转账参数
	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)
	value := new(big.Int)
	value.SetString(amount, 10)

	// 获取nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", err
	}

	// 获取gas价格
	var gp *big.Int
	if gasPrice == "" {
		gp, err = client.SuggestGasPrice(ctx)
		if err != nil {
			return "", err
		}
	} else {
		gp = new(big.Int)
		gp.SetString(gasPrice, 10)
	}

	// 获取私钥
	wallet, err := dao.Wallet.GetByAddress(ctx, from)
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(wallet.PrivateKey)
	if err != nil {
		return "", err
	}

	// 构建交易数据
	data, err := instance.PackTransfer(toAddress, value)
	if err != nil {
		return "", err
	}

	// 预估gas
	if gasLimit == 0 {
		gasLimit, err = client.EstimateGas(ctx, ethereum.CallMsg{
			From: fromAddress,
			To:   &tokenAddress,
			Data: data,
		})
		if err != nil {
			return "", err
		}
		// 增加一点buffer
		gasLimit = gasLimit * 12 / 10
	}

	// 构建交易
	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), gasLimit, gp, data)

	// 获取chainID
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return "", err
	}

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// 发送交易
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	// 保存交易记录
	transaction := &model.Transaction{
		UserId:       wallet.UserId,
		FromAddress:  from,
		ToAddress:    to,
		Amount:       amount,
		TokenAddress: token,
		Hash:         signedTx.Hash().Hex(),
		Nonce:        nonce,
		GasPrice:     gp.String(),
		GasLimit:     gasLimit,
		Data:         common.Bytes2Hex(data),
		Status:       0,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	err = dao.Transaction.Insert(ctx, transaction)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// GetTransactions 获取交易记录
func (s *TransactionLogic) GetTransactions(ctx context.Context, address string, page, pageSize int) ([]*model.Transaction, int, error) {
	return dao.Transaction.GetList(ctx, address, page, pageSize)
}

// UpdateTransactionStatus 更新交易状态
func (s *TransactionLogic) UpdateTransactionStatus(ctx context.Context, hash string) error {
	client := ethclientx.GetClient(ctx)

	// 获取交易收据
	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(hash))
	if err != nil {
		return err
	}

	// 更新状态
	data := g.Map{
		"status":       receipt.Status,
		"block_number": receipt.BlockNumber.Int64(),
		"updated_at":   time.Now().Unix(),
	}

	// 获取区块时间
	block, err := client.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return err
	}
	data["block_time"] = block.Time()

	return dao.Transaction.Update(ctx, hash, data)
}
