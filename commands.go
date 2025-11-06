package main

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
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

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status: %s", res.Status)
	}

	var lst locationAreaList
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&lst); err != nil {
		return err
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
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status: %s", res.Status)
	}

	var lst locationAreaList
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&lst); err != nil {
		return err
	}

	cfg.Next = lst.Next
	cfg.Previous = lst.Previous

	for _, l := range lst.Results {
		fmt.Println(l.Name)
	}
	return nil

}
