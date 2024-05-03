package chromem

import (
	"context"
	"fmt"
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
	c, err := db.GetOrCreateCollection("colours", nil, chromem.NewEmbeddingFuncOllama(embModel, "http://localhost:11434/api"))
	// use openAI embedding
	// c, err := db.GetOrCreateCollection("colours", nil, nil)
	if err != nil {
		panic(err)
	}

	// Embedding not provided, chromem-go will calculate and add
	err = c.AddDocuments(ctx, []chromem.Document{
		{
			ID:      "1",
			Content: "The sky is blue because of Rayleigh scattering.",
		},
		{
			ID:      "2",
			Content: "Leaves are green because chlorophyll absorbs red and blue light.",
		},
	}, runtime.NumCPU())
	if err != nil {
		panic(err)
	}

	query := "Why is the sky blue?"
	if embModel == "nomic-embed-text" {
		// "nomic-embed-text" specific prefix (not required with OpenAI's or other models)
		query = "search_query: " + query
	}
	res, err := c.Query(ctx, query, 1, nil, nil)
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		fmt.Printf("ID: %v\nSimilarity: %v\nContent: %v\n", r.ID, r.Similarity, r.Content)
	}

	// fmt.Printf("ID: %v\nSimilarity: %v\nContent: %v\n", res[0].ID, res[0].Similarity, res[0].Content)
}

func ManualEmbedding() {
	docs := []string{
		"Llamas are members of the camelid family meaning they're pretty closely related to vicuñas and camels",
		"Llamas were first domesticated and used as pack animals 4,000 to 5,000 years ago in the Peruvian highlands",
		"Llamas can grow as much as 6 feet tall though the average llama is between 5 feet 6 inches and 5 feet 9 inches tall",
		"Llamas weigh between 280 and 450 pounds and can carry 25 to 30 percent of their body weight",
		"Llamas are vegetarians and have very efficient digestive systems",
		"Llamas live to be about 20 years old, though some only live for 15 years and others live to be 30 years old",
	}
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

	query := "What animals are llamas related to?" //"How much does a llama weigh?"
	if embModel == "nomic-embed-text" {
		// "nomic-embed-text" specific prefix (not required with OpenAI's or other models)
		query = "search_query: " + query
	}
	res, err := col.Query(ctx, query, 1, nil, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %v\nSimilarity: %v\nContent: %v\n", res[0].ID, res[0].Similarity, res[0].Content)
}
