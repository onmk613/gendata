package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"gendata/internal/generator/data"
)

var placeholderRegex = regexp.MustCompile(`\{([^{}]+)\}`)

func JSONString() string {
	length := rand.Intn(10) + 1
	var emojidate []data.Emoji
	for i := 0; i < length; i++ {
		emojidate = append(emojidate, Emoji())
	}
	bytes, _ := json.MarshalIndent(emojidate, "", "  ")
	return string(bytes)
}

func MarkdownString() string {
	var sb strings.Builder
	sb.WriteString("| Emoji | Description | Category | Aliases | Tags | Unicode Version | iOS Version | Sentences |\n")
	sb.WriteString("|-------|-------------|----------|---------|------|-----------------|-------------|-----------|\n")

	length := rand.Intn(10) + 1
	for i := 0; i < length; i++ {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			Emoji().Emoji,
			Emoji().Description,
			Emoji().Category,
			strings.Join(Emoji().Aliases, ", "),
			strings.Join(Emoji().Tags, ", "),
			Emoji().UnicodeVersion,
			Emoji().IOSVersion,
			Emoji().SentencesEmoji,
		))
	}
	return sb.String()
}

func Emoji() data.Emoji {
	a := data.Emojis[rand.Intn(len(data.Emojis))]
	b := data.SentencesEmoji[rand.Intn(len(data.SentencesEmoji))]

	a.SentencesEmoji = placeholderRegex.ReplaceAllString(b, a.Emoji)

	return a
}

// func init() {
// 	addEmojiLookup()
// }

// type GenerateFunc func(params ...string) string

// var funcRegistry = map[string]GenerateFunc{}

// func RegisterFunc(name string, fn GenerateFunc) {
// 	funcRegistry[name] = fn
// }

// func ProcessTemplate(template string) string {
// 	var placeholderRegex = regexp.MustCompile(`\{([^{}]+)\}`)
// 	result := placeholderRegex.ReplaceAllStringFunc(template, func(match string) string {
// 		var params []string

// 		// 去掉花括号后的字符串分为两部分，第一部分是函数名，第二部分是参数列表（如果有的话）
// 		parts := strings.SplitN(match[1 : len(match)-1], ":", 2)
// 		name := parts[0]

// 		if len(parts) > 1 {
// 			params = strings.Split(parts[1], ",")
// 		}

// 		// 查找函数并调用
// 		fn, ok := funcRegistry[name]
// 		if !ok {
// 			return match
// 		}

// 		return fn(params...)
// 	})

// 	return result
// }

// func addEmojiLookup() {
// 	RegisterFunc("emoji", func(params ...string) string {
// 		a := rand.Intn(len(data.Emojis))
// 		return data.Emojis[a].Emoji
// 	})
// }
