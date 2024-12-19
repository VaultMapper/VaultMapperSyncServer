package main

import (
	"fmt"
	"github.com/NodiumHosting/VaultMapperSyncServer/proto"
)

func main() {
	vaultCell := proto.VaultCell{
		X:         0,
		Z:         0,
		CellType:  proto.CellType_CELLTYPE_ROOM,
		RoomType:  proto.RoomType_ROOMTYPE_START,
		RoomName:  proto.RoomName_ROOMNAME_UNKNOWN,
		Explored:  true,
		Inscribed: false,
		Marked:    false,
	}

	fmt.Println("Hello, World!" + vaultCell.String())
}
