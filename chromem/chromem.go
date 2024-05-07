package chromem

import (
	"context"
	"fmt"
	"github/pbreedt/ollama-embeddings/data"
	"log"
	"runtime"
	"strconv"

	"github.com/ollama/ollama/api"
	"github.com/philippgille/chromem-go"
)

// use "text-embedding-3-small" for OpenAI
// chromem-go:v0.6.0 does not work properly with 'mxbai-embed-large',
// authors did warn that package is not ready for production use before v1.0.0
var embModel string = "nomic-embed-text"

func AutoEmbedding() {
	ctx := context.Background()

	db := chromem.NewDB()
	// use ollama embedding
	c, err := db.GetOrCreateCollection("llamas", nil, chromem.NewEmbeddingFuncOllama(embModel, "http://localhost:11434/api"))
	// use openAI embedding
	// c, err := db.GetOrCreateCollection("colours", nil, nil)
	if err != nil {
		panic(err)
	}

	dt := data.ChromemDocs()
	// Embedding not provided, chromem-go will calculate and add
	err = c.AddDocuments(ctx, dt, runtime.NumCPU())
	if err != nil {
		panic(err)
	}

	query := data.Queries()[0]
	fmt.Printf("Query: %s\n", query)
	if embModel == "nomic-embed-text" {
		// "nomic-embed-text" specific prefix (not required with OpenAI's or other models)
		query = "search_query: " + query
	}
	res, err := c.Query(ctx, query, 1, nil, nil)
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		fmt.Printf("Cromem Result: %s (Similarity:%v)\n", r.Content, r.Similarity)
	}

	// fmt.Printf("ID: %v\nSimilarity: %v\nContent: %v\n", res[0].ID, res[0].Similarity, res[0].Content)
}

func ManualEmbedding() {
	docs := data.EmbedData()
	ctx := context.Background()
	db := chromem.NewDB()
	// db, err := chromem.NewPersistentDB("./emb.db", false)

	// still need to provide NewEmbeddingFuncOllama for use in query
	col, err := db.GetOrCreateCollection("llamas", nil, chromem.NewEmbeddingFuncOllama(embModel, "http://localhost:11434/api"))
	if err != nil {
		panic(err)
	}

	// ollama client used for embedding calc
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal("ollama client error:", err)
	}
	if col.Count() == 0 {
		dbdocs := []chromem.Document{}
		for i, d := range docs {
			embReq := api.EmbeddingRequest{
				Model:  embModel,
				Prompt: d,
			}
			embRes, err := client.Embeddings(ctx, &embReq)
			if err != nil {
				log.Fatal("ollama client error:", err)
			}
			// ollama return emb in float64, chromem works with float32
			f32 := make([]float32, 0)
			for _, f := range embRes.Embedding {
				f32 = append(f32, float32(f))
			}
			dbdocs = append(dbdocs, chromem.Document{
				ID:        strconv.Itoa(i),
				Embedding: f32,
				Content:   d,
			})
		}

		err = col.AddDocuments(ctx, dbdocs, runtime.NumCPU())
		if err != nil {
			panic(err)
		}
	}

	quesions := data.Queries()
	for _, q := range quesions {
		if embModel == "nomic-embed-text" {
			// "nomic-embed-text" specific prefix (not required with OpenAI's or other models)
			q = "search_query: " + q
		}
		res, err := col.Query(ctx, q, 1, nil, nil)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Query: %s\n", q)
		fmt.Printf("Cromem Result: %s (Similarity:%v)\n", res[0].Content, res[0].Similarity)
	}
}
