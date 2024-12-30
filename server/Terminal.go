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
