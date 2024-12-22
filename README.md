Prerequisites:
- Go 1.22.3+ (previous versions might be fine, but this is verified working)
- Protobuf Gen `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

Dev:
- `go mod tidy`
- `go run .` (probably just these)

Updating proto:
- `protoc --go_out=. ./VaultMapperProtocol/vaultmapper.proto`