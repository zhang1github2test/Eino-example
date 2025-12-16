package einoagent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino-ext/components/retriever/es8/search_mode"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"log"
	"os"
)

const (
	indexName          = "eino_example_index"
	fieldContent       = "content"
	fieldContentVector = "content_dense_vector"
	metaData           = "meta"
)

// newRetriever 创建一个新的 Elasticsearch 检索器（Retriever），用于执行向量相似度搜索。
// 它从环境变量中读取 Elasticsearch 的认证信息和地址，并初始化一个嵌入模型来支持语义检索。
//
// 参数：
//   - ctx: 上下文对象，用于控制请求的生命周期。
//
// 返回值：
//   - rtr: 实现了 retriever.Retriever 接口的对象，可用于执行文档检索操作。
//   - err: 如果在创建过程中发生错误，则返回相应的错误信息。
func newRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {

	// 从环境变量获取 Elasticsearch 认证信息与主机地址
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")
	host := os.Getenv("ES_URL")

	// 初始化 Elasticsearch 客户端
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{host},
		Username:  username,
		Password:  password,
	})
	if err != nil {
		log.Fatalf("NewClient of es8 failed, err=%v", err)
	}

	// 初始化文本嵌入模型
	embedder, err := newEmbedding(ctx)
	if err != nil {
		log.Fatalf("NewEmbedding failed, err=%v", err)
	}

	// 构造并初始化 Retriever 对象，配置使用稠密向量余弦相似度进行搜索
	rtr, err = es8.NewRetriever(ctx, &es8.RetrieverConfig{
		Client: client,
		Index:  indexName,
		TopK:   5,
		SearchMode: search_mode.SearchModeDenseVectorSimilarity(
			search_mode.DenseVectorSimilarityTypeCosineSimilarity,
			fieldContentVector,
		),
		// 自定义结果解析函数：将 Elasticsearch 响应转换为标准 Document 结构
		ResultParser: func(ctx context.Context, hit types.Hit) (doc *schema.Document, err error) {
			doc = &schema.Document{
				ID:       *hit.Id_,
				Content:  "",
				MetaData: map[string]any{},
			}

			var src map[string]any
			if err = json.Unmarshal(hit.Source_, &src); err != nil {
				return nil, err
			}

			// 遍历源数据字段，填充到 Document 中
			for field, val := range src {
				switch field {
				case fieldContent:
					doc.Content = val.(string)

				case fieldContentVector:
					var v []float64
					for _, item := range val.([]interface{}) {
						v = append(v, item.(float64))
					}
					doc.WithDenseVector(v)

				case metaData:
					metadata, err := StringToMap(val.(string))
					if err != nil {
						return nil, err
					}
					doc.MetaData = metadata

				default:
					return nil, fmt.Errorf("unexpected field=%s, val=%v", field, val)
				}
			}

			// 设置匹配得分
			if hit.Score_ != nil {
				doc.WithScore(float64(*hit.Score_))
			}

			return doc, nil
		},
		Embedding: embedder,
	})

	return rtr, err
}

// StringToMap 将 JSON 字符串转换为 map[string]any
func StringToMap(s string) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal([]byte(s), &result)
	return result, err
}
