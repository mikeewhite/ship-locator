package domain

type ShipSearchResult struct {
	MMSI int32
	Name string
}

func NewShipSearchResult(mmsi int32, name string) ShipSearchResult {
	shipSearchResult := ShipSearchResult{
		MMSI: mmsi,
		Name: name,
	}
	return shipSearchResult
}
