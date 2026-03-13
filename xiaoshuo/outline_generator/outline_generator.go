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
	generateOutline()
}

const systemPromptTemplate = `你是一位精通千万字级长篇玄幻网文架构的资深主编。你擅长规划 3000 章以上体量的小说结构，深知如何通过“换地图”、“升级体系”和“反派梯队”来维持长篇故事的张力。

你的任务是根据用户提供的【核心设定】，生成一份能支撑 3000 章篇幅的【剧情总纲】。

请严格遵守以下长篇架构规则：
1. **体量控制**：规划必须足以支撑 3000 章内容。避免剧情推进过快。
2. **地图层级**：必须设计至少 5-7 个层级分明的地图，每个地图对应主角的一个大境界阶段。
3. **力量体系**：设计至少 9-13 个修炼境界，确保每个境界能支撑 200-300 章的剧情。
4. **反派梯队**：设计层层递进的反派势力，确保每个地图都有对应的冲突源头。
5. **金手指进化**：金手指必须具有成长性，能在不同地图阶段解锁新功能。
6. **伏笔长线**：必须设计贯穿全书的终极悬念，在最终卷揭晓。

输出结构必须包含：
1. 【长篇架构总览】
2. 【修炼与地图体系】
3. 【剧情分卷总纲】（按地图/境界分卷，每卷注明：预估章节范围）
4. 【反派势力图谱】
5. 【金手指演进路线】`

const userPromptTemplate = `
	请基于以下【核心设定】，构建一份能支撑 **3000 章** 篇幅的小说【剧情总纲】。
	
	=== 核心设定开始 ===
	{{.CoreSettings}}
	=== 核心设定结束 ===
	
	执行要求：
	1. **篇幅规划**：请将剧情划分为 5-8 个大卷，每卷预估章节数需合理分配，总和需匹配 3000 章体量。
	2. **换地图逻辑**：明确每一卷结束时的“换地图”契机，确保主角有理由离开当前地图进入更高级地图。
	3. **战力控制**：确保主角在每一卷结束时刚好突破大境界，避免战力溢出。
	4. **内容聚焦**：不需要生成细纲（章节级），只需要卷宗级的宏观剧情，但必须说明每一卷的核心爽点和高潮。
	5. **长期伏笔**：请在总纲中标注哪些设定是为后期（2000 章后）准备的伏笔。
	
	现在，请生成适配 3000 章篇幅的剧情总纲。
`

func getOutLineMessages() (result []*schema.Message, err error) {
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.GoTemplate,
		// 系统消息模板
		schema.SystemMessage(systemPromptTemplate),
		// 用户消息模板
		schema.UserMessage(userPromptTemplate),
	)
	CoreSettings, err := ReadFileToString("xiaoshuo/Core_settings.md")
	if err != nil {
		return
	}

	// 使用模板生成消息
	result, err = template.Format(context.Background(), map[string]any{
		"CoreSettings": CoreSettings,
	})
	return
}
func generateOutline() {
	messages, err := getOutLineMessages()
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
	reportStream(streamResult)
}
func reportStream(sr *schema.StreamReader[*schema.Message]) error {
	defer sr.Close()

	filename := "xiaoshuo/Outline.md"
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
