package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/herodragmon/pokedex/internal/pokecache"
)

func startRepl(cfg *config) {
  reader := bufio.NewScanner(os.Stdin)
  for {
    fmt.Print("Pokedex > ")
    if !reader.Scan() { return }
    words := cleanInput(reader.Text())
    if len(words) == 0 { continue }

    name, args := words[0], words[1:]
    cmd, ok := getCommands()[name]
    if !ok {
      fmt.Println("Unknown command")
      continue
    }
    if err := cmd.callback(cfg, args); err != nil {
      fmt.Println(err)
    }
  }
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

type config struct {
    Next     *string
    Previous *string
		Cache    pokecache.Cache
}


type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}


func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays names of location",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Goes to the Previous location",
			callback:    commandMapB,
		},
	}
}

