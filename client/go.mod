module grpc-redis/client

go 1.22.3

require (
	github.com/go-redis/redis/v8 v8.11.5
	google.golang.org/grpc v1.64.0
	grpc-redis-client/protos/todo v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace grpc-redis-client/protos/todo => ../protos/todo
