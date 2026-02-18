package util

import (
	"regexp"
	"strings"
)

// 멘션 패턴: @로 시작하고 한글/영문/숫자로 구성
var mentionPattern = regexp.MustCompile(`@([가-힣a-zA-Z0-9_]+)`)

// ParseMentions 댓글 내용에서 멘션된 사용자명 추출
func ParseMentions(content string) []string {
	matches := mentionPattern.FindAllStringSubmatch(content, -1)

	// 중복 제거를 위해 맵 사용
	seen := make(map[string]bool)
	var usernames []string

	for _, match := range matches {
		if len(match) > 1 {
			username := match[1]
			if !seen[username] {
				seen[username] = true
				usernames = append(usernames, username)
			}
		}
	}

	return usernames
}

// HighlightMentions 멘션을 HTML 링크로 변환
func HighlightMentions(content string) string {
	return mentionPattern.ReplaceAllStringFunc(content, func(match string) string {
		username := strings.TrimPrefix(match, "@")
		return `<a href="/users/` + username + `" class="mention">` + match + `</a>`
	})
}
