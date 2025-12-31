package main

import (
	"context"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuildWenGoAgent(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		envSetup   func() // 设置环境变量
		envCleanup func() // 清理环境变量
		args       args
		wantR      compose.Runnable[[]string, *schema.Message]
		wantErr    bool
	}{
		{
			name: "成功构建作文批改代理",
			envSetup: func() {
				os.Setenv("OCR_MODEL_NAME", "gpt-4o")
				os.Setenv("REORDER_MODEL_NAME", "gpt-4o")
				os.Setenv("FeedBack_Model_NAME", "gpt-4o")
			},
			envCleanup: func() {
				os.Unsetenv("OCR_MODEL_NAME")
				os.Unsetenv("REORDER_MODEL_NAME")
				os.Unsetenv("FeedBack_Model_NAME")
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			if tt.envSetup != nil {
				tt.envSetup()
			}
			// 恢复函数
			defer func() {
				if tt.envCleanup != nil {
					tt.envCleanup()
				}
			}()
			gotR, err := BuildWenGoAgent(tt.args.ctx)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NotNil(t, gotR)
			}

		})
	}
}
