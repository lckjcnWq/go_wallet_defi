package logic

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/ethclientx"
	"math/big"
	"strings"
	"time"
)

type ContractLogic struct{}

// Deploy 部署合约
func (c ContractLogic) Deploy(ctx context.Context, name, abiStr, bytecode, from, argsJson, network string) (hash, address string, err error) {
	//1.解析ABI
	parsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return "", "", err
	}
	//2.解析构造函数
	var args []interface{}
	if argsJson != "" {
		if err := json.Unmarshal([]byte(argsJson), &args); err != nil {
			return "", "", err
		}
	}
	//3.编码构造函数
	data := common.FromHex(bytecode)
	if len(args) > 0 {
		input, err := parsed.Pack("", args...)
		if err != nil {
			return "", "", err
		}
		data = append(data, input...)
	}
	//4.获取nonce
	client := ethclientx.GetClient(ctx)
	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(from))
	if err != nil {
		return "", "", err
	}
	//5.获取gasPrice
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", "", err
	}
	//6.估算gas
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: common.HexToAddress(from),
		Data: data,
	})
	if err != nil {
		return "", "", err
	}
	//7.构建交易
	tx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, data)

	//8.签名交易
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", "", err
	}

	wallet, err := dao.Wallet.GetByAddress(ctx, from)
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

	//9.发送交易
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", "", err
	}
	//10.等待交易确认
	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return "", "", err
	}
	//11.保存合约信息
	contract := &model.Contract{
		Name:         name,
		Address:      receipt.ContractAddress.Hex(),
		ABI:          abiStr,
		Bytecode:     bytecode,
		Network:      network,
		DeployHash:   signedTx.Hash().Hex(),
		DeployHeight: receipt.BlockNumber.Int64(),
		Creator:      from,
		Status:       1,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	if err := dao.Contract.Insert(ctx, contract); err != nil {
		return "", "", err
	}

	return signedTx.Hash().Hex(), receipt.ContractAddress.Hex(), nil
}

// Call 调用合约
func (s *ContractLogic) Call(ctx context.Context, address, method, argsJson, from, value string) (interface{}, error) {
	// 获取合约信息
	contract, err := dao.Contract.GetByAddress(ctx, address, "ethereum")
	if err != nil {
		return nil, err
	}

	// 解析ABI
	parsed, err := abi.JSON(strings.NewReader(contract.ABI))
	if err != nil {
		return nil, err
	}

	// 解析参数
	var args []interface{}
	if argsJson != "" {
		if err := json.Unmarshal([]byte(argsJson), &args); err != nil {
			return nil, err
		}
	}

	// 编码调用数据
	data, err := parsed.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	client := ethclientx.GetClient(ctx)

	// 判断是否是只读方法
	if parsed.Methods[method].IsConstant() {
		// 只读方法直接调用
		result, err := client.CallContract(ctx, ethereum.CallMsg{
			From: common.HexToAddress(from),
			To:   &common.HexToAddress(address),
			Data: data,
		}, nil)
		if err != nil {
			return nil, err
		}

		// 解码返回值
		values, err := parsed.Unpack(method, result)
		if err != nil {
			return nil, err
		}
		return values, nil
	}

	// 写入方法需要发送交易
	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(from))
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	val := big.NewInt(0)
	if value != "" {
		val.SetString(value, 10)
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  common.HexToAddress(from),
		To:    &common.HexToAddress(address),
		Value: val,
		Data:  data,
	})
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(address), val, gasLimit, gasPrice, data)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	wallet, err := dao.Wallet.GetByAddress(ctx, from)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(wallet.PrivateKey)
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	// 保存调用记录
	call := &model.ContractCall{
		ContractId: contract.Id,
		Method:     method,
		Params:     argsJson,
		From:       from,
		Value:      value,
		Hash:       signedTx.Hash().Hex(),
		Status:     0,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	if err := dao.Contract.InsertCall(ctx, call); err != nil {
		return nil, err
	}

	// 发送交易
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx.Hash().Hex(), nil
}

// GetEvents 获取合约事件
func (s *ContractLogic) GetEvents(ctx context.Context, address, eventName string, fromBlock, toBlock int64, page, pageSize int) ([]*model.ContractEvent, int, error) {
	// 获取合约信息
	contract, err := dao.Contract.GetByAddress(ctx, address, "ethereum")
	if err != nil {
		return nil, 0, err
	}

	return dao.Contract.GetEvents(ctx, contract.Id, eventName, fromBlock, toBlock, page, pageSize)
}
