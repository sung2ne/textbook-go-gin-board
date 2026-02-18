package batch

import (
	"context"
)

// Chunk 슬라이스를 지정 크기로 분할
func Chunk[T any](items []T, size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[i:end])
	}
	return chunks
}

// ProcessChunks 청크 단위 병렬 처리
func ProcessChunks[T any](
	ctx context.Context,
	items []T,
	chunkSize int,
	numWorkers int,
	handler func(context.Context, []T) error,
) error {
	chunks := Chunk(items, chunkSize)

	processor := NewProcessor(numWorkers, handler)
	return processor.Process(ctx, chunks)
}
