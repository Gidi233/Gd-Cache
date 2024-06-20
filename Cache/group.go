package Cache

import (
	"fmt"
	"log"
	"sync"
)

// Getter 是一个加载指定 key 的数据的接口
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 是一个通过函数实现 Getter 接口的类型
type GetterFunc func(key string) ([]byte, error)

// Get 实现了 Getter 接口的函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 是最核心的数据结构，负责与外部交互，控制缓存存储和获取的主流程
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 创建一个新的 Group 实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g

	return g
}

// GetGroup 返回指定名称的 Group
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	g := groups[name]

	return g
}

// Get 从缓存中获取指定 key 的数据
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 从 mainCache 中查找缓存，如果存在则返回缓存值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[Cache] hit")
		return v, nil
	}

	// 如果缓存不存在，则调用 load 方法加载
	return g.load(key)
}

// func (g *Group) Delete(key, value string) {
// 	if key == "" {
// 		log.Println("[Cache] key is required")
// 		return
// 	}

// 	// Delete the corresponding entry from the cache
// 	g.deleteCache(key)
// }

// func (g *Group) deleteCache(key string) {
// 	g.mainCache.delete(key)
// }

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	// 将数据添加到缓存中
	g.populateCache(key, value)

	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
