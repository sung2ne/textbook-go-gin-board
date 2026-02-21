package service

import "testing"

func BenchmarkSliceGrowth(b *testing.B) {
    for i := 0; i < b.N; i++ {
        SliceGrowth()
    }
}

func BenchmarkPreallocatedSlice(b *testing.B) {
    for i := 0; i < b.N; i++ {
        PreallocatedSlice()
    }
}
