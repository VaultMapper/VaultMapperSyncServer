package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Command struct {
	name string
	// aliases can be nil
	aliases []string
	// usage can be nil but should explain how the command is intended to be used and what it's syntax is
	usage string
	// we first try to execute execStr or if that is nil then we try to execute execArr
	execStr func(args string) (ok bool)
	execArr func(args []string) (ok bool)
}

func (c *Command) exec(args []string) (ok bool) {
	if c.execStr != nil {
		return c.execStr(strings.Join(args, " "))
	}
	if c.execArr != nil {
		return c.execArr(args)
	}
	return false
}

var commands = make([]Command, 0)

func RegisterCommand(name string, aliases []string, usage string, execStr func(args string) (ok bool), execArr func(args []string) (ok bool)) {
	commands = append(commands, Command{
		name,
		aliases,
		usage,
		execStr,
		execArr,
	})
}

func RegisterCommands() {
	RegisterCommand("help", []string{}, "help", nil, handleHelp)
	RegisterCommand("toast", nil, "toast <message>", handleToast, nil)
	RegisterCommand("toastv", nil, "toastv <vault_id> <message>", nil, handleToastV)

	RegisterCommand("troll", nil, "troll global <troll>; troll <vault/player> <vault_id/player_uuid> <troll>; troll list", nil, handleTroll)
}

func handleHelp(args []string) (ok bool) {
	if len(args) == 0 {
		fmt.Print("Available Commands: ")
		for i, cmd := range commands {
			fmt.Print(cmd.name)
			if i < len(commands)-1 {
				fmt.Print(", ")
			} else {
				fmt.Println()
			}
		}
	} else {
		command := args[0]

		for _, cmd := range commands {
			if command == cmd.name {
				log.Println("Command " + cmd.name + " usage: " + cmd.usage)
				return true
			}

			for _, alias := range cmd.aliases {
				if command == alias {
					log.Println("Command " + cmd.name + " usage: " + cmd.usage)
					return true
				}
			}
		}

		log.Println("Cannot list help for unknown command")

		return true
	}

	return true
}

func handleToast(args string) (ok bool) {
	log.Println("Broadcasting toast: " + args)
	HUB.BroadcastToast(args)

	return true
}

func handleToastV(args []string) (ok bool) {
	if len(args) < 2 {
		// not enough args
		return false
	}

	vaultID := args[0]
	text := strings.Join(args[1:], " ")
	log.Println("Broadcasting toast in vault " + vaultID + ": " + text)
	HUB.BroadcastToastInVault(vaultID, text)

	return true
}

func handleTroll(args []string) (ok bool) {
	// troll global <troll>
	// troll <vault/player> <vault_id/player_uuid> <troll>
	// troll list

	if len(args) < 1 {
		return false
	}

	switch args[0] {
	case "list":
		fmt.Print("Available Trolls: ")
		for i, troll := range Trolls {
			fmt.Print(troll.name)
			if i < len(Trolls)-1 {
				fmt.Print(", ")
			} else {
				fmt.Println()
			}
		}
	case "global":
		if len(args) < 2 {
			return false
		}
		troll := FindTroll(args[1])
		if troll == nil {
			log.Println("Unknown Troll")
			return true
		}
		troll.execGlobal()
	case "vault":
		fallthrough
	case "player":
		if len(args) < 3 {
			return false
		}

		troll := FindTroll(args[2])
		if troll == nil {
			log.Println("Unknown Troll")
			return true
		}

		if args[0] == "vault" {
			troll.execVault(args[1])
		} else {
			troll.execPlayer(args[1])
		}
	}

	return true
}

func HandleCommand(commandString string) {
	parts := strings.Split(commandString, " ")

	if len(parts) == 0 {
		log.Println("Invalid Command")
		return
	}

	command := parts[0]

	for _, cmd := range commands {
		if command == cmd.name {
			ok := cmd.exec(parts[1:])
			if !ok {
				log.Println("Failed to execute command " + cmd.name + " (" + command + "). Usage: " + cmd.usage)
			}
			return
		}

		for _, alias := range cmd.aliases {
			if command == alias {
				ok := cmd.exec(parts[1:])
				if !ok {
					log.Println("Failed to execute command " + cmd.name + " (" + command + "). Usage: " + cmd.usage)
				}
				return
			}
		}
	}

	log.Println("Unknown command")
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
