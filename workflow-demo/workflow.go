package workflow_demo

import (
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"golang.org/x/net/context"
	"log"
	"os"
)

func FlowDemo() {
	type LambdaInput struct {
		Input    string `json:"input"`
		Role     string `json:"role"`
		Query    string `json:"query"`
		Output   string `json:"output"`
		MetaData string `json:"meta_data"`
	}

	wf := compose.NewWorkflow[[]*schema.Message, LambdaInput]()

	// 读取环境变量
	baseUrl := os.Getenv("OPENAI_BASE_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")

	ctx := context.TODO()

	model, err := openai.NewChatModel(context.TODO(), &openai.ChatModelConfig{
		Model:   "qwen-vl-plus", // 使用的模型版本
		APIKey:  apiKey,         // OpenAI API 密钥
		BaseURL: baseUrl,
	})

	wf.AddChatModelNode("model", model).AddInput(compose.START)

	lambda1 := compose.InvokableLambda(
		func(ctx context.Context, input LambdaInput) (output LambdaInput, err error) {

			return input, nil
		},
	)

	wf.AddLambdaNode("l1", lambda1).AddInput("model", compose.MapFields("Content", "Input"))
	wf.End().AddInput("l1")
	runnable, err := wf.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	lout, err := runnable.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "你是谁？",
		},
	})

	log.Println(lout)
}
