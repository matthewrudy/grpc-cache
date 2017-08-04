proto:
	protoc --go_out=plugins=grpc:. cache/proto/cache.proto

server:
	go run cache/cmd/server.go

client:
	go run main.go
