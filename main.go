package main

import (
	"fmt"
	"github.com/NodiumHosting/VaultMapperSyncServer/proto"
)

func main() {
	vaultCell := proto.VaultCell{
		X:        0,
		z:        0,
		celltype: proto.CellType_CELLTYPE_ROOM,
	}

	fmt.Println("Hello, World!")
}
