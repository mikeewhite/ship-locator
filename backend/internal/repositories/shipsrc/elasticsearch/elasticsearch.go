package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type Repository struct {
	client    *elasticsearch.TypedClient
	indexName string
}

func New(ctx context.Context, cfg config.Config) (*Repository, error) {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses:     []string{cfg.ElasticsearchAddress},
		EnableMetrics: false,
		Logger:        nil,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create Repository client: %w", err)
	}

	repo := &Repository{
		client:    client,
		indexName: cfg.ElasticsearchIndex,
	}

	// Attempt to create indexes
	if err = repo.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	exists, err := client.Indices.Exists(repo.indexName).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if %s index exists: %w", repo.indexName, err)
	}
	if !exists {
		_, err := client.Indices.Create(repo.indexName).
			Request(&create.Request{
				Mappings: &types.TypeMapping{
					Properties: map[string]types.Property{
						"name": types.NewTextProperty(),
						"mmsi": types.NewTextProperty(),
					},
				},
			}).
			Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s index: %w", repo.indexName, err)
		}
	}

	return repo, nil
}

// TODO - Add metrics, logging, and tracing
func (r *Repository) Search(ctx context.Context, nameOrMMSI string) ([]domain.ShipSearchResult, error) {
	var results []domain.ShipSearchResult

	// create a match_phrase_prefix query to support matching before user has typed the full name
	matchPhasePrefixNameQuery := types.Query{
		MatchPhrasePrefix: map[string]types.MatchPhrasePrefixQuery{
			"name": {
				Query: nameOrMMSI,
			},
		},
	}

	// create a fuzzy query to support matching on misspelled names
	fuzzyNameQuery := types.Query{
		Match: map[string]types.MatchQuery{
			"name": {
				Query:     nameOrMMSI,
				Fuzziness: "AUTO", // e.g. 1 edit distance for strings of length 0-2, 2 edit distance for strings of length 3-5, etc.
			},
		},
	}

	fuzzyMMSIQuery := types.Query{
		Match: map[string]types.MatchQuery{
			"mmsi": {
				Query:     nameOrMMSI,
				Fuzziness: "AUTO", // e.g. 1 edit distance for strings of length 0-2, 2 edit distance for strings of length 3-5, etc.
			},
		},
	}

	resp, err := r.client.Search().
		Index(r.indexName).
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					// A should condition means that each sub-clause is optional, but at least one of them must match.
					// The first sub-clause has the highest precedence (i.e. matches are scored higher).
					Should: []types.Query{
						matchPhasePrefixNameQuery, // highest precedence
						fuzzyNameQuery,
						fuzzyMMSIQuery,
					},
				},
			},
		}).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search for ships: %w", err)
	}

	if resp.Hits.Total != nil && resp.Hits.Total.Value > 0 {
		for _, hit := range resp.Hits.Hits {
			var ssr domain.ShipSearchResult

			source_ := hit.Source_
			err := json.Unmarshal(source_, &ssr)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal ship search result: %w", err)
			}
			clog.Infof("source: %+v", string(source_))
			results = append(results, ssr)
		}
	}

	return results, nil
}

func (r *Repository) Index(ctx context.Context, ships []domain.ShipSearchResult) error {
	// TODO - bulk index ships
	upsert := true
	for _, ship := range ships {
		dto := toShipDTO(ship)
		query, err := dto.toJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal ship search data: %w", err)
		}

		id := fmt.Sprintf("%d", ship.MMSI)
		_, err = r.client.Update(r.indexName, id).Request(&update.Request{
			DocAsUpsert: &upsert,
			Doc:         query,
		}).Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to index ship data: %w", err)
		}
	}
	return nil
}

func (r *Repository) Shutdown(ctx context.Context) error {
	// noop
	return nil
}

func (r *Repository) createIndexes(ctx context.Context) error {
	exists, err := r.client.Indices.Exists(r.indexName).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if %s index exists: %w", r.indexName, err)
	}
	if !exists {
		_, err := r.client.Indices.Create(r.indexName).
			Request(&create.Request{
				Mappings: &types.TypeMapping{
					Properties: map[string]types.Property{
						"name": types.NewTextProperty(),
						"mmsi": types.NewTextProperty(),
					},
				},
			}).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to create %s index: %w", r.indexName, err)
		}
	}
	return nil
}
