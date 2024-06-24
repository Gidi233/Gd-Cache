package Cache

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	pb "github.com/Gidi233/Gd-Cache/CachePB"
	"github.com/Gidi233/Gd-Cache/consistentHash"
	"google.golang.org/grpc"
)

/*
   Server
*/

// 节点间通信前缀,例如 http://example.net/_Cache/
const (
	defaultAddr     = "127.0.0.1:8024"
	defaultReplicas = 100
)

// server 实现了一个 gRPC 服务器
type server struct {
	pb.UnimplementedGroupCacheServer

	addr       string     // format: ip:port
	status     bool       // true: running / false: stop
	mu         sync.Mutex
	consHash   *consistentHash.Map
	clients    map[string]*client
}

func NewServer(addr string) (*server, error) {
	if addr == "" {
		addr = defaultAddr
	}
	if !validPeerAddr(addr) {
		return nil, fmt.Errorf("invalid peer address: %s", addr)
	}
	return &server{addr: addr}, nil
}

func validPeerAddr(addr string) bool {
	token1 := strings.Split(addr, ":")
	if len(token1) != 2 {
		return false
	}
	token2 := strings.Split(token1[0], ".")
	if token1[0] != "localhost" && len(token2) != 4 {
		return false
	}
	return true
}

func (s *server) Get(ctx context.Context, in *pb.CacheRequest) (*pb.CacheResponse, error) {
	group, key := in.GetGroup(), in.GetKey()
	resp := &pb.CacheResponse{}

	log.Printf("[Cache_server %s] recv RPC request - (%s)/(%s)", s.addr, group, key)
	if key == "" {
		return resp, fmt.Errorf("key is empty")
	}
	g := GetGroup(group)
	if g == nil {
		return resp, fmt.Errorf("group %s not found", group)
	}

	view, err := g.Get(key)
	if err != nil {
		return resp, err
	}
	resp.Value = string(view.ByteSlice())
	return resp, nil
}

// Start 启动 Cache 服务
func (s *server) Start() error {
	s.mu.Lock()
	if s.status == true { // running
		s.mu.Unlock()
		return fmt.Errorf("Cache server %s already running", s.addr)
	}
	s.status = true

	port := strings.Split(s.addr, ":")[1]
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGroupCacheServer(grpcServer, s)

	s.mu.Unlock()

	if err := grpcServer.Serve(lis); s.status && err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// SetPeers 将各个远端主机 IP 加入到 Server 中
func (s *server) SetPeers(peersAddr ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.consHash = consistentHash.New(defaultReplicas, nil)
	s.consHash.Add(peersAddr...)
	s.clients = make(map[string]*client)
	for _, addr := range peersAddr {
		if !validPeerAddr(addr) {
			panic(fmt.Errorf("invalid peer address: %s", addr))
		}
		service := fmt.Sprintf("Cache/%s", addr)
		s.clients[addr] = &client{name: service}
	}
}

// Pick 根据一致性哈希选择节点
func (s *server) Pick(key string) (Fetcher, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	peerAddr := s.consHash.Get(key)
	// Pick itself
	if peerAddr == s.addr {
		log.Printf("[Cache_server %s] Pick itself", s.addr)
		return nil, false
	}
	log.Printf("[Cache_server %s] Pick peer %s", s.addr, peerAddr)
	return s.clients[peerAddr], true
}

// Stop 停止 Cache 服务
func (s *server) Stop() {
	s.mu.Lock()
	if s.status == false {
		s.mu.Unlock()
		return
	}

	s.status = false
	s.clients = nil
	s.consHash = nil
	s.mu.Unlock()
}
