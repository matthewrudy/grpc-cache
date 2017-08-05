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
		cache: make(map[string][]byte),
	}
}

type cacheService struct {
	sync.RWMutex
	cache map[string][]byte
}

func (service *cacheService) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	key := req.GetKey()

	if key == "sleep" {
		fmt.Println("sleeping for 10 seconds")
		time.Sleep(time.Second * 10)
	}

	service.RLock()
	val, ok := service.cache[key]
	service.RUnlock()

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

	service.Lock()
	service.cache[key] = val
	service.Unlock()

	return &proto.PutResponse{}, nil
}
