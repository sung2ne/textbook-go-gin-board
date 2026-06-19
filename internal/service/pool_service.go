package service

import (
    "bytes"
    "sync"
)

var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

// GetBuffer - 풀에서 버퍼 가져오기
func GetBuffer() *bytes.Buffer {
    return bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer - 풀에 버퍼 반환
func PutBuffer(buf *bytes.Buffer) {
    buf.Reset()
    bufferPool.Put(buf)
}

// 사용 예시
func ProcessWithPool(data []byte) string {
    buf := GetBuffer()
    defer PutBuffer(buf)

    buf.Write(data)
    // 처리...

    return buf.String()
}
