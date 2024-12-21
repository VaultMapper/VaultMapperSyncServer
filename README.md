Prerequisites:
- Go 1.22.3+ (previous versions might be fine, but this is verified working)
- Protobuf Gen `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- Protobuf GRPC `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

Dev:
- `go mod tidy`
- `go run .` (probably just these)

Updating proto:
- `protoc --go_out=. --go-grpc_out=. ./VaultMapperProtocol/vaultmapper.proto`