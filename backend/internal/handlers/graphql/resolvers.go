package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

func (s *Server) getShipByMMSI(p graphql.ResolveParams) (interface{}, error) {
	mmsi, isOK := p.Args["mmsi"].(int)
	if !isOK {
		return nil, fmt.Errorf("invalid value for mmsi field: '%v'", p.Args["mmsi"])
	}
	ship, err := s.service.Get(p.Context, int32(mmsi))
	if err != nil {
		return nil, fmt.Errorf("error on getting ship with mmsi '%d': %w", mmsi, err)
	}
	return toDTO(ship), nil
}
