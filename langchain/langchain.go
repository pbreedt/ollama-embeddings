package langchain

import (
	"context"
	"fmt"
	"github/pbreedt/ollama-embeddings/data"
	"log"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

/*
A number of embedding providers are available via Langchain
chroma, pgvector (postgres with vector capabilities), milvus, pinecone, weaviate, qdrant
*/

func RunLangchain(pgVectorURL string) {
	// Using Chroma as embeddings store
	// store, errNs := chroma.New(
	// 	// also available chroma.WithEmbedder(ef),
	// 	chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
	// 	chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
	// 	chroma.WithDistanceFunction(chroma_go.COSINE),
	// 	chroma.WithNameSpace(uuid.New().String()),
	// )
	// if errNs != nil {
	// 	log.Fatalf("new: %v\n", errNs)
	// }

	// Use seperate ollama client for embeddings
	emb, err := ollama.New(
		ollama.WithModel("nomic-embed-text"),
	)
	if err != nil {
		log.Fatal("client error:", err)
	}

	// MUST provide embedder to pgvector.New
	e, err := embeddings.NewEmbedder(emb)
	ctx := context.Background()
	store, err := pgvector.New(ctx,
		pgvector.WithConnectionURL(pgVectorURL),
		pgvector.WithEmbedder(e),
		pgvector.WithPreDeleteCollection(true), // delete before create
		pgvector.WithCollectionName("embeddings"),
	)

	type meta = map[string]any

	// Add documents to the vector store.
	_, errAd := store.AddDocuments(context.Background(), data.LangchainDocs())
	if errAd != nil {
		log.Fatalf("AddDocument error: %v\n", errAd)
	}

	q1 := "what is the average weight of a llama?" // Answer not directly provided in data
	docs, err := store.SimilaritySearch(ctx, q1, 1, vectorstores.WithScoreThreshold(0))
	if err != nil {
		log.Fatalf("Search error: %v\n", errAd)
	}
	fmt.Printf("Query: %s\n", q1)
	fmt.Printf("Langchain Emb.Result:\n%s\n", getResponse(docs))

	// Use seperate ollama client for general LLM
	llm, err := ollama.New(
		ollama.WithModel("gemma:2b"),
	)
	if err != nil {
		log.Fatal("client error:", err)
	}

	result, err := llm.Call(ctx, fmt.Sprintf("Using the following information: %s\nAnswer this question: %s", getResponse(docs), q1),
		llms.WithTemperature(0.1),
		// llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		// 	fmt.Print(string(chunk))
		// 	return nil
		// }),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query: %s\n", q1)
	fmt.Printf("Langchain LLM.Result:\n%s\n", result)

	questions := data.Queries()
	for _, q := range questions {
		docs, err = store.SimilaritySearch(ctx, q, 1, vectorstores.WithScoreThreshold(0))
		if err != nil {
			log.Fatalf("Search error: %v\n", errAd)
		}
		fmt.Printf("Query: %s\n", q)
		fmt.Printf("Langchain Emb.Result:\n%s\n", getResponse(docs))

		result, err := llm.Call(ctx, fmt.Sprintf("Using the following information: %s\nAnswer this question: %s", getResponse(docs), q),
			llms.WithTemperature(0.1),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				// 	fmt.Print(string(chunk))
				return nil
			}),
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Query: %s\n", q)
		fmt.Printf("Langchain LLM.Result:\n%s\n", result)
	}

}

func getResponse(docs []schema.Document) string {
	texts := make([]string, len(docs))
	for docI, doc := range docs {
		texts[docI] = doc.PageContent
	}
	return strings.Join(texts, "\n")
}
