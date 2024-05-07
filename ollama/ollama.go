package ollama

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func RunWithoutEmbeddings() {
	model := "gemma:2b"
	prompt1 := "how many planets are there?"
	prompt2 := "Who was the first man to walk on the moon?" // Some limits to gemma:2b's knowledge

	llm, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	// Simplified interface for single prompt
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:\n", completion)

	// Alternative using llms.Model.GenerateContent
	msgs := []llms.MessageContent{
		// {Role: llms.ChatMessageTypeHuman, Parts: []llms.ContentPart{llms.TextPart(prompt2)}},
		llms.TextParts(llms.ChatMessageTypeHuman, prompt2), // short-hand version of line above
	}
	opts := llms.WithModel(model)
	response, err := llm.GenerateContent(ctx, msgs, opts)

	fmt.Println("Response:\n", response.Choices[0].Content)

}

func RunOllamaChat() {

	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		api.Message{
			Role:    "system",
			Content: "Provide very brief, concise responses",
		},
		api.Message{
			Role:    "user",
			Content: "Name some unusual animals",
		},
		api.Message{
			Role:    "assistant",
			Content: "Monotreme, platypus, echidna",
		},
		api.Message{
			Role:    "user",
			Content: "which of these is the most dangerous?",
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    "llama3",
		Messages: messages,
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	}

	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
}
