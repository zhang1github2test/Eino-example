package einoagent

import (
	"context"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// BuildEinoAgent 构建一个用于处理用户消息并生成响应的可运行图结构（Agent）。
// 该函数通过组合多个节点（如 Lambda、ChatTemplate、Retriever 等）构建一个有向图，
// 并编译成可以执行的 Runnable 对象。
//
// 参数:
//   - ctx: 上下文对象，用于控制生命周期和传递元数据。
//
// 返回值:
//   - r: 实现了 compose.Runnable 接口的对象，可用于执行整个流程。
//   - err: 如果在构建过程中发生错误，则返回相应的错误信息；否则为 nil。
func BuildEinoAgent(ctx context.Context) (r compose.Runnable[*UserMessage, *schema.Message], err error) {
	const (
		InputToQuery   = "InputToQuery"
		ChatTemplate   = "ChatTemplate"
		ReactAgent     = "ReactAgent"
		RedisRetriever = "RedisRetriever"
		InputToHistory = "InputToHistory"
	)

	// 创建一个新的图结构，输入类型为 *UserMessage，输出类型为 *schema.Message
	g := compose.NewGraph[*UserMessage, *schema.Message]()

	// 添加将用户输入转换为查询语句的 Lambda 节点
	_ = g.AddLambdaNode(InputToQuery, compose.InvokableLambdaWithOption(newLambda), compose.WithNodeName("UserMessageToQuery"))

	// 初始化聊天模板，并将其作为 ChatTemplate 节点加入图中
	chatTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate, chatTemplateKeyOfChatTemplate)

	// 初始化 ReAct Agent 的 Lambda 函数，并添加到图中
	reactAgentKeyOfLambda, err := newLambda1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(ReactAgent, reactAgentKeyOfLambda, compose.WithNodeName("ReAct Agent"))

	// 初始化 Redis 检索器，并添加到图中，指定其输出键为 "documents"
	redisRetrieverKeyOfRetriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(RedisRetriever, redisRetrieverKeyOfRetriever, compose.WithOutputKey("documents"))

	// 添加将用户输入转为历史变量的 Lambda 节点
	_ = g.AddLambdaNode(InputToHistory, compose.InvokableLambdaWithOption(newLambda2), compose.WithNodeName("UserMessageToVariables"))

	// 定义图中的边关系：START -> InputToQuery 和 START -> InputToHistory 表示流程开始时同时触发这两个节点
	_ = g.AddEdge(compose.START, InputToQuery)
	_ = g.AddEdge(compose.START, InputToHistory)

	// 流程结束节点连接 ReactAgent
	_ = g.AddEdge(ReactAgent, compose.END)

	// 数据流定义：
	// InputToQuery -> RedisRetriever：使用查询结果进行检索
	// RedisRetriever -> ChatTemplate：将检索结果传入聊天模板
	// InputToHistory -> ChatTemplate：将历史上下文传入聊天模板
	// ChatTemplate -> ReactAgent：最终由 ReactAgent 处理生成回复
	_ = g.AddEdge(InputToQuery, RedisRetriever)
	_ = g.AddEdge(RedisRetriever, ChatTemplate)
	_ = g.AddEdge(InputToHistory, ChatTemplate)
	_ = g.AddEdge(ChatTemplate, ReactAgent)

	// 编译图结构为可执行的 Runnable 对象，设置图名称及节点触发模式为所有前驱完成后再触发
	r, err = g.Compile(ctx, compose.WithGraphName("EinoAgent"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}

	return r, err
}
