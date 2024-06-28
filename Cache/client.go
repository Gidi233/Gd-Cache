package Cache

import (
	"context"
	"time"

	pb "github.com/Gidi233/Gd-Cache/CachePB"
	"github.com/Gidi233/Gd-Cache/register"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

// client 模块实现 Cache 访问其他节点获取缓存能力
type client struct {
	name string // 服务名称: Cache/ip:port
}

func (c *client) Fetch(group string, key string) ([]byte, error) {
	// 创建一个 etcd 客户端
	cli, err := clientv3.New(register.DefaultEtcdConfig)
	if err != nil {
		return nil, err
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {

		}
	}(cli)

	// 服务发现
	conn, err := register.EtcdDial(cli, c.name)
	if err != nil {
		return nil, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	grpcClient := pb.NewGroupCacheClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := grpcClient.Get(ctx, &pb.CacheRequest{Group: group, Key: key})
	if err != nil {
		return nil, err
	}

	return []byte(resp.GetValue()), nil
}
