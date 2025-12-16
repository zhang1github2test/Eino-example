package einoagent

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"os"
)

// newModel 创建一个新的工具调用聊天模型实例
// 参数:
//
//	ctx - 上下文对象，用于控制请求的生命周期
//
// 返回值:
//
//	tmd - 工具调用聊天模型接口实例
//	err - 创建过程中可能发生的错误
func newModel(ctx context.Context) (tmd model.ToolCallingChatModel, err error) {
	// 读取环境变量
	baseUrl := os.Getenv("OPENAI_BASE_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")
	modelName := os.Getenv("MODEL_NAME")
	tmd, err = openai.NewChatModel(context.TODO(), &openai.ChatModelConfig{
		Model:   modelName, // 使用的模型版本
		APIKey:  apiKey,    // OpenAI API 密钥
		BaseURL: baseUrl,
	})
	return tmd, err
}
