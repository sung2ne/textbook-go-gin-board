package util

import "testing"

func TestSlugify(t *testing.T) {
    t.Run("basic string", func(t *testing.T) {
        got := Slugify("Hello World")
        want := "hello-world"
        if got != want {
            t.Errorf("Slugify(\"Hello World\") = %q, want %q", got, want)
        }
    })

    t.Run("special characters", func(t *testing.T) {
        got := Slugify("Hello! World?")
        want := "hello-world"
        if got != want {
            t.Errorf("Slugify(\"Hello! World?\") = %q, want %q", got, want)
        }
    })

    t.Run("numbers", func(t *testing.T) {
        got := Slugify("Go 1.21 Release")
        want := "go-1-21-release"
        if got != want {
            t.Errorf("Slugify(\"Go 1.21 Release\") = %q, want %q", got, want)
        }
    })
}
