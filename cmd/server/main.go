package main

import (
	"flag"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"post-comments/pkg/database"
	"post-comments/pkg/generated"
	"post-comments/pkg/resolver"
	"post-comments/pkg/storage"
)

func main() {

	// Define flag for storage type
	storageType := flag.String("storage", "in_memory", "type of storage to use (postgres or in_memory)")
	flag.Parse()

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	var store storage.Storage
	if *storageType == "postgres" {
		cfg := database.LoadDBConfig()
		connect, err := database.NewDB(cfg)
		if err != nil {
			log.Fatalf("failed to initialize db: %s", err.Error())
		}
		defer connect.Close()
		store = storage.NewPostgresStorage(connect)
	} else {
		store = storage.NewInMemoryStorage()
	}

	r := resolver.NewResolver(store)
	// Create a GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	// Create a playground for testing
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", viper.GetString("port"))
	log.Fatal(http.ListenAndServe(viper.GetString("port"), nil))
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
