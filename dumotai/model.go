package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"os"
)

func newModel(ctx context.Context, modelName string) (tmd model.ToolCallingChatModel, err error) {
	// 读取环境变量
	baseUrl := os.Getenv("OPENAI_BASE_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")
	tmd, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   modelName, // 使用的模型版本
		APIKey:  apiKey,    // OpenAI API 密钥
		BaseURL: baseUrl,
	})
	return tmd, err
}
