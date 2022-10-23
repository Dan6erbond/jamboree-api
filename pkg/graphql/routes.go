package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/dan6erbond/jamboree-api/graph/generated"
	"github.com/dan6erbond/jamboree-api/pkg/auth"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dan6erbond/jamboree-api/graph"
)

func RegisterRoutes(router *mux.Router, resolver *graph.Resolver, logger *zap.Logger) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	router.Use(auth.Middleware())

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql")).Methods("GET", "POST", "OPTIONS")
	router.Handle("/graphql", srv).Methods("GET", "POST", "OPTIONS")
}
