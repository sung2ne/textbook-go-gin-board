package util

import "testing"

func BenchmarkSlugify(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Slugify("Hello World 123!")
    }
}
