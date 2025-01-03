package ethclient

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
	"sync"
)

var (
	// 以太坊
	client *ethclient.Client
	once   sync.Once
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
