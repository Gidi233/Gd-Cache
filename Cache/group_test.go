package Cache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"test":     "666",
	"faketest": "777",
	"realtest": "888",
}

func TestCache_Getter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("test")
	if v, _ := f.Get("test"); string(v) != string(expect) {
		t.Fatalf("callback failed")
	}
}

func TestGetterFunc_Get(t *testing.T) {
	loadCounts := make(map[string]int, len(db))

	getter := GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	})

	gee := NewGroup("CacheTest", 2<<10, getter)

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", k)
		} // load from callback function
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
