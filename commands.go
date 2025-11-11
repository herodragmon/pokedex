package main

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"io"
	"math/rand"
)

type locationAreaList struct {
    Results  []struct {
        Name string `json:"name"`
    } `json:"results"`
    Next     *string `json:"next"`
    Previous *string `json:"previous"`
}

func commandExit(_ *config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config, _ []string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandMap(cfg *config, _ []string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}

	cachedData, found := cfg.Cache.Get(url)
	
	if found {
		fmt.Println("Loading location areas from cache...")
		
		var lst locationAreaList
		if err := json.Unmarshal(cachedData, &lst); err != nil {
			return fmt.Errorf("failed to unmarshal cached data: %w", err)
		}

		cfg.Next = lst.Next
		cfg.Previous = lst.Previous

		for _, l := range lst.Results {
			fmt.Println(l.Name)
		}
		
		return nil
	}

	fmt.Println("Loading location areas from API...")
	
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	cfg.Cache.Add(url, bodyBytes)
	
	var lst locationAreaList
	if err = json.Unmarshal(bodyBytes, &lst); err != nil {
		return fmt.Errorf("failed to unmarshal API data: %w", err)
	}

	cfg.Next = lst.Next
	cfg.Previous = lst.Previous

	for _, l := range lst.Results {
		fmt.Println(l.Name)
	}
	
	return nil
}

func commandMapB(cfg *config, _ []string) error {
	if cfg.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	url := *cfg.Previous

	cachedData, found := cfg.Cache.Get(url)

	if found {
		fmt.Println("Loading location areas from cache...")

		var lst locationAreaList
		if err := json.Unmarshal(cachedData, &lst); err != nil {
			return fmt.Errorf("failed to unmarshal cached data: %w", err)
		}

		cfg.Next = lst.Next
		cfg.Previous = lst.Previous

		for _, l := range lst.Results {
			fmt.Println(l.Name)
		}
		
		return nil
	}

	fmt.Println("Loading location areas from API...")

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	cfg.Cache.Add(url, bodyBytes)

	var lst locationAreaList
	if err = json.Unmarshal(bodyBytes, &lst); err != nil {
		return fmt.Errorf("failed to unmarshal API data: %w", err)
	}

	cfg.Next = lst.Next
	cfg.Previous = lst.Previous

	for _, l := range lst.Results {
		fmt.Println(l.Name)
	}

	return nil
}

type locationAreaDetail struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}


func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a location area name (e.g., 'explore eterna-forest')")
	}

	locationName := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + locationName

	fmt.Printf("Exploring %s...\n", locationName)

	cachedData, found := cfg.Cache.Get(url)
	var locationDetail locationAreaDetail
	
	if found {
		fmt.Println("Loading from cache...")
		if err := json.Unmarshal(cachedData, &locationDetail); err != nil {
			return fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
	} else {
		fmt.Println("Loading from API...")
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("location area '%s' not found", locationName)
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status: %s", res.Status)
		}

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		cfg.Cache.Add(url, bodyBytes)

		if err = json.Unmarshal(bodyBytes, &locationDetail); err != nil {
			return fmt.Errorf("failed to unmarshal API data: %w", err)
		}
	}

	fmt.Println("Found Pokémon:")
	if len(locationDetail.PokemonEncounters) == 0 {
		fmt.Println("  (No Pokémon found in this area)")
		return nil
	}

	for _, encounter := range locationDetail.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name (e.g., 'catch pikachu')")
	}

	pokemonName := args[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	if _, ok := cfg.Pokedex[pokemonName]; ok {
		return fmt.Errorf("you have already caught %s", pokemonName)
	}
	
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)


	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("pokemon '%s' not found", pokemonName)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var pokemon Pokemon
	if err = json.Unmarshal(bodyBytes, &pokemon); err != nil {
		return fmt.Errorf("failed to unmarshal API data: %w", err)
	}
	
	const catchThreshold = 50
	catchValue := rand.Intn(pokemon.BaseExperience)

	if catchValue < catchThreshold {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command")
		
		cfg.Pokedex[pokemon.Name] = pokemon
		
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name to inspect")
	}

	pokemonName := args[0]

	pokemon, found := cfg.Pokedex[pokemonName]
	if !found {
		fmt.Println("pokemon not found")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *config, _ []string) error {
	fmt.Println("Your Pokedex:")

	if len(cfg.Pokedex) == 0 {
		fmt.Println("Yo've not caught any pokemon")
		return nil
	}

	for p := range cfg.Pokedex {
		fmt.Printf(" - %s\n", p)
	}
	return nil
}
