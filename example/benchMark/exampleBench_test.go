package benchMark

import (
	"fmt"
	"testing"

	cache "github.com/Gidi233/Gd-Cache"
)

func BenchmarkGetScore(b *testing.B) {
	// 模拟MySQL数据库
	mysql := map[string]string{
		"test":     "666",
		"realtest": "777",
		"faketest": "888",
	}
	group := cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := mysql[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	_, err := group.Get("test")
	if err != nil {
		b.Fatalf("Error getting value: %s", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := group.Get("test")
		if err != nil {
			b.Fatalf("Error getting value: %s", err)
		}
	}
}
