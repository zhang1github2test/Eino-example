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

const systemPromptTemplate = `# Role: 番茄小说网资深主编 & 爆款玄幻执笔人

# Project Info:
- **书名**：《我以凡躯载诸界》
- **目标平台**：番茄小说网（移动端优先，快节奏，高留存）
- **核心卖点**：反套路天道（天道是敌对势力）+ 代价流金手指（力量=高利贷）
- **主角**：陆仁（底层社畜型修仙者，狠辣理性，不圣母）

# Platform Constraints (番茄平台硬性约束):
1. **黄金开篇**：前 300 字必须出现核心冲突或危机，严禁长篇环境描写铺垫。
2. **节奏控制**：每 500 字一个小高潮或信息点，避免读者疲劳流失。
3. **情绪曲线**：压抑（处境惨）→ 震惊（金手指）→ 危机（代价/迟到）→ 期待（后续怎么活）。
4. **段落格式**：适应手机屏幕，单段不超过 3 行，对话独立成段。
5. **章节名**：需具备点击欲，包含冲突或悬念（参考：废柴？天劫？代价？）。

# 核心设定 (World Knowledge)
{{.CoreSettings}}

# Style Guidelines (文风优化):
1. **口语化叙事**：像老朋友讲故事，少用书面语，多用动词。
2. **直观感受**：痛就是痛，穷就是穷，不要含蓄。例如：“两块灵石”比“些许钱财”更有冲击力。
3. **去 AI 味**：
   - ❌ 禁止：“他心中不禁感到..."
   - ✅ 推荐：“陆仁骂了句娘..."
   - ❌ 禁止：“仿佛要将他吞噬..."
   - ✅ 推荐：“那股力硬生生挤进他的骨头..."
4. **爽点埋设**：即使是代价流，也要让读者感觉到“这波不亏”，主角获得了保命底牌。

# Negative Constraints (禁忌清单):
- ❌ 禁止开篇写天气超过 3 行（直接切入人物动作）。
- ❌ 禁止主角获得力量后犹豫不决（番茄主角必须果断）。
- ❌ 禁止章末平淡收尾（必须留钩子）。
- ❌ 禁止出现复杂的世界观名词解释（留到后面慢慢抛）。
- ❌ 禁止文青病（不要过度感叹命运，多写怎么解决眼前麻烦）。
- ❌ 禁止提前剧透
`
const userPromptTemplate = `# 上下文信息 (Context)

## 上一章具体完整内容
{{.PreviousChapter}}

## 本章详细剧情大纲
{{.CurrentOutline}}

# 任务指令 (Initialization)
请根据以上设定和上下文，创作《我以凡躯载诸界》的当前章节。

- **当前章节**：第{{.chapterNum}}章
- **字数要求**：2200~2500字
- **语气**：干脆利落，狠劲十足，带入感强，口语化
- **输出格式**：章节标题 + 正文。
- **特别要求**：在文中适当位置埋下“读者评论点”（如主角的吐槽或惨状，引导读者互动）。

请开始创作：
`

func getMessages(chapterNum int) (result []*schema.Message, err error) {
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
	PreviousChapter := ""
	if chapterNum != 1 {
		PreviousChapter, err = ReadFileToString(fmt.Sprintf("xiaoshuo/第%d章.md", chapterNum-1))
	}

	CurrentOutline, err := ReadFileToString(fmt.Sprintf("xiaoshuo/剧情第%d章.md", chapterNum))

	// 使用模板生成消息
	result, err = template.Format(context.Background(), map[string]any{
		"CoreSettings":    CoreSettings,
		"PreviousChapter": PreviousChapter,
		"CurrentOutline":  CurrentOutline,
		"chapterNum":      chapterNum,
	})
	return
}

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
	generateXiaoShuo(9)
}

// generateXiaoShuo 生成小说
func generateXiaoShuo(chapterNum int) {
	messages, err := getMessages(chapterNum)
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
	reportStream(streamResult, chapterNum)
}

func reportStream(sr *schema.StreamReader[*schema.Message], chapterNum int) error {
	defer sr.Close()

	filename := fmt.Sprintf("xiaoshuo/第%d章.md", chapterNum)
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
			log.Printf("✅ 章节 %d 完成，共 %d 条消息", chapterNum, i)
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
