package searchgraph

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const tracerName = "github.com/mikeewhite/ship-locator/graphql/searchgraph"

func (s *Server) lookupShipByNameOrMMSI(p graphql.ResolveParams) (interface{}, error) {
	tr := otel.Tracer(tracerName)
	// See https://opentelemetry.io/docs/specs/otel/trace/semantic_conventions/instrumentation/graphql/
	ctx, span := tr.Start(p.Context, fmt.Sprintf("%s %s", p.Info.Operation.GetOperation(), p.Info.FieldName))
	defer span.End()

	searchTerm, isOK := p.Args["searchTerm"].(string)
	if !isOK {
		return nil, fmt.Errorf("invalid value for searchTerm field: '%v'", p.Args["searchTerm"])
	}

	span.SetAttributes(attribute.Key("searchTerm").String(searchTerm))
	ships, err := s.shipServiceService.Search(ctx, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("error on searching for ships with searchTerm '%s': %w", searchTerm, err)
	}

	return toShipResultDTOs(ships), nil
}
