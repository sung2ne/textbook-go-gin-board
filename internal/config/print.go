package config

import (
    "fmt"
    "strings"
)

func (c *Config) Print() {
    fmt.Println("=== Configuration ===")
    fmt.Printf("App Name: %s\n", c.App.Name)
    fmt.Printf("Environment: %s\n", c.App.Environment)
    fmt.Printf("Debug: %v\n", c.App.Debug)
    fmt.Printf("Port: %d\n", c.Server.Port)
    fmt.Printf("Database URL: %s\n", maskDatabaseURL(c.Database.URL))
    fmt.Printf("JWT Secret: %s\n", maskSecret(c.JWT.Secret))
    fmt.Printf("Log Level: %s\n", c.Log.Level)
    fmt.Println("====================")
}

func maskDatabaseURL(url string) string {
    // postgres://user:password@host:port/db -> postgres://user:****@host:port/db
    if idx := strings.Index(url, "@"); idx > 0 {
        prefix := url[:strings.Index(url, ":")+3] // postgres://
        suffix := url[idx:]                        // @host:port/db
        return prefix + "****:****" + suffix
    }
    return url
}

func maskSecret(secret string) string {
    if len(secret) <= 8 {
        return "****"
    }
    return secret[:4] + "****" + secret[len(secret)-4:]
}
