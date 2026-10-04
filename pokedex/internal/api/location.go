package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yeikol-l/bootdev/pokedex/internal/pokecache"
)

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationAreaDetail struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	GameIndex         int    `json:"game_index"`
	PokemonEncounters []struct {
		Pokemon Pokemon `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetLocationAreas(cache *pokecache.Cache, cursor int) ([]LocationArea, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", cursor)
	cacheHit, ok := cache.Get(url)
	if ok {
		var locationAreas []LocationArea
		err := json.Unmarshal(cacheHit, &locationAreas)
		if err != nil {
			return nil, err
		}
		return locationAreas, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("La request fallo con la url: %s", url)
	}
	var locationAreas struct {
		Count    int            `json:"count"`
		Next     string         `json:"next"`
		Previous string         `json:"previous"`
		Results  []LocationArea `json:"results"`
	}

	err = json.NewDecoder(res.Body).Decode(&locationAreas)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(locationAreas.Results)
	if err != nil {
		return nil, err
	}
	cache.Add(url, bytes)
	return locationAreas.Results, nil
}

func GetLocationAreaPokemons(cache *pokecache.Cache, areaName string) ([]Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)
	cacheHit, ok := cache.Get(url)
	if ok {
		var pokemons []Pokemon
		err := json.Unmarshal(cacheHit, &pokemons)
		return pokemons, err
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var locationAreaDetail LocationAreaDetail
	err = json.NewDecoder(res.Body).Decode(&locationAreaDetail)
	if err != nil {
		return nil, err
	}
	var pokemons []Pokemon
	for _, v := range locationAreaDetail.PokemonEncounters {
		pokemons = append(pokemons, v.Pokemon)
	}

	bytes, err := json.Marshal(pokemons)
	if err != nil {
		return nil, err
	}
	cache.Add(url, bytes)

	return pokemons, nil
}
