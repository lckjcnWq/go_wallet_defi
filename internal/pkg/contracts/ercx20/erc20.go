package ercx20

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

// 合约ABI
const ERC20ABI = `[
    {
        "constant": true,
        "inputs": [],
        "name": "name",
        "outputs": [{"name": "", "type": "string"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [{"name": "", "type": "string"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [{"name": "", "type": "uint8"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [{"name": "account", "type": "address"}],
        "name": "balanceOf",
        "outputs": [{"name": "", "type": "uint256"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {"name": "recipient", "type": "address"},
            {"name": "amount", "type": "uint256"}
        ],
        "name": "transfer",
        "outputs": [{"name": "", "type": "bool"}],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }
]`

// ERC20 合约实例
type ERC20 struct {
	abi     abi.ABI
	address common.Address
	backend bind.ContractBackend
}

// NewERC20 创建ERC20合约实例
func NewERC20(address common.Address, backend *ethclient.Client) (*ERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}

	return &ERC20{
		abi:     parsed,
		address: address,
		backend: backend,
	}, nil
}

// Name 获取代币名称
func (e *ERC20) Name() (string, error) {
	var result string
	err := e.abi.UnpackIntoInterface(&result, "name", nil)
	return result, err
}

// Symbol 获取代币符号
func (e *ERC20) Symbol() (string, error) {
	var result string
	err := e.abi.UnpackIntoInterface(&result, "symbol", nil)
	return result, err
}

// Decimals 获取代币精度
func (e *ERC20) Decimals() (uint8, error) {
	var result uint8
	err := e.abi.UnpackIntoInterface(&result, "decimals", nil)
	return result, err
}

// BalanceOf 查询代币余额
func (e *ERC20) BalanceOf(account common.Address) (*big.Int, error) {
	var result *big.Int
	err := e.abi.UnpackIntoInterface(&result, "balanceOf", account.Bytes())
	return result, err
}

// Transfer 转账方法
func (e *ERC20) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return nil, nil // 这个方法我们不需要，我们用PackTransfer
}

// PackTransfer 打包transfer方法调用数据
func (e *ERC20) PackTransfer(to common.Address, amount *big.Int) ([]byte, error) {
	return e.abi.Pack("transfer", to, amount)
}
