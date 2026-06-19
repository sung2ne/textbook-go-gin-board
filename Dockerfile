FROM golang:1.22-alpine AS builder

WORKDIR /app

# 의존성 캐싱
COPY go.mod go.sum ./
RUN go mod download

# 소스 복사 및 빌드
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o server ./cmd/server

# 최소 런타임
FROM alpine:3.19

# 보안: 루트가 아닌 사용자로 실행
RUN adduser -D -g '' appuser

WORKDIR /app

# 타임존 (한국 시간 사용 시)
RUN apk --no-cache add tzdata

COPY --from=builder /app/server .

USER appuser

EXPOSE 8080

CMD ["./server"]
