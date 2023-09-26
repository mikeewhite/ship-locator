package graphql

import (
	"github.com/graphql-go/graphql"
)

func (s *Server) getSchemaConfig() graphql.SchemaConfig {
	shipSearchResult := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "ShipSearchResult",
			Fields: graphql.Fields{
				"mmsi": &graphql.Field{
					Type: graphql.Int,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	shipType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Ship",
			Fields: graphql.Fields{
				"mmsi": &graphql.Field{
					Type: graphql.Int,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"latitude": &graphql.Field{
					Type: graphql.Float,
				},
				"longitude": &graphql.Field{
					Type: graphql.Float,
				},
				"lastUpdated": &graphql.Field{
					Type: graphql.DateTime,
				},
			},
		},
	)

	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"shipSearch": &graphql.Field{
					Type: graphql.NewList(shipSearchResult),
					Args: graphql.FieldConfigArgument{
						"searchTerm": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: s.lookupShipByNameOrMMSI,
				},
				"ship": &graphql.Field{
					Type: shipType,
					Args: graphql.FieldConfigArgument{
						"mmsi": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: s.getShipByMMSI,
				},
			},
		})

	return graphql.SchemaConfig{
		Query: rootQuery,
	}
}
