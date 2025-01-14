package v1

import "github.com/gogf/gf/v2/frame/g"

type BatchTransactionsReq struct {
	g.Meta  `path:"/gas/batch" method:"post"`
	ChainId uint64       `json:"chain_id" v:"required"`
	From    string       `json:"from"     v:"required"`
	Calls   []model.Call `json:"calls"    v:"required"`
}

type BatchTransactionsRes struct {
	TxHash string `json:"tx_hash"`
}

type EstimateGasPriceReq struct {
	g.Meta   `path:"/gas/estimate" method:"get"`
	ChainId  uint64 `json:"chain_id" v:"required"`
	Strategy string `json:"strategy" v:"required"`
}

type EstimateGasPriceRes struct {
	GasPrice string `json:"gas_price"`
}

type AccelerateTransactionReq struct {
	g.Meta  `path:"/gas/accelerate" method:"post"`
	ChainId uint64 `json:"chain_id" v:"required"`
	TxHash  string `json:"tx_hash"  v:"required"`
}

type AccelerateTransactionRes struct {
	NewTxHash string `json:"new_tx_hash"`
}

type CancelTransactionReq struct {
	g.Meta  `path:"/gas/cancel" method:"post"`
	ChainId uint64 `json:"chain_id" v:"required"`
	TxHash  string `json:"tx_hash"  v:"required"`
}

type CancelTransactionRes struct {
	NewTxHash string `json:"new_tx_hash"`
}

type GetGasInfoReq struct {
	g.Meta  `path:"/gas/info" method:"get"`
	ChainId uint64 `json:"chain_id" v:"required"`
}

type GetGasInfoRes struct {
	Info map[string]interface{} `json:"info"`
}
