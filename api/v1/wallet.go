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

type ImportWalletReq struct {
	g.Meta   `path:"/wallet/import" method:"post" tags:"钱包管理" summary:"导入钱包"`
	Chain    string `v:"required" dc:"链名称 ETH/BSC"`
	Mnemonic string `v:"required" dc:"助记词"`
}

type ImportWalletRes struct {
	Address string `json:"address" dc:"钱包地址"`
}

// 钱包列表
type ListWalletReq struct {
	g.Meta   `path:"/wallet/list" method:"get" tags:"钱包管理" summary:"钱包列表"`
	Page     int `json:"page" d:"1" v:"min:1" dc:"页码"`
	PageSize int `json:"pageSize" d:"10" v:"max:50" dc:"每页数量"`
}

type ListWalletRes struct {
	List []WalletInfo `json:"list" dc:"钱包列表"`
}

type WalletInfo struct {
	Id        uint64 `json:"id"`
	Address   string `json:"address"`
	Chain     string `json:"chain"`
	Balance   string `json:"balance"`
	CreatedAt int64  `json:"createdAt"`
}

// 查询余额
type GetBalanceReq struct {
	g.Meta  `path:"/wallet/balance" method:"get" tags:"钱包管理" summary:"查询余额"`
	Address string `v:"required" dc:"钱包地址"`
}

type GetBalanceRes struct {
	Balance string `json:"balance" dc:"余额"`
}
