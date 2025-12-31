package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"log"
)

func newLambdaInputMulti(ctx context.Context, imagePaths []string, opts ...any) (output map[string]any, err error) {
	var result []schema.MessageInputPart

	for i, path := range imagePaths {
		// 处理每个路径
		fmt.Printf("Index: %d, Path: %s\n", i, path)
		mimeType, data, err := imageConverter.ImageToBase64(path)
		if err != nil {
			log.Println("")
			continue
		}
		// 根据需要创建 MessageInputPart 对象
		part := schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeImageURL,
			Image: &schema.MessageInputImage{
				MessagePartCommon: schema.MessagePartCommon{
					Base64Data: &data,
					MIMEType:   mimeType,
				},
			},
		}
		result = append(result, part)
	}
	return map[string]any{
		"userInputMultiContent": []*schema.Message{{
			Role:                  schema.User,
			UserInputMultiContent: result,
		},
		},
	}, nil
}

func reOrderCompositionVariables(ctx context.Context, content *schema.Message, opts ...any) (output map[string]any, err error) {
	return map[string]any{
		"content": content.Content,
	}, nil
}

func compositionFeedbackVariables(ctx context.Context, content *schema.Message, opts ...any) (output map[string]any, err error) {
	c := content.Content
	var v Composition
	err = json.Unmarshal([]byte(c), &v)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"title":               v.Title,
		"composition_content": v.Content,
	}, nil
}
