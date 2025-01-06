package ethclientx

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"go-wallet-defi/internal/model"
	"sync"
)

var (
	// 以太坊
	client  *ethclient.Client
	clients = make(map[uint64]*ethclient.Client)
	once    sync.Once
	mutex   sync.RWMutex
)

func GetClient(ctx context.Context) *ethclient.Client {
	once.Do(func() {
		var err error
		rpcURL := g.Cfg().MustGet(ctx, "ethereum.rpc").String()
		client, err = ethclient.Dial(rpcURL)
		if err != nil {
			panic(err)
		}
	})
	return client
}

func GetClientByChainId(ctx context.Context, chainId uint64) (*ethclient.Client, error) {
	mutex.RLock()
	if client, ok := clients[chainId]; ok {
		mutex.RUnlock()
		return client, nil
	}
	mutex.RUnlock()
	//1.获取链信息
	var chain *model.Chain
	err := g.DB().Model("chain").
		Where("chain_id", chainId).
		Where("status", 1).
		Scan(&chain)
	if err != nil {
		return nil, err
	}
	//2.解析RPC地址
	var rpcUrls []string
	err = json.Unmarshal([]byte(chain.RpcUrls), &rpcUrls)
	if err != nil {
		return nil, err
	}
	//3.连接RPC节点
	var client *ethclient.Client
	for _, url := range rpcUrls {
		client, err = ethclient.Dial(url)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	mutex.Lock()
	clients[chainId] = client
	mutex.Unlock()

	return client, nil
}

// CloseAll 关闭所有客户端连接
func CloseAll() {
	mutex.Lock()
	defer mutex.Unlock()
	client.Close()
	for _, client := range clients {
		client.Close()
	}
	clients = make(map[uint64]*ethclient.Client)
}
