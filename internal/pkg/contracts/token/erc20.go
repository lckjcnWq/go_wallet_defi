package token

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC20合约ABI
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
        "inputs": [{"name": "owner", "type": "address"}],
        "name": "balanceOf",
        "outputs": [{"name": "", "type": "uint256"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {"name": "spender", "type": "address"},
            {"name": "amount", "type": "uint256"}
        ],
        "name": "approve",
        "outputs": [{"name": "", "type": "bool"}],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {"name": "owner", "type": "address"},
            {"name": "spender", "type": "address"}
        ],
        "name": "allowance",
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

// ERC20 represents an ERC20 contract
type ERC20 struct {
	address common.Address
	abi     abi.ABI
	client  *ethclient.Client
}

// NewERC20 creates a new instance of ERC20
func NewERC20(address common.Address, client *ethclient.Client) (*ERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}

	return &ERC20{
		address: address,
		abi:     parsed,
		client:  client,
	}, nil
}

// Name returns the token name
func (e *ERC20) Name() (string, error) {
	var result string
	err := e.call("name", &result)
	return result, err
}

// Symbol returns the token symbol
func (e *ERC20) Symbol() (string, error) {
	var result string
	err := e.call("symbol", &result)
	return result, err
}

// Decimals returns the token decimals
func (e *ERC20) Decimals() (uint8, error) {
	var result uint8
	err := e.call("decimals", &result)
	return result, err
}

// BalanceOf returns the token balance of the given address
func (e *ERC20) BalanceOf(owner common.Address) (*big.Int, error) {
	var result *big.Int
	err := e.call("balanceOf", &result, owner)
	return result, err
}

// Approve approves the given spender to spend tokens
func (e *ERC20) PackApprove(spender common.Address, amount *big.Int) ([]byte, error) {
	return e.abi.Pack("approve", spender, amount)
}

// Allowance returns the remaining tokens that spender is allowed to spend
func (e *ERC20) Allowance(owner, spender common.Address) (*big.Int, error) {
	var result *big.Int
	err := e.call("allowance", &result, owner, spender)
	return result, err
}

// Transfer transfers tokens to the given address
func (e *ERC20) PackTransfer(to common.Address, amount *big.Int) ([]byte, error) {
	return e.abi.Pack("transfer", to, amount)
}

// call executes a contract call
func (e *ERC20) call(method string, result interface{}, args ...interface{}) error {
	data, err := e.abi.Pack(method, args...)
	if err != nil {
		return err
	}

	msg := ethereum.CallMsg{
		To:   &e.address,
		Data: data,
	}

	output, err := e.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return err
	}

	return e.abi.UnpackIntoInterface(result, method, output)
}
