package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"math/big"
	"strings"
	"time"
)

type SecurityService struct{}

var Security = &SecurityService{}

// CreateMultiSigWallet 创建多签钱包
func (s *SecurityService) CreateMultiSigWallet(ctx context.Context, chainId uint64, name string, owners []string, threshold int, createdBy string) (address string, err error) {
	// 部署多签合约
	client, err := ethclient.GetClient(chainId)
	if err != nil {
		return "", err
	}

	// 编码构造函数参数
	data, err := contracts.PackMultiSigConstructor(owners, big.NewInt(int64(threshold)))
	if err != nil {
		return "", err
	}

	// 部署合约
	contractAddr, tx, err := client.DeployContract(ctx, data)
	if err != nil {
		return "", err
	}

	// 等待交易确认
	_, err = client.WaitTransaction(ctx, tx.Hash())
	if err != nil {
		return "", err
	}

	// 保存钱包信息
	ownersJson, _ := json.Marshal(owners)
	wallet := &model.MultiSigWallet{
		ChainId:   chainId,
		Address:   contractAddr.Hex(),
		Name:      name,
		Owners:    string(ownersJson),
		Threshold: threshold,
		CreatedBy: createdBy,
		Status:    1,
	}

	err = dao.Security.CreateMultiSigWallet(ctx, wallet)
	if err != nil {
		return "", err
	}

	return contractAddr.Hex(), nil
}

// SubmitTransaction 提交多签交易
func (s *SecurityService) SubmitTransaction(ctx context.Context, walletId uint64, to string, value string, data []byte, description string) (uint64, error) {
	// 获取钱包信息
	wallet, err := dao.Security.GetMultiSigWallet(ctx, walletId)
	if err != nil {
		return 0, err
	}

	// 校验权限
	var owners []string
	err = json.Unmarshal([]byte(wallet.Owners), &owners)
	if err != nil {
		return 0, err
	}

	// 调用合约提交交易
	client, err := ethclient.GetClient(wallet.ChainId)
	if err != nil {
		return 0, err
	}

	valueBig, _ := new(big.Int).SetString(value, 10)
	txData, err := contracts.PackMultiSigSubmit(to, valueBig, data)
	if err != nil {
		return 0, err
	}

	tx, err := client.SendTransaction(ctx, wallet.Address, big.NewInt(0), txData)
	if err != nil {
		return 0, err
	}

	// 保存交易记录
	multiSigTx := &model.MultiSigTransaction{
		WalletId:    walletId,
		ChainId:     wallet.ChainId,
		TxHash:      tx.Hash().Hex(),
		To:          to,
		Value:       value,
		Data:        string(data),
		Description: description,
		Status:      0,
	}

	err = dao.Security.CreateMultiSigTx(ctx, multiSigTx)
	if err != nil {
		return 0, err
	}

	return multiSigTx.Id, nil
}

// ApproveTransaction 确认多签交易
func (s *SecurityService) ApproveTransaction(ctx context.Context, txId uint64, approver string) error {
	// 获取交易信息
	tx, err := dao.Security.GetMultiSigTx(ctx, txId)
	if err != nil {
		return err
	}

	// 获取钱包信息
	wallet, err := dao.Security.GetMultiSigWallet(ctx, tx.WalletId)
	if err != nil {
		return err
	}

	// 校验权限
	var owners []string
	err = json.Unmarshal([]byte(wallet.Owners), &owners)
	if err != nil {
		return err
	}

	isOwner := false
	for _, owner := range owners {
		if owner == approver {
			isOwner = true
			break
		}
	}
	if !isOwner {
		return errors.New("not owner")
	}

	// 检查是否已确认
	var confirmations []string
	if tx.Confirmations != "" {
		err = json.Unmarshal([]byte(tx.Confirmations), &confirmations)
		if err != nil {
			return err
		}
	}

	for _, confirmation := range confirmations {
		if confirmation == approver {
			return errors.New("already confirmed")
		}
	}

	// 调用合约确认交易
	client, err := ethclient.GetClient(wallet.ChainId)
	if err != nil {
		return err
	}

	approveData, err := contracts.PackMultiSigApprove(big.NewInt(int64(txId)))
	if err != nil {
		return err
	}

	approveTx, err := client.SendTransaction(ctx, wallet.Address, big.NewInt(0), approveData)
	if err != nil {
		return err
	}

	// 等待交易确认
	_, err = client.WaitTransaction(ctx, approveTx.Hash())
	if err != nil {
		return err
	}

	// 更新确认信息
	confirmations = append(confirmations, approver)
	confirmationsJson, _ := json.Marshal(confirmations)
	err = dao.Security.UpdateMultiSigTx(ctx, txId, map[string]interface{}{
		"confirmations": string(confirmationsJson),
		"status":        1,
	})
	if err != nil {
		return err
	}

	// 检查是否达到阈值可以执行
	if len(confirmations) >= wallet.Threshold {
		err = s.ExecuteTransaction(ctx, txId)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExecuteTransaction 执行多签交易
func (s *SecurityService) ExecuteTransaction(ctx context.Context, txId uint64) error {
	// 获取交易信息
	tx, err := dao.Security.GetMultiSigTx(ctx, txId)
	if err != nil {
		return err
	}

	// 获取钱包信息
	wallet, err := dao.Security.GetMultiSigWallet(ctx, tx.WalletId)
	if err != nil {
		return err
	}

	// 调用合约执行交易
	client, err := ethclient.GetClient(wallet.ChainId)
	if err != nil {
		return err
	}

	executeData, err := contracts.PackMultiSigExecute(big.NewInt(int64(txId)))
	if err != nil {
		return err
	}

	executeTx, err := client.SendTransaction(ctx, wallet.Address, big.NewInt(0), executeData)
	if err != nil {
		return err
	}

	// 等待交易确认
	_, err = client.WaitTransaction(ctx, executeTx.Hash())
	if err != nil {
		return err
	}

	// 更新交易状态
	err = dao.Security.UpdateMultiSigTx(ctx, txId, map[string]interface{}{
		"status":      2,
		"executed_at": time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

// CreateWhitelist 创建白名单
func (s *SecurityService) CreateWhitelist(ctx context.Context, userId uint64, address, name string, txType int, expireAt *time.Time) error {
	whitelist := &model.Whitelist{
		UserId:   userId,
		Address:  address,
		Name:     name,
		Type:     txType,
		ExpireAt: expireAt,
		Status:   1,
	}
	return dao.Security.CreateWhitelist(ctx, whitelist)
}

// CreateTransactionLimit 创建交易限额
func (s *SecurityService) CreateTransactionLimit(ctx context.Context, userId uint64, tokenAddress string, singleLimit, dailyLimit, weeklyLimit, monthlyLimit string) error {
	limit := &model.TransactionLimit{
		UserId:       userId,
		TokenAddress: tokenAddress,
		SingleLimit:  singleLimit,
		DailyLimit:   dailyLimit,
		WeeklyLimit:  weeklyLimit,
		MonthlyLimit: monthlyLimit,
		Status:       1,
	}
	return dao.Security.CreateTransactionLimit(ctx, limit)
}

// CheckTransactionLimit 检查交易限额
func (s *SecurityService) CheckTransactionLimit(ctx context.Context, userId uint64, tokenAddress string, amount string) error {
	// 获取限额配置
	limit, err := dao.Security.GetTransactionLimit(ctx, userId)
	if err != nil {
		return err
	}

	// 检查单笔限额
	amountBig, _ := new(big.Int).SetString(amount, 10)
	singleLimitBig, _ := new(big.Int).SetString(limit.SingleLimit, 10)
	if amountBig.Cmp(singleLimitBig) > 0 {
		return errors.New("exceed single limit")
	}

	// 检查日限额
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	todayEnd := todayStart.Add(24 * time.Hour)

	dailyAmount, err := dao.Security.GetUserTxAmount(ctx, userId, tokenAddress, todayStart.Unix(), todayEnd.Unix())
	if err != nil {
		return err
	}

	dailyAmountBig, _ := new(big.Int).SetString(dailyAmount, 10)
	dailyLimitBig, _ := new(big.Int).SetString(limit.DailyLimit, 10)

	total := new(big.Int).Add(dailyAmountBig, amountBig)
	if total.Cmp(dailyLimitBig) > 0 {
		return errors.New("exceed daily limit")
	}

	// 检查周限额
	weekStart := todayStart.AddDate(0, 0, -int(todayStart.Weekday()))
	weekAmount, err := dao.Security.GetUserTxAmount(ctx, userId, tokenAddress, weekStart.Unix(), todayEnd.Unix())
	if err != nil {
		return err
	}

	weekAmountBig, _ := new(big.Int).SetString(weekAmount, 10)
	weeklyLimitBig, _ := new(big.Int).SetString(limit.WeeklyLimit, 10)

	total = new(big.Int).Add(weekAmountBig, amountBig)
	if total.Cmp(weeklyLimitBig) > 0 {
		return errors.New("exceed weekly limit")
	}

	// 检查月限额
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	monthAmount, err := dao.Security.GetUserTxAmount(ctx, userId, tokenAddress, monthStart.Unix(), todayEnd.Unix())
	if err != nil {
		return err
	}

	monthAmountBig, _ := new(big.Int).SetString(monthAmount, 10)
	monthlyLimitBig, _ := new(big.Int).SetString(limit.MonthlyLimit, 10)

	total = new(big.Int).Add(monthAmountBig, amountBig)
	if total.Cmp(monthlyLimitBig) > 0 {
		return errors.New("exceed monthly limit")
	}

	return nil
}

// CreateRiskRule 创建风控规则
func (s *SecurityService) CreateRiskRule(ctx context.Context, name string, ruleType int, content map[string]interface{}, action int) error {
	contentJson, _ := json.Marshal(content)
	rule := &model.RiskRule{
		Name:    name,
		Type:    ruleType,
		Content: string(contentJson),
		Action:  action,
		Status:  1,
	}
	return dao.Security.CreateRiskRule(ctx, rule)
}

// CheckRiskRules 检查风控规则
func (s *SecurityService) CheckRiskRules(ctx context.Context, userId uint64, params map[string]interface{}) error {
	// 获取所有启用的规则
	rules, err := dao.Security.GetActiveRiskRules(ctx)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		var content map[string]interface{}
		err = json.Unmarshal([]byte(rule.Content), &content)
		if err != nil {
			continue
		}

		var matched bool
		switch rule.Type {
		case 1: // 地址黑名单
			matched = s.checkAddressBlacklist(params, content)
		case 2: // 金额预警
			matched = s.checkAmountLimit(params, content)
		case 3: // 频率限制
			matched = s.checkFrequencyLimit(ctx, userId, params, content)
		}

		if matched {
			// 记录风控日志
			log := &model.RiskLog{
				UserId:  userId,
				RuleId:  rule.Id,
				Type:    rule.Type,
				Content: fmt.Sprintf("触发规则: %s, 参数: %v", rule.Name, params),
				Result:  rule.Action,
			}
			dao.Security.CreateRiskLog(ctx, log)

			// 根据动作处理
			switch rule.Action {
			case 1: // 拒绝
				return errors.New("transaction rejected by risk control")
			case 2: // 预警
				// 发送预警通知
				s.sendRiskAlert(ctx, userId, rule, params)
			case 3: // 人工审核
				// 创建审核任务
				s.createAuditTask(ctx, userId, rule, params)
			}
		}
	}

	return nil
}

// 检查地址黑名单
func (s *SecurityService) checkAddressBlacklist(params, rule map[string]interface{}) bool {
	addresses, ok := rule["addresses"].([]interface{})
	if !ok {
		return false
	}

	txAddress, ok := params["address"].(string)
	if !ok {
		return false
	}

	for _, addr := range addresses {
		if strings.EqualFold(addr.(string), txAddress) {
			return true
		}
	}

	return false
}

// 检查金额预警
func (s *SecurityService) checkAmountLimit(params, rule map[string]interface{}) bool {
	threshold, ok := rule["threshold"].(string)
	if !ok {
		return false
	}

	amount, ok := params["amount"].(string)
	if !ok {
		return false
	}

	thresholdBig, _ := new(big.Int).SetString(threshold, 10)
	amountBig, _ := new(big.Int).SetString(amount, 10)

	return amountBig.Cmp(thresholdBig) >= 0
}

// 检查频率限制
func (s *SecurityService) checkFrequencyLimit(ctx context.Context, userId uint64, params, rule map[string]interface{}) bool {
	limit, ok := rule["limit"].(float64)
	if !ok {
		return false
	}

	duration, ok := rule["duration"].(float64)
	if !ok {
		return false
	}

	// 获取时间范围内的交易次数
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(duration) * time.Second)

	count, err := dao.Security.GetUserTxCount(ctx, userId, startTime.Unix(), endTime.Unix())
	if err != nil {
		return false
	}

	return float64(count) >= limit
}

// 发送风险预警
func (s *SecurityService) sendRiskAlert(ctx context.Context, userId uint64, rule *model.RiskRule, params map[string]interface{}) {
	// TODO: 实现预警通知发送
	// 1. 邮件通知
	// 2. 短信通知
	// 3. 站内信
}

// 创建审核任务
func (s *SecurityService) createAuditTask(ctx context.Context, userId uint64, rule *model.RiskRule, params map[string]interface{}) {
	// TODO: 创建人工审核任务
}

// GetSecurityInfo 获取安全配置信息
func (s *SecurityService) GetSecurityInfo(ctx context.Context, userId uint64) (map[string]interface{}, error) {
	// 获取多签钱包
	wallets, err := dao.Security.GetUserMultiSigWallets(ctx, userId)
	if err != nil {
		return nil, err
	}

	// 获取白名单
	whitelists, err := dao.Security.GetUserWhitelists(ctx, userId)
	if err != nil {
		return nil, err
	}

	// 获取交易限额
	limits, err := dao.Security.GetUserTransactionLimits(ctx, userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"multi_sig_wallets": wallets,
		"whitelists":        whitelists,
		"tx_limits":         limits,
	}, nil
}
