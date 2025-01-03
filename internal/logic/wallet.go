package logic

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/pkg/cryptox"
	"go-wallet-defi/internal/service"
	"time"
)

// 确保WalletLogic实现了IWallet接口
var _ service.IWallet = (*WalletLogic)(nil)

type WalletLogic struct{}

func (s *WalletLogic) GetBalance(ctx context.Context, address string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *WalletLogic) Create(ctx context.Context, chain string) (*model.Wallet, error) {
	// 生成助记词
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	// 从助记词生成种子
	//seed := bip39.NewSeed(mnemonic, "")

	// 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	// 获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}

	// 获取地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 加密敏感信息
	encryptedMnemonic, err := cryptox.Encrypt(mnemonic)
	if err != nil {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	encryptedPrivateKey, err := cryptox.Encrypt(hexutil.Encode(privateKeyBytes))
	if err != nil {
		return nil, err
	}

	wallet := &model.Wallet{
		Address:    address.Hex(),
		Mnemonic:   encryptedMnemonic,
		PrivateKey: encryptedPrivateKey,
		PublicKey:  hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA)),
		Chain:      chain,
		WalletType: 1,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	// 保存到数据库
	if err := dao.Wallet.Insert(ctx, wallet); err != nil {
		return nil, err
	}

	// 返回时解密助记词(只在创建时返回明文助记词)
	wallet.Mnemonic = mnemonic

	return wallet, nil
}

// Import 导入钱包
func (s *WalletLogic) Import(ctx context.Context, chain, mnemonic string) (*model.Wallet, error) {
	// 验证助记词
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic")
	}

	// 从助记词生成种子
	seed := bip39.NewSeed(mnemonic, "")

	// 生成私钥
	privateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return nil, err
	}

	// 获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}

	// 获取地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 检查地址是否已存在
	exists, err := dao.Wallet.CheckAddressExists(ctx, address.Hex())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("wallet already exists")
	}

	// 加密存储
	encryptedMnemonic, err := cryptox.Encrypt(mnemonic)
	if err != nil {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	encryptedPrivateKey, err := cryptox.Encrypt(hexutil.Encode(privateKeyBytes))
	if err != nil {
		return nil, err
	}

	wallet := &model.Wallet{
		Address:    address.Hex(),
		Mnemonic:   encryptedMnemonic,
		PrivateKey: encryptedPrivateKey,
		PublicKey:  hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA)),
		Chain:      chain,
		WalletType: 2, // 私钥导入
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	// 保存到数据库
	if err := dao.Wallet.Insert(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetList 获取钱包列表
func (s *WalletLogic) GetList(ctx context.Context, page, pageSize int) ([]*model.Wallet, error) {
	wallets, err := dao.Wallet.GetList(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 查询每个钱包的余额
	for _, wallet := range wallets {
		balance, err := s.GetBalance(ctx, wallet.Address)
		if err != nil {
			wallet.Balance = "0"
			continue
		}
		wallet.Balance = balance
	}

	return wallets, nil
}

func init() {
	service.RegisterWallet(New())
}

func New() *WalletLogic {
	return &WalletLogic{}
}
