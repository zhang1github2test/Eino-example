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
	generateOutline(1)
}

const systemPromptTemplate = `# Role
	你是一位拥有 10 年经验的资深网文主编及剧情架构师，擅长规划长篇玄幻/仙侠小说的剧情节奏。你精通“波浪式推进”叙事法，能够确保长篇故事在 3000 章的篇幅中保持张力，不崩坏、不注水。
	
	# Goal
	根据提供的【小说剧情总纲】，为用户生成指定卷数的【分卷详细大纲】。
	
	# Constraints & Rules
	1. **严格一致性**：必须严格遵守总纲中设定的【章节范围】、【境界体系】、【地图转换】和【核心事件】。不得随意更改总纲中已定的关键节点（如高潮章节数、结局走向）。
	2. **节奏控制**：每卷大纲需体现“波浪式推进”，每 200-300 章需规划一个小高潮，卷末必须留有通往下一卷的“钩子”（悬念）。
	3. **爽点落实**：将总纲中的“核心爽点”具体化为剧情事件，确保主角成长曲线清晰。
	4. **伏笔埋设**：在分卷大纲中明确标注需要埋下的长期伏笔（参考总纲第 6 部分）。
	5. **输出格式**：必须使用 Markdown 格式，结构清晰，包含【卷名】、【章节范围】、【核心主题】、【剧情阶段分解】、【关键事件表】、【人物成长】、【伏笔与悬念】。
	
	# Workflow
	1. 分析用户指定的卷数及其在总纲中的定位。
	2. 将该卷的章节范围划分为 3-4 个主要剧情弧光（Arc）。
	3. 为每个弧光分配具体的章节区间、境界变化、地图场景和核心冲突。
	4. 检查是否符合总纲中的反派势力层级和金手指进化阶段。
	5. 输出最终大纲。
	
	# Tone
	专业、逻辑严密、具有创作启发性。`

const userPromptTemplate = `
	# Context Data (小说剧情总纲)
	请将以下内容作为核心参考依据，不可违背。总纲中应包含多卷的整体规划，请从中提取与【第{{.volumeNumber}}卷】相关的特定信息：
	"""
	  {{.outlineContent}}
	"""
	
	# Task
	请为我生成 **第{{.volumeNumber}}卷** 的详细分卷大纲。
	
	# Requirements
	1. **章节细化**：
	   - 请根据总纲规划，确定本卷的章节范围（例如：若总纲设定每卷约 700 章，请自动计算本卷的起始与结束章节）。
	   - 将本卷划分为 3-4 个剧情弧光，每个弧光需标明具体章节区间。
	   
	2. **境界对应**：
	   - 依据总纲中本卷对应的【境界阶段】，确保剧情推进与境界提升相匹配。
	   - 每个大境界突破需对应一个剧情高潮，境界名称以总纲设定为准。
	
	3. **反派互动**：
	   - 根据总纲【反派势力图谱】，提取本卷主要对手所在的势力层级（如：底层爪牙、中层干部、高层领袖等）。
	   - 设计具体的反派刁难与打脸情节，确保反派强度与本卷主角实力相匹配。
	
	4. **金手指演进**：
	   - 根据【金手指演进路线】，规划本卷金手指的具体成长幅度（例如：修复度从 X% 提升至 Y%）。
	   - 请规划解锁功能的具体剧情节点，确保不与前后卷冲突。
	
	5. **悬念结尾**：
	   - 卷末必须完成总纲为本卷设定的“核心大事件”（如：飞升、换地图、揭秘等）。
	   - 留下关于后续剧情或世界观深层秘密的悬念钩子。
	
	# Output Format
	请严格按照以下结构输出：
	## 第{{.volumeNumber}}卷：[卷名] 详细大纲
	### 1. 卷首信息
	- 章节范围：[根据总纲推算的本卷起止章节]
	- 对应境界：[本卷涉及的具体境界名称]
	- 对应地图：[本卷主要活动区域]
	- 核心主题：[本卷的核心冲突或目标]
	
	### 2. 剧情弧光分解
	- **弧光一：[名称] (第 X-X 章)**
	  - 剧情概要：
	  - 关键事件：
	  - 境界变化：
	  - 爽点/冲突：
	- **弧光二：[名称] (第 X-X 章)**
	  ...
	（根据实际划分的弧光数量继续输出）
	
	### 3. 关键节点规划表
	| 章节 | 事件名称 | 涉及人物 | 金手指状态 | 备注 |
	| :--- | :--- | :--- | :--- | :--- |
	| [关键章] | [本卷重要事件] | [相关角色] | [当前进度] | [总纲指定高潮/伏笔] |
	| ... | ... | ... | ... | ... |
	
	### 4. 伏笔与悬念
	- 本卷埋下的伏笔：[需与总纲长期线索呼应]
	- 本卷结尾钩子：[吸引读者阅读下一卷的关键点]
	
	### 5. 主编建议
	- 针对本卷写作节奏的具体建议（如：哪里需要放缓日常，哪里需要加快节奏）。
	- 注意本卷与上一卷及下一卷的衔接平滑度。
`

func getOutLineMessages(volumeNumber int) (result []*schema.Message, err error) {
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.GoTemplate,
		// 系统消息模板
		schema.SystemMessage(systemPromptTemplate),
		// 用户消息模板
		schema.UserMessage(userPromptTemplate),
	)
	outlineContent, err := ReadFileToString("xiaoshuo/Outline.md")
	if err != nil {
		return
	}

	// 使用模板生成消息
	result, err = template.Format(context.Background(), map[string]any{
		"volumeNumber":   volumeNumber,
		"outlineContent": outlineContent,
	})
	return
}
func generateOutline(volumeNumber int) {
	messages, err := getOutLineMessages(volumeNumber)
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
	reportStream(streamResult, volumeNumber)
}
func reportStream(sr *schema.StreamReader[*schema.Message], volumeNumber int) error {
	defer sr.Close()

	filename := fmt.Sprintf("xiaoshuo/第%d卷剧情大纲.md", volumeNumber)
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
