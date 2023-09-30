package searchgraph

import "github.com/mikeewhite/ship-locator/backend/internal/core/domain"

type ShipSearchResult struct {
	MMSI int32  `json:"mmsi"`
	Name string `json:"name"`
}

func toShipResultDTOs(searchResults []domain.ShipSearchResult) []ShipSearchResult {
	dtos := make([]ShipSearchResult, len(searchResults))
	for i := 0; i < len(searchResults); i++ {
		dtos[i] = ShipSearchResult{
			MMSI: searchResults[i].MMSI,
			Name: searchResults[i].Name,
		}
	}

	return dtos
}
