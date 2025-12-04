package knowledgeindexing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

const (
	indexName               = "eino_example_index"
	fieldContent            = "content"
	fieldContentDenseVector = "content_dense_vector"
	metaData                = "meta"
)

// newIndexer creates a new Elasticsearch indexer instance with embedding support.
// It initializes the ES client, creates the index if needed, and sets up the document indexing pipeline.
//
// Parameters:
//   - ctx: context for controlling the lifecycle of the operation
//
// Returns:
//   - idr: the created indexer instance for indexing documents
//   - err: error if the indexer creation fails, nil otherwise
func newIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	// es supports multiple ways to connect
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")
	host := os.Getenv("ES_URL")

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{host},
		Username:  username,
		Password:  password,
	})
	if err != nil {
		log.Fatalf("NewClient of es8 failed, err=%v", err)
	}

	// create index if needed.
	// comment out the code if index has been created.
	if err = createIndex(ctx, client); err != nil {
		log.Fatalf("createIndex of es8 failed, err=%v", err)
	}
	embedder, err := newEmbedding(ctx)

	// create es indexer component
	indexer, err := es8.NewIndexer(ctx, &es8.IndexerConfig{
		Client:    client,
		Index:     indexName,
		BatchSize: 10,
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (field2Value map[string]es8.FieldValue, err error) {
			metadataBytes, err := json.Marshal(doc.MetaData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal metadata: %w", err)
			}
			return map[string]es8.FieldValue{
				fieldContent: {
					Value:    doc.Content,
					EmbedKey: fieldContentDenseVector, // vectorize doc content and save vector to field "content_vector"
				},
				metaData: {
					Value: string(metadataBytes),
				},
			}, nil
		},
		Embedding: embedder, // replace it with real embedding component
	})
	if err != nil {
		log.Fatalf("NewIndexer of es8 failed, err=%v", err)
	}
	return indexer, err
}
