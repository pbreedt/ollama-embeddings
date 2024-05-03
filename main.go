package main

import (
	"github/pbreedt/ollama-embeddings/chromem"
)

func main() {
	chromem.AutoEmbedding()

	chromem.ManualEmbedding()

	chromem.RunLLMWithEmbedding()
}
