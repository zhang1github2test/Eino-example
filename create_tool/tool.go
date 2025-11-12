package create_tool

import (
	"context"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"log"
)

type TodoAddParams struct {
	Content   string `json:"content" jsonschema:"desc=The content of the todo item"`
	StartedAt int64  `json:"started_at" jsonschema:"desc=The started time of the todo item, in unix timestamp"`
	Deadline  int64  `json:"deadline" jsonschema:"desc=The deadline of the todo item, in unix timestamp"`
}

// 处理函数
func AddTodoFunc(_ context.Context, params *TodoAddParams) (string, error) {
	// Mock处理逻辑
	return `{"msg": "add todo success"}`, nil
}

// 工具创建的第一种方式：使用NewTool方式创建
func GetAddTodoTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "add_todo",
		Desc: "Add a todo item",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"content": {
				Desc:     "The content of the todo item",
				Type:     schema.String,
				Required: true,
			},
			"started_at": {
				Desc: "The started time of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
			"deadline": {
				Desc: "The deadline of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
		}),
	}

	// 使用NewTool创建工具
	return utils.NewTool(info, AddTodoFunc)
}

func UpdateTodoFunc(_ context.Context, params *TodoAddParams) (string, error) {
	log.Println("update todo params: %+v", params)
	// Mock处理逻辑
	return `{"msg": "update todo success"}`, nil
}

// 工具创建的第二种方式：使用InferTool方式创建
func GetUpdateTodoTool() tool.InvokableTool {
	// 使用 InferTool 创建工具
	updateTool, err := utils.InferTool(
		"update_todo", // tool name
		"Update a todo item, eg: content,deadline...", // tool description
		UpdateTodoFunc)
	if err != nil {
		panic(err)
	}
	return updateTool
}

type ListTodoTool struct{}

func (lt *ListTodoTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "list_todo",
		Desc: "List all todo items",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"finished": {
				Desc:     "filter todo items if finished",
				Type:     schema.Boolean,
				Required: false,
			},
		}),
	}, nil
}

func (lt *ListTodoTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	log.Println("查询所有待办事项，参数：%s", argumentsInJSON)
	// Mock调用逻辑
	return `{"todos": [{"id": "1", "content": "在2024年12月10日之前完成Eino项目演示文稿的准备工作", "started_at": 1717401600, "deadline": 1717488000, "done": false},{"id": "2", "content": "在2025年11月9日之前完成学习股票K线数据", "started_at": 1717401600, "deadline": 1717488000, "done": false}]}`, nil
}
