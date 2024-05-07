package main

import (
	"github/pbreedt/ollama-embeddings/chromem"
	"github/pbreedt/ollama-embeddings/langchain"
)

func main() {
	pgVectorURL := langchain.RunPGVector()
	langchain.RunLangchain(pgVectorURL)
	langchain.TerminateContainer()

	chromem.AutoEmbedding()
	chromem.ManualEmbedding()
	chromem.RunLLMWithEmbedding()

	// ollama.RunWithoutEmbeddings()
	// ollama.RunOllamaChat()
}
