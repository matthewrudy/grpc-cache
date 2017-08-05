package cache

import (
	"fmt"
	"sync"
	"time"

	proto "github.com/matthewrudy/grpc-cache/cache/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewService() proto.CacheServer {
	return &cacheService{
		cacheStore: cacheStore{
			store: make(map[string][]byte),
		},
	}
}

type cacheStore struct {
	sync.RWMutex
	store map[string][]byte
}

func (c *cacheStore) get(key string) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok := c.store[key]
	return val, ok
}

func (c *cacheStore) set(key string, val []byte) {
	c.Lock()
	defer c.Unlock()

	c.store[key] = val
}

type cacheService struct {
	cacheStore
}

func (service *cacheService) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	key := req.GetKey()

	if key == "sleep" {
		fmt.Println("sleeping for 10 seconds")
		time.Sleep(time.Second * 10)
	}

	val, ok := service.cacheStore.get(key)

	if !ok {
		return nil, status.Errorf(codes.NotFound, "key not found %s", key)
	}
	fmt.Printf("get key=%s val=%s\n", key, val)
	return &proto.GetResponse{
		Key: key,
		Val: val,
	}, nil
}

func (service *cacheService) Put(ctx context.Context, req *proto.PutRequest) (*proto.PutResponse, error) {
	key := req.GetKey()
	val := req.GetVal()
	fmt.Printf("set key=%s val=%s\n", key, val)

	service.cacheStore.set(key, val)

	return &proto.PutResponse{}, nil
}
