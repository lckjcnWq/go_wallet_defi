package controller

import (
	"context"
	v1 "go-wallet-defi/api/v1"
	"go-wallet-defi/internal/service"
)

type SecurityController struct{}

// CreateMultiSigWallet 创建多签钱包
func (c *SecurityController) CreateMultiSigWallet(ctx context.Context, req *v1.CreateMultiSigWalletReq) (res *v1.CreateMultiSigWalletRes, err error) {
	address, err := service.Security.CreateMultiSigWallet(
		ctx,
		req.ChainId,
		req.Name,
		req.Owners,
		req.Threshold,
		req.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &v1.CreateMultiSigWalletRes{
		Address: address,
	}, nil
}

// SubmitTransaction 提交多签交易
func (c *SecurityController) SubmitTransaction(ctx context.Context, req *v1.SubmitTransactionReq) (res *v1.SubmitTransactionRes, err error) {
	txId, err := service.Security.SubmitTransaction(
		ctx,
		req.WalletId,
		req.To,
		req.Value,
		req.Data,
		req.Description,
	)
	if err != nil {
		return nil, err
	}

	return &v1.SubmitTransactionRes{
		TxId: txId,
	}, nil
}

// ApproveTransaction 确认多签交易
func (c *SecurityController) ApproveTransaction(ctx context.Context, req *v1.ApproveTransactionReq) (res *v1.ApproveTransactionRes, err error) {
	err = service.Security.ApproveTransaction(
		ctx,
		req.TxId,
		req.Approver,
	)
	if err != nil {
		return nil, err
	}

	return &v1.ApproveTransactionRes{}, nil
}

// CreateWhitelist 创建白名单
func (c *SecurityController) CreateWhitelist(ctx context.Context, req *v1.CreateWhitelistReq) (res *v1.CreateWhitelistRes, err error) {
	err = service.Security.CreateWhitelist(
		ctx,
		req.UserId,
		req.Address,
		req.Name,
		req.Type,
		req.ExpireAt,
	)
	if err != nil {
		return nil, err
	}

	return &v1.CreateWhitelistRes{}, nil
}

// CreateTransactionLimit 创建交易限额
func (c *SecurityController) CreateTransactionLimit(ctx context.Context, req *v1.CreateTransactionLimitReq) (res *v1.CreateTransactionLimitRes, err error) {
	err = service.Security.CreateTransactionLimit(
		ctx,
		req.UserId,
		req.TokenAddress,
		req.SingleLimit,
		req.DailyLimit,
		req.WeeklyLimit,
		req.MonthlyLimit,
	)
	if err != nil {
		return nil, err
	}

	return &v1.CreateTransactionLimitRes{}, nil
}

// CreateRiskRule 创建风控规则
func (c *SecurityController) CreateRiskRule(ctx context.Context, req *v1.CreateRiskRuleReq) (res *v1.CreateRiskRuleRes, err error) {
	err = service.Security.CreateRiskRule(
		ctx,
		req.Name,
		req.Type,
		req.Content,
		req.Action,
	)
	if err != nil {
		return nil, err
	}

	return &v1.CreateRiskRuleRes{}, nil
}

// GetSecurityInfo 获取安全配置信息
func (c *SecurityController) GetSecurityInfo(ctx context.Context, req *v1.GetSecurityInfoReq) (res *v1.GetSecurityInfoRes, err error) {
	info, err := service.Security.GetSecurityInfo(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &v1.GetSecurityInfoRes{
		Info: info,
	}, nil
}
