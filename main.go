package main

import (
	"cmp"
	"fmt"
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
	parseEnv()
	fmt.Println(ipAddress, port)

	VMServer.InitDB()

	VMServer.RunTerminal()

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
