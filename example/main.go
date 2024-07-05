package main

import (
	"fmt"
	"log"
	"sync"

	cache "github.com/Gidi233/Gd-Cache"
)

func main() {
	// 模拟MySQL数据库
	var mysql = map[string]string{
		"test":     "666",
		"realtest": "777",
		"faketest": "888",
	}
	// 新建cache实例
	group := cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[MySQL] search key", key)
			if v, ok := mysql[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	// 启动一个服务实例
	var addr string = "localhost:8088"
	svr, err := cache.NewServer(addr)
	if err != nil {
		log.Fatal(err)
	}
	// 设置同伴节点 IP (包括自己)
	svr.SetPeers(addr)
	// 将服务与 cache 绑定 因为 cache 和 server 是解耦合的
	group.RegisterPeers(svr)
	log.Println("cache is running at", addr)

	// 启动服务(注册服务至 etcd / 计算一致性哈希)
	go func() {
		// Start将不会return 除非服务stop或者抛出error
		err = svr.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 发出几个Get请求
	var wg sync.WaitGroup
	wg.Add(4)
	go GetScore(group, &wg)
	go GetScore(group, &wg)
	go GetScore(group, &wg)
	go GetScore(group, &wg)
	wg.Wait()

	for i := 0; i < 1000; i++ {
		GetRealScore(group)
	}
}

func GetScore(group *cache.Group, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("get test...")
	view, err := group.Get("test")
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println(view.String())
}

func GetRealScore(group *cache.Group) {
	view, err := group.Get("realtest")
	log.Printf("get realtest...")
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println(view.String())
}
