Dev:
- go mod tidy
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
- ... idk

Updating proto:
- `protoc --go_out=. --go-grpc_out=. ./VaultMapperProtocol/vaultmapper.proto`