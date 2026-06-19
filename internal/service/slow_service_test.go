package service

import "testing"

func BenchmarkSlowHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        SlowHash("test")
    }
}

func BenchmarkInefficientConcat(b *testing.B) {
    items := make([]string, 1000)
    for i := range items {
        items[i] = "item"
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        InefficientConcat(items)
    }
}

func BenchmarkEfficientConcat(b *testing.B) {
    items := make([]string, 1000)
    for i := range items {
        items[i] = "item"
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        EfficientConcat(items)
    }
}
