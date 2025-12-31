package main

// Composition 表示一篇小学生作文
type Composition struct {
	// Title 是作文的标题，例如 "缅怀先烈 逐梦未来"
	Title string `json:"title"`

	// Content 是作文的正文内容，保留原始段落、换行、拼音、错别字等
	// 使用 string 而非 []string，以完整保留原始格式（包括 \n 等换行符）
	Content string `json:"content"`
}

// CompositionFeedback 表示一篇学生作文的完整批改反馈结果
type CompositionFeedback struct {
	// Score 作文总分，取值范围为 0-100
	Score int `json:"score"`

	// Strengths 作文的主要优点，通常包含 1-3 条
	Strengths []string `json:"strengths"`

	// AreasForImprovement 需要改进的方面，通常包含 1-3 条建设性建议
	AreasForImprovement []string `json:"areas_for_improvement"`

	// Corrections 包含具体的错别字、语法和语言优化建议
	Corrections Corrections `json:"corrections"`

	// OverallComment 整体评语，一段鼓励性、总结性的文字
	OverallComment string `json:"overall_comment"`
}

// Corrections 包含三类具体的修改建议：错别字、语法错误、语言润色
type Corrections struct {
	// Typos 错别字列表，每项包含原文、修正及上下文
	Typos []TypoCorrection `json:"typos"`

	// Grammar 语法或用词不当的修正建议
	Grammar []GrammarCorrection `json:"grammar"`

	// Enhancements 语言表达优化建议，提升文采与感染力
	Enhancements []EnhancementSuggestion `json:"enhancements"`
}

// TypoCorrection 表示一个错别字修正项
type TypoCorrection struct {
	// Original 出现的错别字
	Original string `json:"original"`

	// Corrected 正确的写法
	Corrected string `json:"corrected"`

	// Context 包含该错字的短句上下文（用于定位）
	Context string `json:"context"`
}

// GrammarCorrection 表示一个语法或表达不当的修正项
type GrammarCorrection struct {
	// OriginalSentence 原始存在问题的句子
	OriginalSentence string `json:"original_sentence"`

	// CorrectedSentence 修正后的句子
	CorrectedSentence string `json:"corrected_sentence"`

	// Reason 修改原因，简要说明问题所在
	Reason string `json:"reason"`
}

// EnhancementSuggestion 表示一个语言润色或文采提升建议
type EnhancementSuggestion struct {
	// OriginalSentence 原始句子（通常语法正确但表达平实）
	OriginalSentence string `json:"original_sentence"`

	// EnhancedSentence 优化后的建议句式
	EnhancedSentence string `json:"enhanced_sentence"`

	// Comment 对优化点的简要点评，说明提升效果
	Comment string `json:"comment"`
}
