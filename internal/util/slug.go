package util

import (
    "regexp"
    "strings"
)

var (
    nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)
)

// Slugify는 문자열을 URL 친화적인 슬러그로 변환합니다.
func Slugify(s string) string {
    s = strings.ToLower(s)
    s = nonAlphanumeric.ReplaceAllString(s, "-")
    s = strings.Trim(s, "-")
    return s
}
