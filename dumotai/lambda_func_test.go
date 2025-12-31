package main

import (
	"context"
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MockImageConverter 用于测试的模拟实现
type MockImageConverter struct {
	responses map[string]struct {
		mimeType string
		data     string
		err      error
	}
}

func (m *MockImageConverter) ImageToBase64(path string) (string, string, error) {
	response, exists := m.responses[path]
	if !exists {
		return "", "", nil // 默认返回空值，无错误
	}
	return response.mimeType, response.data, response.err
}

// 测试函数实现
func Test_newLambdaInputMulti(t *testing.T) {
	// 保存原始转换器
	originalConverter := imageConverter
	defer func() {
		imageConverter = originalConverter // 恢复原始转换器
	}()

	type args struct {
		ctx        context.Context
		imagePaths []string
		opts       []any
	}

	tests := []struct {
		name          string
		args          args
		mockResponses map[string]struct {
			mimeType string
			data     string
			err      error
		}
		wantOutput map[string]any
		wantErr    bool
	}{
		{
			name: "正常图片转换",
			args: args{
				ctx:        context.Background(),
				imagePaths: []string{"image1.jpg", "image2.png"},
				opts:       []any{},
			},
			mockResponses: map[string]struct {
				mimeType string
				data     string
				err      error
			}{
				"image1.jpg": {"image/jpeg", "base64data1", nil},
				"image2.png": {"image/png", "base64data2", nil},
			},
			wantOutput: map[string]any{
				"userInputMultiContent": []*schema.Message{{
					Role: schema.User,
					UserInputMultiContent: []schema.MessageInputPart{
						{
							Type: schema.ChatMessagePartTypeImageURL,
							Image: &schema.MessageInputImage{
								MessagePartCommon: schema.MessagePartCommon{
									Base64Data: stringPtr("base64data1"),
									MIMEType:   "image/jpeg",
								},
							},
						},
						{
							Type: schema.ChatMessagePartTypeImageURL,
							Image: &schema.MessageInputImage{
								MessagePartCommon: schema.MessagePartCommon{
									Base64Data: stringPtr("base64data2"),
									MIMEType:   "image/png",
								},
							},
						},
					},
				}},
			},
			wantErr: false,
		},
		{
			name: "包含无效图片路径",
			args: args{
				ctx:        context.Background(),
				imagePaths: []string{"valid.jpg", "invalid.jpg", "another.jpg"},
				opts:       []any{},
			},
			mockResponses: map[string]struct {
				mimeType string
				data     string
				err      error
			}{
				"valid.jpg":   {"image/jpeg", "base64data1", nil},
				"invalid.jpg": {"", "", assert.AnError}, // 模拟错误
				"another.jpg": {"image/jpeg", "base64data2", nil},
			},
			wantOutput: map[string]any{
				"userInputMultiContent": []*schema.Message{{
					Role: schema.User,
					UserInputMultiContent: []schema.MessageInputPart{
						{
							Type: schema.ChatMessagePartTypeImageURL,
							Image: &schema.MessageInputImage{
								MessagePartCommon: schema.MessagePartCommon{
									Base64Data: stringPtr("base64data1"),
									MIMEType:   "image/jpeg",
								},
							},
						},
						{
							Type: schema.ChatMessagePartTypeImageURL,
							Image: &schema.MessageInputImage{
								MessagePartCommon: schema.MessagePartCommon{
									Base64Data: stringPtr("base64data2"),
									MIMEType:   "image/jpeg",
								},
							},
						},
					},
				}},
			},
			wantErr: false,
		},
		{
			name: "空图片路径数组",
			args: args{
				ctx:        context.Background(),
				imagePaths: []string{},
				opts:       []any{},
			},
			mockResponses: map[string]struct {
				mimeType string
				data     string
				err      error
			}{},
			wantOutput: map[string]any{
				"userInputMultiContent": []*schema.Message{{
					Role:                  schema.User,
					UserInputMultiContent: nil,
				}},
			},
			wantErr: false,
		},
		{
			name: "所有路径都无效",
			args: args{
				ctx:        context.Background(),
				imagePaths: []string{"invalid1.jpg", "invalid2.jpg"},
				opts:       []any{},
			},
			mockResponses: map[string]struct {
				mimeType string
				data     string
				err      error
			}{

				"invalid1.jpg": {"", "", assert.AnError},
				"invalid2.jpg": {"", "", assert.AnError},
			},
			wantOutput: map[string]any{
				"userInputMultiContent": []*schema.Message{{
					Role:                  schema.User,
					UserInputMultiContent: nil,
				}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock转换器
			mockConverter := &MockImageConverter{
				responses: tt.mockResponses,
			}
			imageConverter = mockConverter

			gotOutput, err := newLambdaInputMulti(tt.args.ctx, tt.args.imagePaths, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("newLambdaInputMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 比较输出
			if !assert.ObjectsAreEqual(tt.wantOutput, gotOutput) {
				t.Errorf("newLambdaInputMulti() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

// 辅助函数：创建字符串指针
func stringPtr(s string) *string {
	return &s
}
