package logger

import (
	"regexp"
)

var sensitiveFields = []string{
	"password",
	"token",
	"secret",
	"api_key",
	"credit_card",
}

var sensitivePattern = regexp.MustCompile(
	`"(password|token|secret|api_key|credit_card)"\s*:\s*"[^"]*"`,
)

func SanitizeJSON(data []byte) []byte {
	return sensitivePattern.ReplaceAll(data, []byte(`"$1":"[REDACTED]"`))
}

func SanitizeMap(m map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range m {
		if isSensitive(k) {
			result[k] = "[REDACTED]"
		} else if nested, ok := v.(map[string]any); ok {
			result[k] = SanitizeMap(nested)
		} else {
			result[k] = v
		}
	}
	return result
}

func isSensitive(field string) bool {
	for _, s := range sensitiveFields {
		if field == s {
			return true
		}
	}
	return false
}
