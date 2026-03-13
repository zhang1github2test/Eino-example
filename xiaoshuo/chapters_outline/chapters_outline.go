package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
)

func ReadFileToString(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	generateChaptersOutline(1, 1, 50)
}

func getOutLineMessages(volumeNumber, startChapter, endChapter int) (result []*schema.Message, err error) {
	systemPromptTemplate, err := ReadFileToString("xiaoshuo/chapters_outline/SystemPrompt.md")
	userPromptTemplate, err := ReadFileToString("xiaoshuo/chapters_outline/UserPromptTemplate.md")
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.GoTemplate,
		// 系统消息模板
		schema.SystemMessage(systemPromptTemplate),
		// 用户消息模板
		schema.UserMessage(userPromptTemplate),
	)
	VolumeOutline, err := ReadFileToString("xiaoshuo/Outline.md")
	if err != nil {
		return
	}
	PreviousBatchStateSnapshot := ""
	if startChapter != 1 {
		PreviousBatchStateSnapshot = ""
	}

	// 使用模板生成消息
	result, err = template.Format(context.Background(), map[string]any{
		"VolumeNumber":               volumeNumber,
		"VolumeOutline":              VolumeOutline,
		"BookTitle":                  "我以凡躯载诸界",
		"StartChapter":               startChapter,
		"EndChapter":                 endChapter,
		"PreviousBatchStateSnapshot": PreviousBatchStateSnapshot,
	})
	return
}
func generateChaptersOutline(volumeNumber, startChapter, endChapter int) {
	messages, err := getOutLineMessages(volumeNumber, startChapter, endChapter)
	if err != nil {
		panic(err)
	}
	// 读取环境变量
	baseUrl := os.Getenv("OPENAI_BASE_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")
	modelName := os.Getenv("MODEL_NAME")
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		Model:   modelName, // 使用的模型版本
		APIKey:  apiKey,    // OpenAI API 密钥
		BaseURL: baseUrl,
	})

	streamResult, err := chatModel.Stream(context.Background(), messages)
	reportStream(streamResult, volumeNumber, startChapter, endChapter)
}
func reportStream(sr *schema.StreamReader[*schema.Message], volumeNumber, startChapter, endChapter int) error {
	defer sr.Close()

	filename := fmt.Sprintf("xiaoshuo/第%d卷%d至%d章大纲.md", volumeNumber, startChapter, endChapter)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 4096) // 4KB 缓冲区

	i := 0
	for {
		message, err := sr.Recv()
		if err == io.EOF {
			writer.Flush()
			log.Printf("✅ 剧情大纲生成完成，共 %d 条消息", i)
			return nil
		}
		if err != nil {
			writer.Flush() // 出错时也尝试刷新
			return fmt.Errorf("接收失败: %w", err)
		}
		if message != nil && message.Content != "" {
			if _, err := writer.WriteString(message.Content); err != nil {
				return fmt.Errorf("写入失败: %w", err)
			}

			// 每 20 条消息刷新一次，平衡性能与安全性
			if i%20 == 0 {
				if err := writer.Flush(); err != nil {
					log.Printf("⚠️ 刷新缓冲区警告: %v", err)
				}
			}
		}
		i++
		meta := message.ResponseMeta
		log.Println("返回的元数据", meta.Usage)
	}
}
