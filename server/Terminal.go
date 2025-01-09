package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func HandleCommand(command string) {
	if strings.HasPrefix(command, "toast ") {
		restOfCommand := strings.TrimPrefix(command, "toast ")
		log.Println("Broadcasting toast: " + restOfCommand)
		HUB.BroadcastToast(restOfCommand)
	}

	if strings.HasPrefix(command, "toastv ") {
		restOfCommand := strings.TrimPrefix(command, "toastv ")
		split := strings.Split(restOfCommand, " ")
		vaultID := split[0]
		text := strings.Join(split[1:], " ")
		log.Println("Broadcasting toast in vault " + vaultID + ": " + text)
		HUB.BroadcastToastInVault(vaultID, text)
	}
}

func StartTerminal() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Terminal started")
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		HandleCommand(strings.TrimSpace(text))
	}
}

func RunTerminal() {
	go StartTerminal()
}
