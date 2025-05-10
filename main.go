package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"math"
	"math/rand/v2"
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
	callback    func(args []string) error
}

var commandMap map[string]cliCommand
var pokedex map[string]pokemonapi.Pokemon

const locationAreaEndpoint = "https://pokeapi.co/api/v2/location-area/"
const pokemonEndpoint = "https://pokeapi.co/api/v2/pokemon/"

func initializeApp() {
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
				next:		   "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
				prev:		   "",
			},
			callback:    commandMapList,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays location names in Pokemon world, 20 at a time, going backward",
			config: 		 &config{},
			callback:    commandMapbList,
		},
		"explore": {
			name:        "explore",
			description: "Displays all the pokemon in an area",
			config: 		 &config{},
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attemtps to catch a pokemon",
			config: 		 &config{},
			callback:    commandCatchPokemon,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect details of a captured pokemon",
			config: 		 &config{},
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all the pokemon captured in the pokedex",
			config: 		 &config{},
			callback:    commandListPokedex,
		},
	}

	pokedex = map[string]pokemonapi.Pokemon{}
}

func main() {
	initializeApp()
	pokemonapi.InitializePackage()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan() 
		userText := scanner.Text()
		cleanText := cleanInput(userText)
		command, ok := commandMap[cleanText[0]]
		if ok {
			err := command.callback(cleanText)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(_ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ []string) error {
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

func commandMapList(_ []string) error {
	cmd := commandMap["map"]

	if cmd.config == nil || cmd.config.next == "" {
		fmt.Println("No locations available")
		return nil
	}

	return displayLocationAreas(cmd, cmd.config.next)
}

func commandMapbList(_ []string) error {
	cmd := commandMap["map"]

	if cmd.config == nil || cmd.config.prev == "" {
		fmt.Println("No previous locations available")
		return nil
	}

	return displayLocationAreas(cmd, cmd.config.prev)
}

func getPokemonNames(resp *pokemonapi.ExploreResponse) []string {
	var pokemonNames []string

	for _, encounter := range resp.PokemonEncounters {
		pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
	}

	return pokemonNames
}

func commandExplore(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("explore command requires area name as argument")
	}

	fmt.Printf("Exploring %s...\n", args[1])

	exploreUrl := locationAreaEndpoint + args[1]
	
	exploreResponse, err := pokemonapi.GetExploreArea(exploreUrl)
	if err != nil {
		return fmt.Errorf("Error fetching area information: %w", err)
	}

	pokemonNames := getPokemonNames(exploreResponse)
	
	if len(pokemonNames) > 0 {
		fmt.Println("Found pokemon:")
		for _, name := range pokemonNames {
			fmt.Printf("- %s\n", name)
		}
	} else {
		fmt.Println("No pokemon found")
	}

	return nil
}

func calculateCaptureChance(experience int) float64 {
	baseChance := .9
	decayFactor := .01
	successChance := baseChance * math.Exp(-decayFactor*float64(experience))
	minChance := 0.1
	
	if successChance < minChance {
		successChance = minChance
	}

	return successChance
}

func attemptCapture(experience int) bool {
	successChance := calculateCaptureChance(experience)
	randomValue := rand.Float64()

	return randomValue < successChance
}

func commandCatchPokemon(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("catch command requires pokemon name as argument")
	}

	pokemonUrl := pokemonEndpoint + args[1]

	fmt.Printf("Throwing a Pokeball at %s...\n", args[1])

	pokemon, err := pokemonapi.GetPokemon(pokemonUrl)
	if err != nil {
		return fmt.Errorf("Error fetching pokemon: %w", err)
	}

	captured := attemptCapture(pokemon.BaseExperience)
	if captured {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		pokedex[pokemon.Name] = *pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandListPokedex(_ []string) error {
	fmt.Println("Your Pokedex:")
	for _, val := range pokedex {
		fmt.Printf("- %s\n", val.Name)
	}

	return nil
}

func commandInspect(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("inspect command requires pokemon name as argument")
	}

	if pokemon, ok := pokedex[args[1]]; !ok {
		fmt.Println("you have not caught that pokemoh")
	} else {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, pokeType := range pokemon.Types {
			fmt.Printf("  -%s\n", pokeType.Type.Name)
		}
	}

	return nil
}

func cleanInput(text string) []string {
	textLowered := strings.ToLower(text)
	words := strings.Fields(textLowered)
	return words
}
