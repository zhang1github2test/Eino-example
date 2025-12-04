package knowledgeindexing

import (
	"context"
	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino/components/document"
)

// newLoader creates a new document loader instance with default file loader configuration.
// It initializes a file loader with the provided context and empty configuration.
//
// Parameters:
//
//	ctx - context for the loader operations
//
// Returns:
//
//	ldr - the created document loader instance
//	err - error if loader creation fails, nil otherwise
func newLoader(ctx context.Context) (ldr document.Loader, err error) {
	config := &file.FileLoaderConfig{}
	ldr, err = file.NewFileLoader(ctx, config)
	if err != nil {
		return nil, err
	}
	return ldr, nil
}
