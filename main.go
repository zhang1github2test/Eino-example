package main

import (
	"Eino-example/create_tool"
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/tool/httprequest/get"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
)

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	// 读取环境变量
	baseUrl := os.Getenv("OPENAI_BASE_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")

	ctx := context.TODO()

	chatModel, err := openai.NewChatModel(context.TODO(), &openai.ChatModelConfig{
		Model:   "qwen-flash", // 使用的模型版本
		APIKey:  apiKey,       // OpenAI API 密钥
		BaseURL: baseUrl,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建工具
	//searchTool, err := duckduckgo.NewTextSearchTool(context.TODO(), &duckduckgo.Config{})
	//if err != nil {
	//	log.Fatal(err)
	//
	//}

	// 第一种创建工具的方式：使用NewTool方式创建
	addTodoTool := create_tool.GetAddTodoTool()

	//第二种创建工具的方式：使用InferTool方式创建
	updateTodoTool := create_tool.GetUpdateTodoTool()

	//第三种创建工具的方式：实现InvokableTool接口
	listTodoTool := &create_tool.ListTodoTool{}

	// 第四种创建工具的方式：使用官方现有的工具
	newTool, err := get.NewTool(ctx, &get.Config{})
	if err != nil {
		log.Fatal(err)

	}

	// 创建ToolsNode
	conf := compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{newTool, addTodoTool, updateTodoTool, listTodoTool}, // 工具可以是 InvokableTool 或 StreamableTool
	}

	// 创建 agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig:      conf,
	})

	// 运行示例
	resp, err := agent.Generate(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "任务内容是：学习股票K线数据，开始时间是现在，结束是2天后。如果任务存在则更新任务，否则创建任务。",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	fmt.Println(resp.Content)

}

func reportStream(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()

	i := 0
	for {
		message, err := sr.Recv()
		if err == io.EOF { // 流式输出结束
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		log.Printf("message[%d]: %+v\n", i, message)
		i++
	}
}
