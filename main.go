package main

import (
	"context"
	"fmt"
	"os"
	"time"

	proto "github.com/matthewrudy/grpc-cache/cache/proto"
	"google.golang.org/grpc"
)

func main() {
	if err := runClient(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
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

	fmt.Printf("get key=%s val=%s\n", resp.GetKey(), resp.GetVal())
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

	fmt.Printf("put key=%s val=%s\n", key, val)
	return nil
}

func runClient() error {
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	cache := proto.NewCacheClient(conn)

	get(cache, "foo")
	put(cache, "foo", "bar")
	get(cache, "foo")
	get(cache, "sleep")

	return nil
}
