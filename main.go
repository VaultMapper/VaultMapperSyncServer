package main

import (
	"cmp"
	"fmt"
	pb "github.com/NodiumHosting/VaultMapperSyncServer/proto"
	VMServer "github.com/NodiumHosting/VaultMapperSyncServer/server"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var (
	ipAddress string = "0.0.0.0"
	port      int    = 42069
)

func main() {
	vaultCell := pb.VaultCell{
		X:         0,
		Z:         0,
		CellType:  pb.CellType_CELLTYPE_ROOM,
		RoomType:  pb.RoomType_ROOMTYPE_START,
		RoomName:  pb.RoomName_ROOMNAME_UNKNOWN,
		Explored:  true,
		Inscribed: false,
		Marked:    false,
	}

	fmt.Println("Hello, World!\n" + vaultCell.String())

	parseEnv()
	fmt.Println(ipAddress, port)

	VMServer.Serve(ipAddress, port)

	/*
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		s.Serve(lis)
		pb.RegisterVMServiceServer(s, &server)*/
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

/*
type server struct {
	pb.UnimplementedVMServiceServer
}

func (s *server) NewVaultCell(ctx context.Context, req *pb.VaultCell) (*pb.VaultCell, error) {
	// Handle the new cell
	return &pb.VaultCell{}, nil
}*/
