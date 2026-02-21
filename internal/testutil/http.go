package testutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// ParseJSON은 JSON 바이트를 타입 T로 파싱합니다.
func ParseJSON[T any](t *testing.T, data []byte) T {
	var result T
	err := json.Unmarshal(data, &result)
	require.NoError(t, err)
	return result
}
