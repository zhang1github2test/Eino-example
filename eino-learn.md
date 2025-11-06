
Eino 是基于 Golang 的 AI 应用开发框架。
Eino 官方文档 https://cloudwego.cn/zh/docs/eino/
仓库地址：https://github.com/cloudwego/eino，https://github.com/cloudwego/eino-ext
# 使用Eino开发第一个简单AI应用
```go
package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
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
		log.Println("Error loading .env file:", err)
	}
	// 读取环境变量
	baseUrl := os.Getenv("OPENAI_BASE_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")

	chatModel, err := openai.NewChatModel(context.TODO(), &openai.ChatModelConfig{
		Model:   "qwen-flash", // 使用的模型版本
		APIKey:  apiKey,       // OpenAI API 密钥
		BaseURL: baseUrl,
	})
	if err != nil {
		log.Print(err)
	}

	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)

	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})
	// 完整输出示例
	result, err := chatModel.Generate(context.TODO(), messages)
	if err != nil {
		log.Print(err)
	}
	println(result.Content)

	// 流式处理
	streamResult2, err := chatModel.Stream(context.TODO(), messages)
	if err != nil {
		log.Print(err)
		return
	}
	reportStream(streamResult2)

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

```
> 1、这里使用了godotenv加载环境变量，你需要在项目根目录下创建一个.env文件，文件内容如下：
OPENAI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1
OPENAI_API_KEY=sk-xxx-your-api-key
2、这里使用了openai.NewChatModel创建一个ChatModel实例，这个实例会调用OpenAI的API接口，并返回一个结果。
3、这里使用了prompt.FromMessages创建一个模板，这个模板会根据传入的参数生成一个消息列表。
4、这里使用了template.Format生成一个消息列表，这个消息列表会根据传入的参数生成一个消息列表。
5、这里使用了chatModel.Generate生成一个结果，这个结果会根据传入的参数生成一个结果。
6、这里使用了chatModel.Stream生成一个流式结果，这个结果会根据传入的参数生成一个流式结果。


## 大模型调用工具的示例
```go
package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/tool/httprequest/get"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
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

	newTool, err := get.NewTool(ctx, &get.Config{})
	if err != nil {
		log.Fatal(err)

	}

	// 创建ToolsNode
	conf := &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{newTool}, // 工具可以是 InvokableTool 或 StreamableTool
	}
	toolsNode, err := compose.NewToolNode(context.Background(), conf)

	// 获取工具信息并绑定到 ChatModel
	toolInfos := make([]*schema.ToolInfo, 0, 1)
	for _, tool := range conf.Tools {
		info, err := tool.Info(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		toolInfos = append(toolInfos, info)
	}
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		log.Fatal(err)
	}

	// 构建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(toolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// 运行示例
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "搜索一下 cloudwego/eino 的仓库地址",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	for _, msg := range resp {
		fmt.Println(msg.Content)
	}

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

```
上面的agent中，将会调用一个，输出工具调用的结果。

### 2.1 Enio框架中工具创建的四种方法
#### 1. 使用NewTool方式创建工具
   优点：直观，易于理解
   缺点： ToolInfo 中手动定义参数信息和实际的参数结构（TodoAddParams）需要在两个地方中定义，维护困难。
#### 2. 使用 InferTool 构建
   通过结构体的 tag 来定义参数信息，就能实现参数结构体和描述信息同源，无需维护两份信息。
#### 3. 使用实现 Tool 接口
   需要更多自定义逻辑的场景，可以通过实现 Tool 接口来创建：