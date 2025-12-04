package knowledgeindexing

import (
	"context"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
)

// BuildKnowledgeIndexing 构建知识索引流程图
// 该函数通过编排不同的节点来构建一个从文档源到Elasticsearch索引的处理流程。
// 流程包括:
// 1. 使用 FileLoader 加载文件数据
// 2. 通过 MarkdownSplitter 对加载的文档进行切片处理
// 3. 利用 EsIndexer 将处理后的文档片段索引至 Elasticsearch
//
// 参数:
//   - ctx: 上下文，用于控制请求生命周期
//
// 返回值:
//   - r: 可执行的流程图实例，输入为 document.Source，输出为 []string
//   - err: 错误信息，如果构建过程中出现错误则返回具体错误信息
func BuildKnowledgeIndexing(ctx context.Context) (r compose.Runnable[document.Source, []string], err error) {
	const (
		FileLoader       = "FileLoader"
		MarkdownSplitter = "MarkdownSplitter"
		EsIndexer        = "EsIndexer"
	)
	graph := compose.NewGraph[document.Source, []string]()

	fileLoaderKeyOfLoader, err := newLoader(ctx)
	if err != nil {
		return nil, err
	}
	_ = graph.AddLoaderNode(FileLoader, fileLoaderKeyOfLoader)
	markdownSplitterKeyOfTransformer, err := NewTransformer(ctx)
	if err != nil {
		return nil, err
	}

	_ = graph.AddDocumentTransformerNode(MarkdownSplitter, markdownSplitterKeyOfTransformer)
	esIndexerKeyOfIndexer, err := newIndexer(ctx)
	if err != nil {
		return nil, err
	}
	_ = graph.AddIndexerNode(EsIndexer, esIndexerKeyOfIndexer)
	_ = graph.AddEdge(compose.START, FileLoader)
	_ = graph.AddEdge(FileLoader, MarkdownSplitter)
	_ = graph.AddEdge(MarkdownSplitter, EsIndexer)
	_ = graph.AddEdge(EsIndexer, compose.END)
	compile, err := graph.Compile(ctx)

	return compile, err
}
