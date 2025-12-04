package main

import (
	"Eino-example/knowledgeindexing"
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/document"
	"github.com/joho/godotenv"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	err = indexMarkdownFiles(context.Background(), "knowledgeindexing/eino-doc")
	if err != nil {
		log.Fatal(err)
	}

}

func indexMarkdownFiles(ctx context.Context, dir string) error {
	runner, err := knowledgeindexing.BuildKnowledgeIndexing(ctx)
	if err != nil {
		return fmt.Errorf("build index graph failed: %w", err)
	}

	// 遍历 dir 下的所有 markdown 文件
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir failed: %w", err)
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			fmt.Printf("[skip] not a markdown file: %s\n", path)
			return nil
		}

		fmt.Printf("[start] indexing file: %s\n", path)

		ids, err := runner.Invoke(ctx, document.Source{URI: path})
		if err != nil {
			return fmt.Errorf("invoke index graph failed: %w", err)
		}

		fmt.Printf("[done] indexing file: %s, len of parts: %d\n", path, len(ids))

		return nil
	})

	return err
}
