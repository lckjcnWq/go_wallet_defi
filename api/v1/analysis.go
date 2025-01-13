package v1

import "github.com/gogf/gf/v2/frame/g"

type AnalyzeAddressReq struct {
	g.Meta  `path:"/analysis/analyze" method:"post"`
	ChainId uint64 `json:"chain_id" v:"required"`
	Address string `json:"address"  v:"required"`
}

type AnalyzeAddressRes struct{}

type GetAnalysisReportReq struct {
	g.Meta  `path:"/analysis/report" method:"get"`
	Address string `json:"address" v:"required"`
}

type GetAnalysisReportRes struct {
	Report map[string]interface{} `json:"report"`
}

type UpdateMarketDataReq struct {
	g.Meta `path:"/analysis/market/update" method:"post"`
}

type UpdateMarketDataRes struct{}
