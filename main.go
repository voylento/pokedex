package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"github.com/voylento/pokedexcli/internal/pokemonapi"
)

type config struct {
	next 			string
	prev			string
}

type cliCommand struct {
	name        string
	description string
	config	    *config
	callback    func() error
}

var commandMap map[string]cliCommand

func init() {
	commandMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			config: 		 nil,
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			config: 		 nil,
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays location names in Pokemon world, 20 at a time, going forward",
			config: 		 &config{
				next:		"https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
				prev:		"",
			},
			callback:    commandMapList,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays location names in Pokemon world, 20 at a time, going backward",
			config: 		 &config{},
			callback:    commandMapbList,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan() 
		userText := scanner.Text()
		cleanText := cleanInput(userText)
		if len(cleanText) == 0 {
			continue
		}
		command, ok := commandMap[cleanText[0]]
		if ok {
			err := command.callback()
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, val := range commandMap {
		fmt.Printf("%-8s %v\n", val.name+":", val.description)
	}
	return nil
}

func displayLocationAreas(cmd cliCommand, url string) error {
	locationResponse, err := pokemonapi.GetLocationAreas(url)
	if err != nil {
		return fmt.Errorf("Error fetching location areas: %w", err)
	}

	if next, ok := locationResponse.Next.(string); ok {
		cmd.config.next = next
	} else {
		cmd.config.next = ""
	}
	
	if prev, ok := locationResponse.Previous.(string); ok {
		cmd.config.prev = prev
	} else {
		cmd.config.prev = ""
	}

	for _, item := range locationResponse.Results {
		fmt.Printf("%s\n", item.Name)
	}

	return nil
}

func commandMapList() error {
	cmd := commandMap["map"]

	if cmd.config == nil || cmd.config.next == "" {
		fmt.Println("No locations available")
		return nil
	}

	return displayLocationAreas(cmd, cmd.config.next)
}

func commandMapbList() error {
	cmd := commandMap["map"]

	if cmd.config == nil || cmd.config.prev == "" {
		fmt.Println("No previous locations available")
		return nil
	}

	return displayLocationAreas(cmd, cmd.config.prev)
}

func cleanInput(text string) []string {
	textLowered := strings.ToLower(text)
	words := strings.Fields(textLowered)
	return words
}
