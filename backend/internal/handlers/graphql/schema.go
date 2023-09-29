package graphql

import (
	_ "embed"

	"github.com/graphql-go/graphql"
)

type service struct {
	Name    string
	Version string
	Schema  string
}

var name = "ship-data-service"
var version = "0.0.1"

//go:embed ship_service_schema.graphql
var schema string

func (s *Server) getSchemaConfig() graphql.SchemaConfig {
	serviceType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Service",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"version": &graphql.Field{
					Type: graphql.String,
				},
				"schema": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

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
				"service": &graphql.Field{
					Type: serviceType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return service{
							Name:    name,
							Version: version,
							Schema:  schema,
						}, nil
					},
				},
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
