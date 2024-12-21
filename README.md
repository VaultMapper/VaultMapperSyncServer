Dev:
- go mod tidy
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
- ... idk

Updating proto:
- `protoc --go_out=. ./VaultMapperProtocol/vaultmapper.proto`