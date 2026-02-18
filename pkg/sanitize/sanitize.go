package sanitize

import (
	"html"
	"regexp"
	"strings"
)

func HTML(s string) string {
	return html.EscapeString(s)
}

func StripTags(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func Normalize(s string) string {
	s = TrimSpace(s)
	s = StripTags(s)
	return s
}
