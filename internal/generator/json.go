package generator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"gendata/internal/generator/data"
)

var placeholderRegex = regexp.MustCompile(`\{([^{}]+)\}`)

func (g *Random) JSONString() string {
	length := g.IntN(10) + 1
	var emojidate []data.Emoji
	for i := 0; i < length; i++ {
		emojidate = append(emojidate, g.Emoji())
	}
	bytes, _ := json.MarshalIndent(emojidate, "", "  ")
	return string(bytes)
}

func (g *Random) MarkdownString() string {
	var sb strings.Builder
	sb.WriteString("| Emoji | Description | Category | Aliases | Tags | Unicode Version | iOS Version | Sentences |\n")
	sb.WriteString("|-------|-------------|----------|---------|------|-----------------|-------------|-----------|\n")

	length := g.IntN(10) + 1
	for i := 0; i < length; i++ {
		e := g.Emoji()
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			e.Emoji,
			e.Description,
			e.Category,
			strings.Join(e.Aliases, ", "),
			strings.Join(e.Tags, ", "),
			e.UnicodeVersion,
			e.IOSVersion,
			e.SentencesEmoji,
		))
	}
	return sb.String()
}

func (g *Random) Emoji() data.Emoji {
	a := data.Emojis[g.IntN(len(data.Emojis))]
	b := data.SentencesEmoji[g.IntN(len(data.SentencesEmoji))]

	a.SentencesEmoji = placeholderRegex.ReplaceAllString(b, a.Emoji)

	return a
}
