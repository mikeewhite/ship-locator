package ports

type CollectorService interface {
	Process(mmsi int32, shipName string, latitude, longitude float64) error
}
