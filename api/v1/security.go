package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type CreateMultiSigWalletReq struct {
	g.Meta    `path:"/security/wallet/create" method:"post"`
	ChainId   uint64   `json:"chain_id"   v:"required"`
	Name      string   `json:"name"       v:"required"`
	Owners    []string `json:"owners"     v:"required"`
	Threshold int      `json:"threshold"  v:"required"`
	CreatedBy string   `json:"created_by" v:"required"`
}

type CreateMultiSigWalletRes struct {
	Address string `json:"address"`
}

type SubmitTransactionReq struct {
	g.Meta      `path:"/security/tx/submit" method:"post"`
	WalletId    uint64 `json:"wallet_id"   v:"required"`
	To          string `json:"to"          v:"required"`
	Value       string `json:"value"       v:"required"`
	Data        []byte `json:"data"`
	Description string `json:"description"`
}

type SubmitTransactionRes struct {
	TxId uint64 `json:"tx_id"`
}

type ApproveTransactionReq struct {
	g.Meta   `path:"/security/tx/approve" method:"post"`
	TxId     uint64 `json:"tx_id"     v:"required"`
	Approver string `json:"approver"   v:"required"`
}

type ApproveTransactionRes struct{}

type CreateWhitelistReq struct {
	g.Meta   `path:"/security/whitelist/create" method:"post"`
	UserId   uint64     `json:"user_id"    v:"required"`
	Address  string     `json:"address"    v:"required"`
	Name     string     `json:"name"       v:"required"`
	Type     int        `json:"type"       v:"required"`
	ExpireAt *time.Time `json:"expire_at"`
}

type CreateWhitelistRes struct{}

type CreateTransactionLimitReq struct {
	g.Meta       `path:"/security/limit/create" method:"post"`
	UserId       uint64 `json:"user_id"        v:"required"`
	TokenAddress string `json:"token_address"  v:"required"`
	SingleLimit  string `json:"single_limit"   v:"required"`
	DailyLimit   string `json:"daily_limit"    v:"required"`
	WeeklyLimit  string `json:"weekly_limit"   v:"required"`
	MonthlyLimit string `json:"monthly_limit"  v:"required"`
}

type CreateTransactionLimitRes struct{}

type CreateRiskRuleReq struct {
	g.Meta  `path:"/security/risk/create" method:"post"`
	Name    string                 `json:"name"    v:"required"`
	Type    int                    `json:"type"    v:"required"`
	Content map[string]interface{} `json:"content" v:"required"`
	Action  int                    `json:"action"  v:"required"`
}

type CreateRiskRuleRes struct{}

type GetSecurityInfoReq struct {
	g.Meta `path:"/security/info" method:"get"`
	UserId uint64 `json:"user_id" v:"required"`
}

type GetSecurityInfoRes struct {
	Info map[string]interface{} `json:"info"`
}
