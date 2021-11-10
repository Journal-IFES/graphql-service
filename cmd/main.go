package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Journal-IFES/graphql-service/internal/graphqlfields"
	"github.com/Journal-IFES/graphql-service/internal/postgres"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func signalHandler(ctx context.Context, cancel context.CancelFunc, s *http.Server, sigchan chan os.Signal) {
	for {
		select {
		case sig := <-sigchan:
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println(sig.String())
				if err := s.Shutdown(ctx); err != nil {
					log.Fatal(err)
				}
				cancel()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	router := gin.Default()

	s := &http.Server{
		Addr:           ":8081",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := postgres.InitPostgresDB("postgres://jifes:12345@127.0.0.1/jifes?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer postgres.ClosePostgresDB()

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: *graphqlfields.ModelsFields()}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: *graphqlfields.ModelsMutations()}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	router.Any("/graphql", gin.WrapH(h))

	done := make(chan uint8)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}

		close(done)
	}()

	go signalHandler(ctx, cancel, s, sigc)

	<-done
}
