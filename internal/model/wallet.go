package model

type Wallet struct {
	Id         uint64 `json:"id"           description:"钱包ID"`
	UserId     uint64 `json:"userId"       description:"用户ID"`
	Address    string `json:"address"      description:"钱包地址"`
	Mnemonic   string `json:"mnemonic"     description:"助记词(加密存储)"`
	PrivateKey string `json:"privateKey"   description:"私钥(加密存储)"`
	PublicKey  string `json:"publicKey"    description:"公钥"`
	Chain      string `json:"chain"        description:"链名称 ETH/BSC"`
	WalletType int    `json:"walletType"   description:"钱包类型 1:助记词 2:私钥导入"`
	CreatedAt  int64  `json:"createdAt"    description:"创建时间"`
	UpdatedAt  int64  `json:"updatedAt"    description:"更新时间"`
}
