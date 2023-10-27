package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kubeagi/arcadia/graphql-server/go-server/graph"
	"github.com/kubeagi/arcadia/graphql-server/go-server/pkg/client"
)

func main() {
	//srv := handler.NewDefaultServer(generated.NewExecutableSchema(starwars.NewResolver()))
	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	cfg := ctrl.GetConfigOrDie()

	client.InitClient(cfg)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
		rc := graphql.GetFieldContext(ctx)
		fmt.Println("Entered", rc.Object, rc.Field.Name)
		res, err = next(ctx)
		fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
		return res, err
	})

	http.Handle("/", playground.Handler("Arcadia", "/query"))
	http.Handle("/query", srv)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
