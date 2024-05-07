package data

import (
	"strconv"

	"github.com/philippgille/chromem-go"
	"github.com/tmc/langchaingo/schema"
)

const (
	HowMuchDoesLlamaWeigh = "How much does a llama weigh?"
	WhatAnimalsAreRelated = "What animals are llamas related to?"
	HowTallIsAvgLlama     = "How tall is the average llama?"
)

func Queries() []string {
	return []string{
		HowMuchDoesLlamaWeigh,
		WhatAnimalsAreRelated,
		HowTallIsAvgLlama,
	}
}

func EmbedData() []string {
	return []string{
		"Llamas are members of the camelid family meaning they're pretty closely related to vicuñas and camels",
		"Llamas were first domesticated and used as pack animals 4,000 to 5,000 years ago in the Peruvian highlands",
		"Llamas can grow as much as 6 feet tall though the average llama is between 5 feet 6 inches and 5 feet 9 inches tall",
		"Llamas weigh between 280 and 450 pounds and can carry 25 to 30 percent of their body weight",
		"Llamas are vegetarians and have very efficient digestive systems",
		"Llamas live to be about 20 years old, though some only live for 15 years and others live to be 30 years old",
	}
}

func LangchainDocs() []schema.Document {
	var docs []schema.Document
	for _, d := range EmbedData() {
		docs = append(docs, schema.Document{PageContent: d})
	}

	return docs
}

func ChromemDocs() []chromem.Document {
	var docs []chromem.Document
	for i, d := range EmbedData() {
		docs = append(docs, chromem.Document{ID: strconv.Itoa(i), Content: d})
	}

	return docs
}
