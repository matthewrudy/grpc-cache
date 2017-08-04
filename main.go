package main

import (
	"context"
	"log"
	"time"

	proto "github.com/matthewrudy/grpc-cache/cache/proto"
	"google.golang.org/grpc"
)

func main() {
	if err := runClient(); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

func get(cache proto.CacheClient, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	resp, err := cache.Get(ctx, &proto.GetRequest{
		Key: key,
	})

	if err != nil {
		return err
	}

	log.Printf("get key=%s val=%s", resp.GetKey(), resp.GetVal())
	return nil
}

func put(cache proto.CacheClient, key string, val string) error {
	_, err := cache.Put(context.Background(), &proto.PutRequest{
		Key: key,
		Val: []byte(val),
	})

	if err != nil {
		return err
	}

	log.Printf("put key=%s val=%s", key, val)
	return nil
}

func runClient() error {
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	cache := proto.NewCacheClient(conn)

	if err = get(cache, "foo"); err != nil {
		log.Fatalf("getting foo failed: %v", err)
	}
	if err = put(cache, "foo", "bar"); err != nil {
		log.Fatalf("setting foo failed: %v", err)
	}
	if err = get(cache, "foo"); err != nil {
		log.Fatalf("getting foo failed: %v", err)
	}
	if err = get(cache, "sleep"); err != nil {
		log.Fatalf("getting sleep failed: %v", err)
	}

	return nil
}
