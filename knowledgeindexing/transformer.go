package knowledgeindexing

import (
	"context"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
)

// NewTransformer creates a new document transformer for processing markdown headers.
// It initializes a header splitter with specific configuration to handle title headers.
//
// Parameters:
//
//	ctx - The context for the transformer creation
//
// Returns:
//
//	tfr - The created document transformer instance
//	err - Error if transformer creation fails, nil otherwise
func NewTransformer(ctx context.Context) (tfr document.Transformer, err error) {
	// Configure header processing to map '#' to 'title' field
	config := &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "title",
			"##":  "subtitle",
			"###": "subsubtitle",
		},
		TrimHeaders: false}
	tfr, err = markdown.NewHeaderSplitter(ctx, config)
	if err != nil {
		return nil, err
	}
	return tfr, nil
}
