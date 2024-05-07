package langchain

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var pgVectorContainer *postgres.PostgresContainer

func RunPGVector() string {
	pgVectorContainer, err := postgres.RunContainer(
		context.Background(),
		testcontainers.WithImage("docker.io/pgvector/pgvector:pg16"),
		postgres.WithDatabase("db_test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("passw0rd!"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil && strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
		log.Fatalf("starting pgvector container error: %v", err)
	}

	pgvectorURL, err := pgVectorContainer.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		log.Fatalf("connection string error: %v", err)
	}

	return pgvectorURL
}

func TerminateContainer() {
	if pgVectorContainer != nil {
		pgVectorContainer.Terminate(context.Background())
	}
}
