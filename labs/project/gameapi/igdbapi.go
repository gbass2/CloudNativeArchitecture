package gameapi

import (
	"github.com/Henry-Sarabia/igdb/v2"
	"fmt"
)

var (
  key = "bjf52rgd650jp9r6i1ah3dqjh6sjw0"
	token = "291rkwhbn61ew66uzpd3xxohzeysia"
	client = igdb.NewClient(key, token, nil) // Http client to connecto igdb api.
)

// Queries IGDB for game information for a given name.
// Returns a slice of a map containg the returned data.
func SearchGame(name string) ([]map[string]interface{}, error) {
		var data []map[string]interface{} // Slice of map of the returned data.

		// Set the api return parameters.
		gameParams := igdb.ComposeOptions( // Options for quering games
				igdb.SetLimit(3),
				igdb.SetFields("*"),
		)

		platParams := igdb.ComposeOptions( // Options for quering platforms.
				igdb.SetFields("*"),
		)

		dateParams := igdb.ComposeOptions( // Options for quering platforms.
				igdb.SetFields("*"),
				igdb.SetLimit(1),
		)

		involParams := igdb.ComposeOptions( // Options for quering involved companies.
				igdb.SetFields("company","developer"),
		)

		companyParams := igdb.ComposeOptions( // Options for quering companies.
				igdb.SetFields("name"),
		)

		// Query api.
		games, err := client.Games.Search(name, gameParams)

		if err != nil {
				return data, err
		}

		// Take the response and parse it into the map
		// Get the data from games that does not require another query
		for _, game := range games {
				result := make(map[string]interface{})
				result["name"] = game.Name
				result["summary"] = game.Summary
				result["ratings"] = fmt.Sprintf("%.2f", game.AggregatedRating)

				// Query for platforms the games are on and add it to a slice.
				platforms, err := client.Platforms.List(game.Platforms, platParams)
				var platsSlice []string // Temporarily holds the platforms associated with game

				if err != nil {
					  result["platforms"] = nil
						goto date
				}

				for _, platform := range platforms {
						platsSlice = append(platsSlice, platform.Name)
				}

				// Add the platforms to the map.
				result["platforms"] = platsSlice

				date:
				// Get the release date of the game
				releaseDates, err := client.ReleaseDates.List(game.ReleaseDates, dateParams)

				if err != nil {
						result["date"] = nil
						goto companies
				}

				result["date"] = releaseDates[0].Human

				companies:
				// Query for involved companies the games are on and add it to a slice.
				companies, err := client.Companies.List([]int{-1}, companyParams)
				involvedCompanies, err := client.InvolvedCompanies.List(game.InvolvedCompanies, involParams)

				var involSlice []int // Temporarily holds the platforms associated with game
				var compSlice []string // Temporarily holds the platforms associated with game

				if err != nil {
					  result["companies"] = nil
						goto end
				}

				for _, company := range involvedCompanies {
						if company.Developer {
							involSlice = append(involSlice, company.Company)
						}
				}

				// Query for companies the games are on and add it to a slice.
				companies, err = client.Companies.List(involSlice, companyParams)

				if err != nil {
						result["companies"] = nil
						goto end
				}

				for _, company := range companies {
						compSlice = append(compSlice, company.Name)
				}

				// Add the platforms to the map.
				result["companies"] = compSlice

				end:
				data = append(data, result)
		}

	  return data, nil
}
