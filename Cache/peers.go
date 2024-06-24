package Cache

// Picker 定义了获取分布式节点的能力
type Picker interface {
	Pick(key string) (peer Fetcher, ok bool)
}

// Fetcher 定义了从远端获取缓存的能力, 每个 Peer 都应该实现此接口
type Fetcher interface {
	Fetch(group string, key string) ([]byte, error)
}
