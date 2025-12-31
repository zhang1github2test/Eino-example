package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"os"
)

// BuildWenGoAgent 构建作文批改AI代理
// 该函数创建一个包含OCR识别、作文重排序和作文批改三个步骤的图结构
func BuildWenGoAgent(ctx context.Context) (compose.Runnable[[]string, *schema.Message], error) {
	const (
		// 图节点名称常量
		imagePathsToQuery = "ImagePathsToQuery" // 图片路径查询节点
		ocrTemplate       = "OcrTemplate"       // OCR模板节点
		ocrModel          = "OcrModel"          // OCR模型节点
		reorderVariables  = "ReorderVariables"  // 重排序变量节点
		reorderTemplate   = "ReorderTemplate"   // 重排序模板节点
		reorderModel      = "ReorderModel"      // 重排序模型节点
		feedbackVariables = "FeedbackVariables" // 反馈变量节点
		feedbackTemplate  = "FeedbackTemplate"  // 反馈模板节点
		feedbackModel     = "FeedbackModel"     // 反馈模型节点
	)

	// 创建一个新的图结构，输入类型为 []string，输出类型为 *schema.Message
	g := compose.NewGraph[[]string, *schema.Message]()

	// 添加图片路径转换Lambda节点
	if err := g.AddLambdaNode(imagePathsToQuery, compose.InvokableLambdaWithOption(newLambdaInputMulti),
		compose.WithNodeName(imagePathsToQuery)); err != nil {
		return nil, fmt.Errorf("failed to add ImagePathsToQuery node: %w", err)
	}

	// 添加OCR模板节点
	orcTemplate, err := newOcrChatTemplate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create OCR template: %w", err)
	}
	if err := g.AddChatTemplateNode(ocrTemplate, orcTemplate, compose.WithNodeName(ocrTemplate)); err != nil {
		return nil, fmt.Errorf("failed to add OCR template node: %w", err)
	}

	// 添加OCR模型节点
	if err := addModelNode(ctx, g, os.Getenv("OCR_MODEL_NAME"), ocrModel); err != nil {
		return nil, fmt.Errorf("failed to create OCR model node: %w", err)
	}

	// 添加重排序变量转换Lambda节点
	if err := g.AddLambdaNode(reorderVariables, compose.InvokableLambdaWithOption(reOrderCompositionVariables),
		compose.WithNodeName(reorderVariables)); err != nil {
		return nil, fmt.Errorf("failed to add ReOrderVariables node: %w", err)
	}

	// 添加重排序模板节点
	reOrderTemplate, err := newReorderCompositionCTemplate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reorder template: %w", err)
	}
	if err := g.AddChatTemplateNode(reorderTemplate, reOrderTemplate,
		compose.WithNodeName(reorderTemplate)); err != nil {
		return nil, fmt.Errorf("failed to add reorder template node: %w", err)
	}

	// 添加重排序模型节点
	if err := addModelNode(ctx, g, os.Getenv("REORDER_MODEL_NAME"), reorderModel); err != nil {
		return nil, fmt.Errorf("failed to create reorder model node: %w", err)
	}

	// 添加反馈变量转换Lambda节点
	if err := g.AddLambdaNode(feedbackVariables, compose.InvokableLambdaWithOption(compositionFeedbackVariables),
		compose.WithNodeName(feedbackVariables)); err != nil {
		return nil, fmt.Errorf("failed to add FeedBackVariables node: %w", err)
	}

	// 添加反馈模板节点
	feedBackTemplate, err := newCompositionFeedbackTemplate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create feedback template: %w", err)
	}
	if err := g.AddChatTemplateNode(feedbackTemplate, feedBackTemplate, compose.WithNodeName(feedbackTemplate)); err != nil {
		return nil, fmt.Errorf("failed to add feedback template node: %w", err)
	}

	// 添加反馈模型节点
	if err := addModelNode(ctx, g, os.Getenv("FeedBack_Model_NAME"), feedbackModel); err != nil {
		return nil, fmt.Errorf("failed to create feedback model node: %w", err)
	}

	// 添加图边
	edges := [][2]string{
		{compose.START, imagePathsToQuery},
		{imagePathsToQuery, ocrTemplate},
		{ocrTemplate, ocrModel},
		{ocrModel, reorderVariables},
		{reorderVariables, reorderTemplate},
		{reorderTemplate, reorderModel},
		{reorderModel, feedbackVariables},
		{feedbackVariables, feedbackTemplate},
		{feedbackTemplate, feedbackModel},
		{feedbackModel, compose.END},
	}

	for _, edge := range edges {
		if err := g.AddEdge(edge[0], edge[1]); err != nil {
			return nil, fmt.Errorf("failed to add edge from %s to %s: %w", edge[0], edge[1], err)
		}
	}

	// 编译图结构为可执行的 Runnable 对象
	r, err := g.Compile(ctx,
		compose.WithGraphName("WenGoAgent"),
		compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, fmt.Errorf("failed to compile graph: %w", err)
	}

	return r, nil
}

// addModelNode 是一个辅助函数，用于创建模型节点
// 它封装了模型创建、Lambda创建和节点添加的逻辑
func addModelNode(ctx context.Context, g *compose.Graph[[]string, *schema.Message], modelName, nodeName string) error {
	if modelName == "" {
		return fmt.Errorf("model name is required for node %s", nodeName)
	}

	model, err := newModel(ctx, modelName)
	if err != nil {
		return fmt.Errorf("failed to create model %s: %w", modelName, err)
	}

	lambda, err := compose.AnyLambda(model.Generate, model.Stream, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create lambda for model %s: %w", modelName, err)
	}

	if err := g.AddLambdaNode(nodeName, lambda, compose.WithNodeName(nodeName)); err != nil {
		return fmt.Errorf("failed to add model node %s: %w", nodeName, err)
	}

	return nil
}
