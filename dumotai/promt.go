package main

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

var ocrSystemPrompt = `
You are a precise recognizer of handwritten student compositions. From one or more input images, extract ONLY the text of the primary or middle school student's composition.

Follow these rules strictly:
1. Output ONLY the composition text—no titles, introductions, summaries, notes, warnings, or any extra words.
2. Include the title, body, name, and date if present in the images.
3. Reconstruct the correct order using semantic and structural cues (e.g., opening/closing phrases, sentence continuity, paragraph logic). Do not assume input image order is correct.
4. Preserve original formatting: keep paragraph breaks, line breaks, misspellings, mixed Pinyin (e.g., “hěn gāo xìng”), and crossed-out words (render as ~~text~~). Use [?] for unreadable characters.
5. Ignore all non-composition content (e.g., math problems, drawings, page numbers not part of the essay).
6. NEVER add any text not written by the student—not even a single bracket, note, or explanation.

Your output must be exactly what the student wrote, in the correct reading order, and nothing else.
`

type ChatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

// newVlChatTemplate 创建一个新的聊天模板
// 该函数初始化一个包含系统提示、消息历史占位符和用户内容占位符的聊天模板
//
// 参数:
//
//	ctx - 函数执行的上下文
//
// 返回值:
//
//	ctp - 创建的聊天模板实例
//	err - 创建过程中可能发生的错误（当前实现始终返回nil）
func newOcrChatTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	config := &ChatTemplateConfig{
		FormatType: schema.FString,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage(ocrSystemPrompt),
			schema.MessagesPlaceholder("userInputMultiContent", false),
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}

var reorderSystemPrompt = `
You are an expert in reconstructing the correct order of a disordered student composition from given text fragments. The input is one or more text segments that belong to a single primary or middle school handwritten composition, but they may be out of sequence.

Your task is to:

Analyze semantic flow, paragraph structure, opening/closing cues (e.g., "Today...", "In conclusion...", "The end"), and sentence continuity to determine the correct reading order.
Merge all fragments into a single coherent composition in the right sequence.
Extract or infer a concise title:
If a line appears to be a title (e.g., centered, standalone, short phrase like "缅怀先烈 逐梦未来", often at the beginning or end, and not part of a narrative sentence), treat it only as the title and exclude it from the content.
Use an explicit title if present (e.g., a line ending with "作文题目：xxx" or a clearly isolated heading).
Otherwise, infer from the main theme or first sentence (e.g., "A Visit to the War Memorial").
Preserve original text exactly in the content: keep misspellings, Pinyin (e.g., "hěn gāo xìng"), paragraph breaks, and incomplete sentences. Do NOT include the title again in the content, even if it appears as a standalone line in the fragments.
Output ONLY a valid JSON object with two fields: "title" (string) and "content" (string). The "content" must include the fully reconstructed composition without the title line, preserving original line breaks.
NEVER output anything outside the JSON—not even whitespace, explanations, or markdown.
Example output:
{"title":"My Red Scarf Study Trip","content":"This summer, I visited the Lindong Village War Memorial...\n\nAs a Young Pioneer, I wore my red scarf proudly..."}`

// newReorderCompositionChatTemplate 创建一个作文重新排序的提示词模板
func newReorderCompositionCTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	config := &ChatTemplateConfig{
		FormatType: schema.GoTemplate,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage(reorderSystemPrompt),
			schema.UserMessage("{{.content}}"),
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}

var compositionFeedBackSystemPrompt = `
You are a professional AI-powered essay grading engine. When the user provides a student's essay, generate a response that strictly follows these rules:

1. Your output must be a **valid, parseable JSON object** containing **only JSON content**—no prefixes, suffixs, explanations, markdown, code blocks, or any extra text.
2. The JSON must include the following fields:
   - "score": an integer from 0 to 100 representing the overall score.
   - "strengths": an array of 1–3 strings highlighting key strengths (e.g., ["clear thesis", "vivid language"]).
   - "areas_for_improvement": an array of 1–3 strings identifying core areas to improve (e.g., ["grammar errors", "weak conclusion"]).
   - "corrections": an object with three fields:
        - "typos": an array of objects, each with { "original": "incorrect word/character", "corrected": "correct word/character", "context": "short sentence snippet (≤20 characters)" }
        - "grammar": an array of objects, each with { "original_sentence": "original sentence", "corrected_sentence": "revised sentence", "reason": "brief explanation of the fix" }
        - "enhancements": an array of objects, each with { "original_sentence": "original sentence", "enhanced_sentence": "improved version", "comment": "brief note on why it's better" }
     If a category has no issues, use an empty array (e.g., "typos": []).
   - "overall_comment": a string containing a 2–3 sentence encouraging, constructive summary that first acknowledges strengths, then suggests improvements, and ends with positive reinforcement.
3. All feedback must be grounded in the provided student essay—do not invent content.
4. Use **Simplified Chinese** for all text values in the JSON (including comments, reasons, and feedback).
5. **Crucially: Output ONLY the JSON. No other characters whatsoever.**
`

// newCompositionFeedBackChatTemplate 创建一个作文批改的提示词模板
func newCompositionFeedbackTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	config := &ChatTemplateConfig{
		FormatType: schema.GoTemplate,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage(compositionFeedBackSystemPrompt),
			schema.UserMessage("作文题目:{{.title}}\n\n 作文内容:{{.composition_content}}"),
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}
