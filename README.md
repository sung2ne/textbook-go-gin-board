# 소설처럼 읽는 Go 언어 - 실습 코드

[위키독스 교재](https://wikidocs.net/book/18987)의 실습 코드 저장소입니다.

Go + Gin 프레임워크로 게시판 REST API를 단계별로 구현합니다.

## 사용 방법

원하는 챕터의 브랜치를 체크아웃하면 해당 시점까지의 완성된 프로젝트를 받을 수 있습니다.

```bash
# 저장소 클론
git clone https://github.com/sung2ne/textbook-go-gin-board.git
cd textbook-go-gin-board

# 원하는 챕터로 이동
git checkout part04/chapter-02   # PART 04의 02장까지 완성된 코드
```

## 브랜치 목록

각 브랜치는 해당 챕터까지의 코드가 누적 적용되어 독립적으로 빌드 가능합니다.

### PART 04. 게시판 만들기

| 브랜치 | 내용 |
|--------|------|
| `part04/chapter-01` | 게시판 설계 |
| `part04/chapter-02` | 게시글 CRUD |
| `part04/chapter-03` | 페이징과 검색 |
| `part04/chapter-04` | 댓글 기능 |

### PART 05. 미들웨어

| 브랜치 | 내용 |
|--------|------|
| `part05/chapter-01` | 미들웨어 기초 (RequestID, Timing) |
| `part05/chapter-02` | 로깅과 모니터링 (zerolog, Prometheus) |
| `part05/chapter-03` | 에러 처리 (AppError, RFC 7807, Recovery) |
| `part05/chapter-04` | CORS와 보안 (Rate Limit, Secure Headers) |

### PART 06. JWT 인증

| 브랜치 | 내용 |
|--------|------|
| `part06/chapter-01` | JWT 기초 (Claims, TokenService) |
| `part06/chapter-02` | 사용자 인증 (User, Password, AuthService) |
| `part06/chapter-03` | 인증 미들웨어 (Auth, Role, Owner) |
| `part06/chapter-04` | 토큰 관리 (Redis TokenStore, Blacklist) |

### PART 07. 인증된 게시판

| 브랜치 | 내용 |
|--------|------|
| `part07/chapter-01` | 게시판 인증 연동 |
| `part07/chapter-02` | 댓글 알림과 멘션 |
| `part07/chapter-03` | 권한 관리 (RBAC, Admin) |
| `part07/chapter-04` | 마이페이지 (Profile, Password, Withdraw) |

### PART 08. 동시성 프로그래밍

| 브랜치 | 내용 |
|--------|------|
| `part08/chapter-01` | 고루틴 기초 |
| `part08/chapter-02` | 채널 |
| `part08/chapter-03` | 동시성 패턴 (WorkerPool, Retry, Backoff) |
| `part08/chapter-04` | 웹 서버 동시성 (Worker Queue, Email, WebSocket, Scheduler, Batch) |

### PART 09. API 문서화

| 브랜치 | 내용 |
|--------|------|
| `part09/chapter-01` | Swagger 설정 |
| `part09/chapter-02` | API 문서 작성 (핸들러 어노테이션) |
| `part09/chapter-03` | 문서 고급 기능 (Security, Tags) |

## 기술 스택

- **언어**: Go 1.21+
- **웹 프레임워크**: Gin
- **ORM**: GORM
- **데이터베이스**: MySQL, Redis
- **인증**: JWT (golang-jwt)
- **API 문서**: Swagger (swaggo/swag)
- **로깅**: zerolog
- **모니터링**: Prometheus

## 실행 방법

```bash
# 의존성 설치
go mod tidy

# 설정 파일 확인
# config/config.yaml에서 DB/Redis 연결 정보 설정

# 서버 실행
go run cmd/api/main.go

# Swagger UI (PART 09 이후)
# http://localhost:8080/swagger/index.html
```

## 프로젝트 구조

```
├── cmd/api/main.go          # 엔트리포인트
├── config/                  # 설정 파일
├── internal/
│   ├── auth/                # JWT, 비밀번호 서비스
│   ├── config/              # 설정 로더
│   ├── database/            # DB 연결 (MySQL, Redis)
│   ├── domain/              # 도메인 모델
│   ├── dto/                 # 요청/응답 DTO
│   ├── handler/             # HTTP 핸들러
│   ├── middleware/           # 미들웨어
│   ├── repository/          # 데이터 접근 계층
│   ├── router/              # 라우팅
│   └── service/             # 비즈니스 로직
├── pkg/
│   ├── pool/                # 워커 풀 (PART 08)
│   ├── retry/               # 재시도 패턴 (PART 08)
│   └── batch/               # 배치 처리 (PART 08)
└── docs/                    # Swagger 문서 (PART 09)
```

## 라이선스

이 저장소는 [소설처럼 읽는 Go 언어](https://wikidocs.net/book/18987) 교재의 실습 코드입니다.
