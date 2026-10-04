package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yeikol-l/bootdev/pokedex/internal/pokecache"
)

type PokemonDetails struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	PastAbilities []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Abilities []struct {
			IsHidden bool        `json:"is_hidden"`
			Slot     int         `json:"slot"`
			Ability  interface{} `json:"ability"`
		} `json:"abilities"`
	} `json:"past_abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order interface{} `json:"order"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		Other struct {
			Home struct {
				FrontShiny       string      `json:"front_shiny"`
				FrontFemale      interface{} `json:"front_female"`
				FrontDefault     string      `json:"front_default"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"home"`
			Showdown struct {
				BackShiny        string      `json:"back_shiny"`
				BackFemale       interface{} `json:"back_female"`
				FrontShiny       string      `json:"front_shiny"`
				BackDefault      string      `json:"back_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontDefault     string      `json:"front_default"`
				BackShinyFemale  interface{} `json:"back_shiny_female"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"showdown"`
			DreamWorld struct {
				FrontFemale  interface{} `json:"front_female"`
				FrontDefault string      `json:"front_default"`
			} `json:"dream_world"`
			OfficialArtwork struct {
				Versions struct {
					GenerationI struct {
						RedAndBlue struct {
							FrontDefault string `json:"front_default"`
						} `json:"red-and-blue"`
						RedAndGreen struct {
							FrontDefault string `json:"front_default"`
						} `json:"red-and-green"`
					} `json:"generation-i"`
					GenerationIi struct {
						GoldAndSilver struct {
							FrontDefault interface{} `json:"front_default"`
						} `json:"gold-and-silver"`
					} `json:"generation-ii"`
				} `json:"versions"`
				FrontShiny   string `json:"front_shiny"`
				FrontDefault string `json:"front_default"`
			} `json:"official-artwork"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				Yellow struct {
					BackGbc              string `json:"back_gbc"`
					BackGray             string `json:"back_gray"`
					FrontGbc             string `json:"front_gbc"`
					FrontGray            string `json:"front_gray"`
					BackDefault          string `json:"back_default"`
					FrontDefault         string `json:"front_default"`
					BackTransparent      string `json:"back_transparent"`
					FrontTransparent     string `json:"front_transparent"`
					BackTransparentGray  string `json:"back_transparent_gray"`
					FrontTransparentGray string `json:"front_transparent_gray"`
				} `json:"yellow"`
				RedBlue struct {
					BackGray             string `json:"back_gray"`
					FrontGray            string `json:"front_gray"`
					BackDefault          string `json:"back_default"`
					FrontDefault         string `json:"front_default"`
					BackTransparent      string `json:"back_transparent"`
					FrontTransparent     string `json:"front_transparent"`
					BackTransparentGray  string `json:"back_transparent_gray"`
					FrontTransparentGray string `json:"front_transparent_gray"`
				} `json:"red-blue"`
				RedGreenJapan struct {
					BackGray     string `json:"back_gray"`
					FrontGray    string `json:"front_gray"`
					BackDefault  string `json:"back_default"`
					FrontDefault string `json:"front_default"`
				} `json:"red-green-japan"`
			} `json:"generation-i"`
			GenerationV struct {
				Icons struct {
					Animated struct {
						FrontDefault string `json:"front_default"`
					} `json:"animated"`
					FrontDefault string `json:"front_default"`
				} `json:"icons"`
				BlackWhite struct {
					Animated struct {
						BackShiny        string      `json:"back_shiny"`
						BackFemale       interface{} `json:"back_female"`
						FrontShiny       string      `json:"front_shiny"`
						BackDefault      string      `json:"back_default"`
						FrontFemale      interface{} `json:"front_female"`
						FrontDefault     string      `json:"front_default"`
						BackShinyFemale  interface{} `json:"back_shiny_female"`
						FrontShinyFemale interface{} `json:"front_shiny_female"`
					} `json:"animated"`
					BackShiny        string      `json:"back_shiny"`
					BackFemale       interface{} `json:"back_female"`
					FrontShiny       string      `json:"front_shiny"`
					BackDefault      string      `json:"back_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontDefault     string      `json:"front_default"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationIi struct {
				Gold struct {
					BackShiny             string `json:"back_shiny"`
					FrontShiny            string `json:"front_shiny"`
					BackDefault           string `json:"back_default"`
					FrontDefault          string `json:"front_default"`
					BackTransparent       string `json:"back_transparent"`
					FrontTransparent      string `json:"front_transparent"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
				} `json:"gold"`
				Silver struct {
					BackShiny             string `json:"back_shiny"`
					FrontShiny            string `json:"front_shiny"`
					BackDefault           string `json:"back_default"`
					FrontDefault          string `json:"front_default"`
					BackTransparent       string `json:"back_transparent"`
					FrontTransparent      string `json:"front_transparent"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
				} `json:"silver"`
				Crystal struct {
					Animated struct {
						FrontShiny   string `json:"front_shiny"`
						FrontDefault string `json:"front_default"`
					} `json:"animated"`
					BackShiny             string `json:"back_shiny"`
					FrontShiny            string `json:"front_shiny"`
					BackDefault           string `json:"back_default"`
					FrontDefault          string `json:"front_default"`
					BackTransparent       string `json:"back_transparent"`
					FrontTransparent      string `json:"front_transparent"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
				} `json:"crystal"`
			} `json:"generation-ii"`
			GenerationIv struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
				} `json:"icons"`
				Platinum struct {
					Animated struct {
						FrontShiny       string      `json:"front_shiny"`
						FrontFemale      interface{} `json:"front_female"`
						FrontDefault     string      `json:"front_default"`
						FrontShinyFemale interface{} `json:"front_shiny_female"`
					} `json:"animated"`
					BackShiny        string      `json:"back_shiny"`
					BackFemale       interface{} `json:"back_female"`
					FrontShiny       string      `json:"front_shiny"`
					BackDefault      string      `json:"back_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontDefault     string      `json:"front_default"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"platinum"`
				DiamondPearl struct {
					Animated struct {
						FrontShiny       string      `json:"front_shiny"`
						FrontFemale      interface{} `json:"front_female"`
						FrontDefault     string      `json:"front_default"`
						FrontShinyFemale interface{} `json:"front_shiny_female"`
					} `json:"animated"`
					BackShiny        string      `json:"back_shiny"`
					BackFemale       interface{} `json:"back_female"`
					FrontShiny       string      `json:"front_shiny"`
					BackDefault      string      `json:"back_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontDefault     string      `json:"front_default"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					Animated struct {
						FrontShiny       string      `json:"front_shiny"`
						FrontFemale      interface{} `json:"front_female"`
						FrontDefault     string      `json:"front_default"`
						FrontShinyFemale interface{} `json:"front_shiny_female"`
					} `json:"animated"`
					BackShiny        string      `json:"back_shiny"`
					BackFemale       interface{} `json:"back_female"`
					FrontShiny       string      `json:"front_shiny"`
					BackDefault      string      `json:"back_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontDefault     string      `json:"front_default"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
			} `json:"generation-iv"`
			GenerationIx struct {
				Champions struct {
					FrontShiny   string `json:"front_shiny"`
					FrontDefault string `json:"front_default"`
				} `json:"champions"`
				ScarletViolet struct {
					FrontFemale  interface{} `json:"front_female"`
					FrontDefault string      `json:"front_default"`
				} `json:"scarlet-violet"`
			} `json:"generation-ix"`
			GenerationVi struct {
				XY struct {
					BackShiny        string      `json:"back_shiny"`
					BackFemale       interface{} `json:"back_female"`
					FrontShiny       string      `json:"front_shiny"`
					BackDefault      string      `json:"back_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontDefault     string      `json:"front_default"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"x-y"`
				Icons struct {
					FrontFemale  interface{} `json:"front_female"`
					FrontDefault string      `json:"front_default"`
				} `json:"icons"`
				OmegarubyAlphasapphire struct {
					BackShiny        string      `json:"back_shiny"`
					BackFemale       interface{} `json:"back_female"`
					FrontShiny       string      `json:"front_shiny"`
					BackDefault      string      `json:"back_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontDefault     string      `json:"front_default"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
			} `json:"generation-vi"`
			GenerationIii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
				} `json:"icons"`
				Emerald struct {
					Animated struct {
						BackShiny    string `json:"back_shiny"`
						FrontShiny   string `json:"front_shiny"`
						BackDefault  string `json:"back_default"`
						FrontDefault string `json:"front_default"`
					} `json:"animated"`
					BackShiny    string `json:"back_shiny"`
					FrontShiny   string `json:"front_shiny"`
					BackDefault  string `json:"back_default"`
					FrontDefault string `json:"front_default"`
				} `json:"emerald"`
				RubySapphire struct {
					BackShiny    string `json:"back_shiny"`
					FrontShiny   string `json:"front_shiny"`
					BackDefault  string `json:"back_default"`
					FrontDefault string `json:"front_default"`
				} `json:"ruby-sapphire"`
				FireredLeafgreen struct {
					BackShiny    string `json:"back_shiny"`
					FrontShiny   string `json:"front_shiny"`
					BackDefault  string `json:"back_default"`
					FrontDefault string `json:"front_default"`
				} `json:"firered-leafgreen"`
			} `json:"generation-iii"`
			GenerationVii struct {
				Icons struct {
					FrontFemale  interface{} `json:"front_female"`
					FrontDefault string      `json:"front_default"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					BackShiny        string      `json:"back_shiny"`
					BackFemale       interface{} `json:"back_female"`
					FrontShiny       string      `json:"front_shiny"`
					BackDefault      string      `json:"back_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontDefault     string      `json:"front_default"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
				LetsGoPikachuLetsGoEevee struct {
					Icons struct {
						FrontDefault string `json:"front_default"`
					} `json:"icons"`
					BackShiny        string      `json:"back_shiny"`
					BackFemale       interface{} `json:"back_female"`
					FrontShiny       string      `json:"front_shiny"`
					BackDefault      string      `json:"back_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontDefault     string      `json:"front_default"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"lets-go-pikachu-lets-go-eevee"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontFemale  interface{} `json:"front_female"`
					FrontDefault string      `json:"front_default"`
				} `json:"icons"`
				BrilliantDiamondShiningPearl struct {
					FrontDefault string `json:"front_default"`
				} `json:"brilliant-diamond-shining-pearl"`
			} `json:"generation-viii"`
		} `json:"versions"`
		BackShiny        string      `json:"back_shiny"`
		BackFemale       interface{} `json:"back_female"`
		FrontShiny       string      `json:"front_shiny"`
		BackDefault      string      `json:"back_default"`
		FrontFemale      interface{} `json:"front_female"`
		FrontDefault     string      `json:"front_default"`
		BackShinyFemale  interface{} `json:"back_shiny_female"`
		FrontShinyFemale interface{} `json:"front_shiny_female"`
	} `json:"sprites"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	PastStats []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Stats []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		} `json:"stats"`
	} `json:"past_stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []interface{} `json:"past_types"`
}

func GetPokemonDetail(c *pokecache.Cache, name string) (PokemonDetails, error) {
	var pokemon PokemonDetails
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	cacheHit, ok := c.Get(url)
	if ok {
		err := json.Unmarshal(cacheHit, &pokemon)
		return pokemon, err
	}
	res, err := http.Get(url)
	if err != nil {
		return pokemon, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return pokemon, fmt.Errorf("No se encontro el pokemon: %s", name)
	}
	err = json.NewDecoder(res.Body).Decode(&pokemon)
	return pokemon, err
}
