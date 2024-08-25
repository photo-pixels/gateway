package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/photo-pixels/gateway/internal/auth"
	"github.com/photo-pixels/gateway/internal/clients"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/photo-pixels/gateway/internal/graph"
	"github.com/photo-pixels/platform/log"
)

type Config struct {
	Host            string `yaml:"host"`
	HttpPort        int    `yaml:"port"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"`
}

type GraphQLServer struct {
	cfg    Config
	server *http.Server
	logger log.Logger
	auth   *auth.Auth
	cc     *clients.ServiceClientsContainer
}

func NewGraphQLServer(logger log.Logger,
	serverConfig Config,
	auth *auth.Auth,
	cc *clients.ServiceClientsContainer,
) *GraphQLServer {
	return &GraphQLServer{
		cfg:    serverConfig,
		logger: logger.Named("graphql_server"),
		auth:   auth,
		cc:     cc,
	}
}

// Start запустить сервер
func (s *GraphQLServer) Start(ctx context.Context) error {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			ServiceClientsContainer: s.cc,
			AuthService:             s.auth,
		},
		Directives: graph.DirectiveRoot{
			IsAuthenticated: s.auth.IsAuthenticated,
			IsToken:         s.auth.IsToken,
			SkipAuthenticate: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
				return next(ctx)
			},
		},
	}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", auth.TokenMiddleware(auth.SessionMiddleware(srv)))

	netAddress := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.HttpPort)
	s.server = &http.Server{
		Addr:    netAddress,
		Handler: nil,
	}

	go func() {
		s.logger.Infof("start server at %s", netAddress)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Errorf("failed to start server: %v", err)
		}
	}()

	return nil
}

// Stop остановить сервер
func (s *GraphQLServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.ShutdownTimeout))
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Errorf("failed to shutdown server: %v", err)
	}
}
