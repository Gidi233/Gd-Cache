package Cache

import (
	"context"
	"time"

	pb "github.com/Gidi233/Gd-Cache/CachePB"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// client 模块实现 Cache 访问其他节点获取缓存能力
type client struct {
	name string // 服务名称: Cache/ip:port
}

func (c *client) Fetch(group string, key string) ([]byte, error) {	
	conn, err := grpc.Dial(c.name, grpc.WithTransportCredentials(insecure.NewCredentials()))

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
