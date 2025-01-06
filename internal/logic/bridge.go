package logic

import (
	"context"
	_ "crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/contracts/bridge"
	"go-wallet-defi/internal/pkg/contracts/token"
	"go-wallet-defi/internal/pkg/ethclientx"
	"math/big"
	"strings"
	"time"
)

type BridgeLogic struct{}

func (s BridgeLogic) CrossTransfer(ctx context.Context, fromChainId, toChainId uint64, fromAddress, toAddress, tokenAddress, amount string) (hash string, nonce uint64, err error) {
	//1.获取来源链信息
	fromChain, err := dao.Chain.GetById(ctx, fromChainId)
	if err != nil {
		return "", 0, err
	}
	//2.获取目标链信息
	toChain, err := dao.Chain.GetById(ctx, toChainId)
	if err != nil {
		return "", 0, err
	}
	//3.检查跨链桥合约地址
	if fromChain.BridgeAddress == "" || toChain.BridgeAddress == "" {
		return "", 0, errors.New("bridge contract not configured")
	}
	//4.获取客户端
	client, err := ethclientx.GetClientByChainId(ctx, fromChain.ChainId)
	if err != nil {
		return "", 0, err
	}
	//5.解析ABI
	parsed, err := abi.JSON(strings.NewReader(bridge.BridgeABI))
	if err != nil {
		return "", 0, err
	}
	//6.如果是代币，需要先approve
	if tokenAddress != "" {
		//6.1获取代币合约
		token, err := token.NewERC20(common.HexToAddress(tokenAddress), client)
		if err != nil {
			return "", 0, err
		}

		// 调用approve方法
		approveData, err := token.PackApprove(common.HexToAddress(fromChain.BridgeAddress), new(big.Int).SetString(amount, 10))
		if err != nil {
			return "", 0, err
		}

		//6.2发送approve交易
		approveHash, err := s.sendTransaction(ctx, client, fromAddress, tokenAddress, big.NewInt(0), approveData)
		if err != nil {
			return "", 0, err
		}

		//6.3等待approve交易确认
		_, err = s.waitTransaction(ctx, client, approveHash)
		if err != nil {
			return "", 0, err
		}
	}

	// 构造lock方法调用数据
	value := big.NewInt(0)
	if tokenAddress == "" {
		value, _ = new(big.Int).SetString(amount, 10)
		tokenAddress = "0x0000000000000000000000000000000000000000"
	}

	data, err := parsed.Pack("lock",
		common.HexToAddress(tokenAddress),
		interface{}(new(big.Int).SetString(amount, 10)),
		new(big.Int).SetUint64(toChainId),
		common.HexToAddress(toAddress),
		new(big.Int).SetUint64(nonce),
	)
	if err != nil {
		return "", 0, err
	}

	// 发送交易
	hash, err = s.sendTransaction(ctx, client, fromAddress, fromChain.BridgeAddress, value, data)
	if err != nil {
		return "", 0, err
	}

	// 保存跨链交易记录
	transfer := &model.CrossTransfer{
		FromChainId:  fromChainId,
		ToChainId:    toChainId,
		FromAddress:  fromAddress,
		ToAddress:    toAddress,
		TokenAddress: tokenAddress,
		Amount:       amount,
		Fee:          "0",
		Nonce:        nonce,
		FromHash:     hash,
		Status:       0,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	err = dao.Chain.InsertCrossTransfer(ctx, transfer)
	if err != nil {
		return "", 0, err
	}

	return hash, nonce, nil
}

// GetCrossTransfers 获取跨链交易列表
func (s *BridgeLogic) GetCrossTransfers(ctx context.Context, fromChainId, toChainId uint64, address string, status, page, pageSize int) ([]*model.CrossTransfer, int, error) {
	return dao.Chain.GetCrossTransferList(ctx, fromChainId, toChainId, address, status, page, pageSize)
}

// ProcessLockEvent 处理锁定事件
func (s *BridgeLogic) ProcessLockEvent(ctx context.Context, chainId uint64, token, from string, amount string, toChainId uint64, toAddress string, nonce uint64, hash string) error {
	// 查找跨链交易记录
	transfers, _, err := dao.Chain.GetCrossTransferList(ctx, chainId, toChainId, "", 0, 1, 10)
	if err != nil {
		return err
	}

	var transfer *model.CrossTransfer
	for _, t := range transfers {
		if t.Nonce == nonce {
			transfer = t
			break
		}
	}
	if transfer == nil {
		return errors.New("transfer not found")
	}

	// 更新状态为已锁定
	err = dao.Chain.UpdateCrossTransfer(ctx, transfer.Id, g.Map{
		"status":     1,
		"from_hash":  hash,
		"updated_at": time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	// 获取目标链信息
	toChain, err := dao.Chain.GetByChainId(ctx, toChainId)
	if err != nil {
		return err
	}

	// 获取验证者私钥
	privateKey, err := crypto.HexToECDSA("validator private key")
	if err != nil {
		return err
	}

	// 计算解锁签名
	message := crypto.Keccak256(
		common.HexToAddress(token).Bytes(),
		new(big.Int).SetString(amount, 10).Bytes(),
		common.HexToAddress(toAddress).Bytes(),
		new(big.Int).SetUint64(chainId).Bytes(),
		new(big.Int).SetUint64(nonce).Bytes(),
	)

	signature, err := crypto.Sign(message, privateKey)
	if err != nil {
		return err
	}

	// 获取目标链客户端
	client, err := ethclientx.GetClientByChainId(ctx, toChainId)
	if err != nil {
		return err
	}

	// 解析ABI
	parsed, err := abi.JSON(strings.NewReader(bridge.BridgeABI))
	if err != nil {
		return err
	}

	// 构造unlock方法调用数据
	data, err := parsed.Pack("unlock",
		common.HexToAddress(token),
		new(big.Int).SetString(amount, 10),
		common.HexToAddress(toAddress),
		new(big.Int).SetUint64(chainId),
		new(big.Int).SetUint64(nonce),
		signature,
	)
	if err != nil {
		return err
	}

	// 发送解锁交易
	validator := crypto.PubkeyToAddress(privateKey.PublicKey)
	toHash, err := s.sendTransaction(ctx, client, validator.Hex(), toChain.BridgeAddress, big.NewInt(0), data)
	if err != nil {
		return err
	}

	// 更新解锁交易哈希
	return dao.Chain.UpdateCrossTransfer(ctx, transfer.Id, g.Map{
		"to_hash":    toHash,
		"updated_at": time.Now().Unix(),
	})
}

// ProcessUnlockEvent 处理解锁事件
func (s *BridgeLogic) ProcessUnlockEvent(ctx context.Context, chainId uint64, token, to string, amount string, fromChainId uint64, nonce uint64, hash string) error {
	// 查找跨链交易记录
	transfers, _, err := dao.Chain.GetCrossTransferList(ctx, fromChainId, chainId, "", 1, 1, 10)
	if err != nil {
		return err
	}

	var transfer *model.CrossTransfer
	for _, t := range transfers {
		if t.Nonce == nonce {
			transfer = t
			break
		}
	}
	if transfer == nil {
		return errors.New("transfer not found")
	}

	// 更新状态为已完成
	return dao.Chain.UpdateCrossTransfer(ctx, transfer.Id, g.Map{
		"status":     3,
		"to_hash":    hash,
		"updated_at": time.Now().Unix(),
	})
}

// 发送交易
func (s *BridgeLogic) sendTransaction(ctx context.Context, client *ethclient.Client, from, to string, value *big.Int, data []byte) (string, error) {
	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(from))
	if err != nil {
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  common.HexToAddress(from),
		To:    &common.HexToAddress(to),
		Value: value,
		Data:  data,
	})
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(to), value, gasLimit, gasPrice, data)

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
		return "", err
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// 等待交易确认
func (s *BridgeLogic) waitTransaction(ctx context.Context, client *ethclient.Client, hash string) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(ctx, common.HexToHash(hash))
		if err != nil {
			if err == ethereum.NotFound {
				time.Sleep(time.Second)
				continue
			}
			return nil, err
		}
		return receipt, nil
	}
}
