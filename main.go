package main

import (
	"cmp"
	"fmt"
	pb "github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/NodiumHosting/VaultMapperSyncServer/render"
	VMServer "github.com/NodiumHosting/VaultMapperSyncServer/server"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"strconv"
)

var (
	ipAddress string = "0.0.0.0"
	port      int    = 42069
)

func main() {
	cells := make([]*pb.VaultCell, 0)
	for i := -24; i <= 24; i++ {
		for j := -24; j <= 24; j++ {
			rt := pb.RoomType_ROOMTYPE_BASIC
			if i == 0 && j == 0 {
				rt = pb.RoomType_ROOMTYPE_START
			}
			ct := pb.CellType_CELLTYPE_UNKNOWN
			//only if both are even - room
			if i%2 == 0 && j%2 == 0 {
				ct = pb.CellType_CELLTYPE_ROOM
			}
			// if one is even - tunnel
			if i%2 == 0 && j%2 != 0 {
				ct = pb.CellType_CELLTYPE_TUNNEL_Z
			} else if i%2 != 0 && j%2 == 0 {
				ct = pb.CellType_CELLTYPE_TUNNEL_X
			}
			// if both are odd - void
			if i%2 != 0 && j%2 != 0 {
				continue
			}

			rn := pb.RoomName_ROOMNAME_UNKNOWN
			if ct == pb.CellType_CELLTYPE_ROOM {
				rnd := rand.Int() % 100
				aOrB := rand.Int() % 2
				if rnd < 15 {
					if aOrB == 0 {
						rt = pb.RoomType_ROOMTYPE_OMEGA
					} else {
						rt = pb.RoomType_ROOMTYPE_CHALLENGE
					}
					rname := rand.Int()%17 + 1
					rn = pb.RoomName(rname)
				}
			}
			cells = append(cells, &pb.VaultCell{
				X:         int32(i),
				Z:         int32(j),
				CellType:  ct,
				RoomType:  rt,
				RoomName:  rn,
				Explored:  true,
				Inscribed: false,
				Marked:    false,
			})
		}
	}
	err, data := render.RenderVault(cells)
	if err != nil {
		return
	}
	//save to out.png
	f, err := os.Create("out.png")
	if err != nil {
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	if _, err = f.Write(data); err != nil {
		return
	}

	parseEnv()
	fmt.Println(ipAddress, port)

	VMServer.Run(ipAddress, port)
}

// parseEnv() parses environment variables and reverts to defaults if necessary
func parseEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env, using defaults")
	} else {
		ipAddress = cmp.Or(os.Getenv("IP_ADDRESS"), ipAddress) // default to 127.0.0.1 if not set in env file

		parseInt, err := strconv.Atoi(os.Getenv("PORT"))
		if err == nil { // if the Atoi works, the port exists
			if parseInt > 0 && parseInt < 65536 { // if the port is valid
				port = parseInt
			}
		}
	}
}
