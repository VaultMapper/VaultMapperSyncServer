package server

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/NodiumHosting/VaultMapperSyncServer/assets"
	pb "github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/NodiumHosting/VaultMapperSyncServer/util"
	"google.golang.org/protobuf/proto"
)

type Troll struct {
	name       string
	execGlobal func()
	execVault  func(vaultId string)
	execPlayer func(playerUuid string)
}

var Trolls = []Troll{
	{"bad_apple", badAppleGlobal, badAppleVault, badApplePlayer},
}

func FindTroll(trollName string) *Troll {
	for _, troll := range Trolls {
		if trollName == troll.name {
			return &troll
		}
	}

	return nil
}

func badAppleGlobal() {
	go HUB.Vaults.Range(func(key, value any) bool {
		vaultId := key.(string)
		vault := value.(*Vault)

		// TODO: remove these logs
		fmt.Println("vault " + vaultId)

		go vault.Connections.Range(func(key, value any) bool {
			playerUuid := key.(string)
			connection := value.(*Connection)

			// TODO: remove these logs
			fmt.Println("vault " + vaultId + " player " + playerUuid)

			go badApple(connection)

			return true
		})

		return true
	})
}

func badAppleVault(vaultId string) {
	// TODO: implement
}

func badApplePlayer(playerUuid string) {
	// TODO: implement
}

// TODO: check if player still connected, else we get a crash for sending to a closed chan

func badApple(connection *Connection) {
	width := assets.BadApple.Width
	height := assets.BadApple.Height

	util.RunAtFramerate(assets.BadApple.FPS, assets.BadApple.Frames, func(frame int) {
		frameData := assets.BadApple.GetFrameData(frame)
		pixelIndex := 0

		frameData.Iterate(func(bits byte) {
			for i := 0; i < 8; i++ {
				pixel := (bits >> (7 - i)) & 1

				x := (pixelIndex%width)*2 + 1
				z := (pixelIndex/width)*2 + 1

				var rt pb.RoomType

				switch pixel {
				case 0:
					rt = pb.RoomType_ROOMTYPE_BASIC
				case 1:
					rt = pb.RoomType_ROOMTYPE_ORE
				default:
					rt = pb.RoomType_ROOMTYPE_BASIC
				}

				rn := ""

				if pixel == 1 {
					rn = "the_vault:vault/rooms/common/ore1"
				}

				message := pb.Message{Type: pb.MessageType_VAULT_CELL, VaultCell: &pb.VaultCell{
					X:         int32(x - width),
					Z:         int32(z - height),
					CellType:  pb.CellType_CELLTYPE_ROOM,
					RoomType:  rt,
					RoomName:  rn,
					Explored:  true,
					Inscribed: false,
					Marked:    false,
				}}
				messageBuffer, err := proto.Marshal(&message)
				if err != nil {
					log.Println("error marshaling message for bad apple")
					continue
				}

				connection.Send <- messageBuffer

				pixelIndex++
			}
		})

		for x := range width * 2 {
			for z := range height * 2 {
				if x%2 == 0 || z%2 == 0 {
					message := pb.Message{Type: pb.MessageType_VAULT_CELL, VaultCell: &pb.VaultCell{
						X:         int32(x - width),
						Z:         int32(z - height),
						CellType:  pb.CellType_CELLTYPE_UNKNOWN,
						RoomType:  pb.RoomType_ROOMTYPE_UNKNOWN,
						RoomName:  "",
						Explored:  false,
						Inscribed: false,
						Marked:    false,
					}}
					messageBuffer, err := proto.Marshal(&message)
					if err != nil {
						log.Println("error marshaling message for bad apple")
						continue
					}

					connection.Send <- messageBuffer
					continue
				}
			}
		}
	})
}
