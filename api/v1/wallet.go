package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CreateWalletReq struct {
	g.Meta `path:"/wallet/create" method:"post" tags:"钱包管理" summary:"创建钱包"`
	Chain  string `v:"required" dc:"链名称 ETH/BSC"`
}

type CreateWalletRes struct {
	Address  string `json:"address" dc:"钱包地址"`
	Mnemonic string `json:"mnemonic" dc:"助记词"`
}