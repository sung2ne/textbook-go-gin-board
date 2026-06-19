package testutil

import (
    "encoding/json"
    "testing"

    "github.com/stretchr/testify/require"
)

func ParseJSON[T any](t *testing.T, data []byte) T {
    var result T
    err := json.Unmarshal(data, &result)
    require.NoError(t, err)
    return result
}

// 사용 예시
func TestGetPosts(t *testing.T) {
    // ... 요청 실행 ...

    posts := testutil.ParseJSON[[]PostResponse](t, rec.Body.Bytes())
    assert.Len(t, posts, 10)
}
