package graphql

import (
	"github.com/graphql-go/graphql"
)

func (s *Server) getSchemaConfig() graphql.SchemaConfig {
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

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
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
		Query: queryType,
	}
}
