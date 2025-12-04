package knowledgeindexing

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/exists"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/densevectorsimilarity"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// createIndex 创建 Elasticsearch 索引
// 该函数用于在 Elasticsearch 中创建一个新的索引，并定义索引的映射结构。
// 映射包括以下字段：
// - fieldContent: 文本类型字段，用于存储文档内容
// - fieldExtraLocation: 文本类型字段，用于存储额外的位置信息
// - fieldContentDenseVector: 密集向量类型字段，用于存储文档内容的向量表示
//   - Dims: 向量维度，设置为1024（与嵌入维度相同）
//   - Index: 是否建立索引，设置为true
//   - Similarity: 相似度计算方法，使用余弦相似度
//
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - client: Elasticsearch 客户端实例
//
// 返回值:
//   - error: 如果创建索引过程中发生错误，则返回相应的错误信息
func createIndex(ctx context.Context, client *elasticsearch.Client) error {
	existsResp, err := exists.NewExistsFunc(client)(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if existsResp {
		return nil
	}

	_, err = create.NewCreateFunc(client)(indexName).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: map[string]types.Property{
				fieldContent: types.NewTextProperty(),
				metaData:     types.NewTextProperty(),
				fieldContentDenseVector: &types.DenseVectorProperty{
					Dims:       of(1024), // same as embedding dimensions
					Index:      of(true),
					Similarity: of(densevectorsimilarity.Cosine),
				},
			},
		},
	}).Do(ctx)

	return err
}

func of[T any](v T) *T {
	return &v
}
