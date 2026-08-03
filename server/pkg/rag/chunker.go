package rag

import (
	"regexp"
	"strings"
	"unicode"
)

// Chunk 知识库文本块
type Chunk struct {
	CourseID     string
	CourseTitle  string
	ChapterID    string
	ChapterTitle string
	Heading      string
	Text         string
}

// Section Markdown 切分后的段落
type Section struct {
	Heading string
	Body    string
}

type mdSection struct {
	heading, body string
}

// SplitByHeadings 按 Markdown 标题切分章节内容
func SplitByHeadings(markdown string) []Section {
	lines := strings.Split(markdown, "\n")
	var sections []mdSection
	current := mdSection{}
	inCode := false
	headingRe := regexp.MustCompile(`^#{1,3}\s`)
	stripRe := regexp.MustCompile(`^#+\s*`)
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCode = !inCode
		}
		if !inCode && headingRe.MatchString(line) {
			if strings.TrimSpace(current.body) != "" {
				sections = append(sections, current)
			}
			current = mdSection{heading: stripRe.ReplaceAllString(line, ""), body: ""}
		} else {
			current.body += line + "\n"
		}
	}
	if strings.TrimSpace(current.body) != "" {
		sections = append(sections, current)
	}
	out := make([]Section, len(sections))
	for i, s := range sections {
		out[i] = Section{Heading: s.heading, Body: s.body}
	}
	return out
}

// StripMarkdown 去除 Markdown 格式，保留纯文本
func StripMarkdown(text string) string {
	re := regexp.MustCompile("(?s)```.*?```")
	text = re.ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\|.*\|`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`[#>*` + "`" + `_[\]()!-]`).ReplaceAllString(text, " ")
	return strings.Join(strings.Fields(text), " ")
}

// Tokenize 中文滑窗 + 英文分词，用于关键词检索
func Tokenize(query string) map[string]bool {
	tokens := map[string]bool{}
	lower := strings.ToLower(query)
	for _, m := range regexp.MustCompile(`[a-z0-9-]{2,}`).FindAllString(lower, -1) {
		tokens[m] = true
	}
	runes := []rune(query)
	for i := 0; i < len(runes); i++ {
		for l := 2; l <= 4 && i+l <= len(runes); l++ {
			sub := string(runes[i : i+l])
			if isMostlyChinese(sub) {
				tokens[sub] = true
			}
		}
	}
	stop := map[string]bool{"什么": true, "怎么": true, "如何": true, "的": true, "了": true, "是": true}
	for k := range tokens {
		if stop[k] {
			delete(tokens, k)
		}
	}
	return tokens
}

func isMostlyChinese(s string) bool {
	n := 0
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			n++
		}
	}
	return n >= len(s)/2
}

// KeywordScore 计算查询与文本的关键词匹配分数
func KeywordScore(text, heading, query string) float64 {
	tokens := Tokenize(query)
	if len(tokens) == 0 {
		return 0
	}
	lower := strings.ToLower(text + " " + heading)
	score := 0.0
	for tok := range tokens {
		if strings.Contains(lower, tok) {
			score += 1
		}
	}
	return score / float64(len(tokens))
}
