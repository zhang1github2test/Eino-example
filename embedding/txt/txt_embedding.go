package main

import (
	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"log"
	"os"
)

func main() {
	context.TODO()
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	ctx := context.TODO()
	embedder, err := dashscope.NewEmbedder(ctx, &dashscope.EmbeddingConfig{
		APIKey: apiKey,
		Model:  "text-embedding-v3",
	})
	if err != nil {
		log.Printf("new embedder error: %v\n", err)
		return
	}
	var texts []string
	texts = append(texts, "1. Eiffel Tower: Located in Paris, France, it is one of the most famous landmarks in the world, designed by Gustave Eiffel and built in 1889.")
	embedStrings, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		log.Printf("embed strings error: %v\n", err)
		return
	}
	log.Printf("embed strings: %v\n", embedStrings)
}
