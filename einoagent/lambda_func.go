package einoagent

import (
	"context"
	"time"
)

// newLambda 创建一个新的 Lambda 函数处理器
// 该函数接收用户消息并直接返回消息中的查询内容
//
// 参数:
//
//	ctx - 上下文对象，用于控制函数执行的生命周期
//	input - 用户消息指针，包含待处理的查询数据
//	opts - 可变参数，用于传递额外的配置选项
//
// 返回值:
//
//	output - 处理结果字符串，即原始查询内容
//	err - 错误信息，当前实现始终返回 nil
func newLambda(ctx context.Context, input *UserMessage, opts ...any) (output string, err error) {
	return input.Query, nil
}

// newLambda2 创建一个新的lambda函数处理用户消息
//
// 参数:
//
//	ctx - 上下文对象，用于控制请求的生命周期
//	input - 用户消息输入，包含查询内容和历史记录
//	opts - 可变参数列表，用于传递额外选项
//
// 返回值:
//
//	output - 包含处理结果的映射表，包含content、history和date字段
//	err - 错误信息，如果处理成功则返回nil
func newLambda2(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	return map[string]any{
		"content": input.Query,
		"history": input.History,
		"date":    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
