package main

import (
	"cmp"
	"fmt"
	"github.com/NodiumHosting/VaultMapperSyncServer/icons"
	VMServer "github.com/NodiumHosting/VaultMapperSyncServer/server"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	ipAddress string = "0.0.0.0"
	port      int    = 42069
)

func main() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(multiWriter)

	parseEnv()
	log.Println(ipAddress, port)

	icons.Init()

	VMServer.InitDB()

	VMServer.RunTerminal()

	VMServer.CleanDB()   // clean the database on startup
	VMServer.StartCron() // start the cron job to clean the database every day at midnight

	VMServer.Run(ipAddress, port)

	VMServer.StopCron()
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
