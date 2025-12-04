package knowledgeindexing

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino/components/embedding"
	"os"
)

// newEmbedding creates a new embedding embedder instance using DashScope API.
// It initializes the embedder with the API key from environment variable and default model configuration.
//
// Parameters:
//
//	ctx - context for the embedding operations
//
// Returns:
//
//	embedding.Embedder - the created embedder instance
//	error - error if embedder creation fails, nil otherwise
func newEmbedding(ctx context.Context) (embedding.Embedder, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("DASHSCOPE_API_KEY")

	// Create new DashScope embedder with specified configuration
	embedder, err := dashscope.NewEmbedder(ctx, &dashscope.EmbeddingConfig{
		APIKey: apiKey,
		Model:  "text-embedding-v3",
	})
	if err != nil {
		return nil, fmt.Errorf("new embedder error: %w", err)
	}
	return embedder, nil
}
