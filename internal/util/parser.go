package util

import (
    "errors"
    "strconv"
    "strings"
)

// ParsePagination은 page와 size 문자열을 파싱합니다.
func ParsePagination(pageStr, sizeStr string) (page, size int, err error) {
    page = 1   // 기본값
    size = 10  // 기본값

    if pageStr != "" {
        page, err = strconv.Atoi(strings.TrimSpace(pageStr))
        if err != nil || page < 1 {
            return 0, 0, errors.New("invalid page")
        }
    }

    if sizeStr != "" {
        size, err = strconv.Atoi(strings.TrimSpace(sizeStr))
        if err != nil || size < 1 || size > 100 {
            return 0, 0, errors.New("invalid size")
        }
    }

    return page, size, nil
}
