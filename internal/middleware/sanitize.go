package middleware

var sensitiveFields = map[string]bool{
    "password":     true,
    "token":        true,
    "secret":       true,
    "credit_card":  true,
    "card_number":  true,
}

func SanitizeForLog(data map[string]any) map[string]any {
    result := make(map[string]any)
    for k, v := range data {
        if sensitiveFields[k] {
            result[k] = "[REDACTED]"
        } else {
            result[k] = v
        }
    }
    return result
}
