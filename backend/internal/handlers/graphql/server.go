package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/rs/cors"

	"github.com/mikeewhite/ship-locator/backend/internal/core/ports"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type Server struct {
	httpServer http.Server
	service    ports.ShipService
	schema     *graphql.Schema
}

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operationName"`
	Variables map[string]interface{} `json:"variables"`
}

const endpoint = "/graphql"

func New(cfg config.Config, service ports.ShipService) (*Server, error) {
	s := &Server{
		service: service,
	}

	schema, err := graphql.NewSchema(s.getSchemaConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to initialise schema: %w", err)
	}
	s.schema = &schema

	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, s.HandleQuery)
	s.httpServer = http.Server{
		Addr:    cfg.GraphQLAddress,
		Handler: cors.Default().Handler(mux),
	}

	return s, nil
}

func (s *Server) Serve(ctx context.Context) error {
	clog.Infof("Starting GraphQL server at %s", s.httpServer.Addr)
	go func() {
		<-ctx.Done()
		clog.Info("Stopping GraphQL server")
		s.Shutdown()
	}()
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		clog.Errorf("failed to gracefully shutdown graphql server: %w", err)
	}
}

func (s *Server) HandleQuery(w http.ResponseWriter, r *http.Request) {
	var p postData
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(400)
		return
	}
	result := s.executeQuery(p)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		clog.Errorf("error on encoding graphql response: %s", err.Error())
	}
}

func (s *Server) executeQuery(postData postData) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         *s.schema,
		RequestString:  postData.Query,
		VariableValues: postData.Variables,
		OperationName:  postData.Operation,
		Context:        context.Background(),
	})
	if len(result.Errors) > 0 {
		clog.Errorf("error returned from graphQL API: %v", result.Errors)
	}
	return result
}
