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

func (w WalletLogic) Import(ctx context.Context, chain, mnemonic string) (*model.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (w WalletLogic) GetBalance(ctx context.Context, address string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	service.RegisterWallet(New())
}

func New() *WalletLogic {
	return &WalletLogic{}
}
