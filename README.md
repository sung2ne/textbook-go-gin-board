# 소설처럼 읽는 Go 언어 - 실습 코드 저장소

> **교재**: [소설처럼 읽는 Go 언어](https://wikidocs.net/book/18987) (위키독스)
> **버전**: v1.3.1

Go와 Gin 프레임워크로 게시판 API를 단계별로 만들어가는 교재의 실습 코드 저장소입니다.

## 기술 스택

| 분류 | 기술 |
|------|------|
| 언어 | Go 1.21+ |
| 웹 프레임워크 | Gin |
| ORM | GORM |
| 데이터베이스 | PostgreSQL |
| 캐시 | Redis |
| 인증 | JWT |
| 문서화 | Swagger (swaggo) |

## 시작하기

### 저장소 클론

```bash
git clone https://github.com/sung2ne/textbook-go-gin-board.git
cd textbook-go-gin-board
```

### 원하는 챕터로 이동

각 브랜치에는 해당 챕터까지의 코드가 누적 적용되어 있습니다.

```bash
# PART 04의 02장까지 완성된 코드
git checkout part04/chapter-02

# PART 07의 04장까지 완성된 코드
git checkout part07/chapter-04
```

### 프로젝트 실행

```bash
# 의존성 설치
go mod tidy

# 설정 파일 확인 (config/config.yaml)
# DB, Redis 연결 정보를 환경에 맞게 수정

# 서버 실행
go run cmd/api/main.go

# Swagger UI (PART 09 이후)
# http://localhost:8080/swagger/index.html
```

### Docker로 실행

```bash
# 개발 환경 (PostgreSQL + Redis 포함)
docker-compose up -d
```

---

## 브랜치 목록

> PART 01~03은 독립적인 예제 코드 위주라 별도 브랜치가 없습니다. PART 04부터 게시판 프로젝트가 시작됩니다.

총 **35개** 브랜치가 PART별로 제공됩니다.

### PART 04. 비인증 게시판

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

### PART 10. 테스트

| 브랜치 | 내용 |
|--------|------|
| `part10/chapter-01` | 단위 테스트 (Mock, testify) |
| `part10/chapter-02` | 통합 테스트 (testutil, scenarios) |
| `part10/chapter-03` | 테스트 자동화 (벤치마크) |

### PART 11. 배포

| 브랜치 | 내용 |
|--------|------|
| `part11/chapter-01` | Docker (Dockerfile, docker-compose) |
| `part11/chapter-02` | 프로덕션 설정 (환경 변수, 로거, 워커) |
| `part11/chapter-03` | 클라우드 배포 (Cloud Run, Fly.io, GitHub Actions) |

### PART 12. 성능 최적화

| 브랜치 | 내용 |
|--------|------|
| `part12/chapter-01` | 프로파일링 (pprof, trace) |
| `part12/chapter-02` | 캐싱 (인메모리, Redis, HTTP 캐싱) |
| `part12/chapter-03` | 데이터베이스 최적화 (인덱스, 쿼리, 레플리카) |

### PART 13. 프로젝트 마무리

| 브랜치 | 내용 |
|--------|------|
| `part13/chapter-01` | 코드 품질 (lint, 리뷰) |
| `part13/chapter-02` | 보안 점검 |
| `part13/chapter-03` | 다음 단계 |

---

## 프로젝트 구조

```
goboardapi/
├── cmd/api/          # 애플리케이션 진입점
├── config/           # 설정 파일
├── internal/
│   ├── domain/       # 도메인 모델 (Entity)
│   ├── dto/          # 데이터 전송 객체
│   ├── repository/   # 데이터 접근 계층
│   ├── service/      # 비즈니스 로직
│   ├── handler/      # HTTP 핸들러
│   └── middleware/   # Gin 미들웨어
└── pkg/              # 공통 유틸리티
```

## 활용 팁

**교재를 따라가며 직접 코딩하는 것을 추천합니다.** 저장소의 코드는 다음과 같은 상황에서 활용하세요.

- 코드가 정상 동작하지 않을 때 비교 대상으로 활용
- 특정 챕터부터 학습을 시작하고 싶을 때 해당 브랜치에서 출발
- 전체 프로젝트 구조를 한눈에 파악하고 싶을 때 참고

## 업데이트 이력

| 날짜 | 버전 | 내용 |
|------|------|------|
| 2026-02-22 | v1.3.1 | 버그 수정: PART 08 select 문 예제 데드락 수정 (독자 문의 반영) |
| 2026-02-21 | v1.3.0 | 기술 감사 반영: 모듈명 통일 (17개 파일), 편집 노트 3건, 시각자료 12건 변환 |
| 2026-02-21 | v1.2.0 | PART 10~13 챕터별 브랜치 배포 (12개 추가, 총 36개) |
| 2026-02-18 | v1.1.0 | GitHub 저장소 공개, PART 04~09 챕터별 브랜치 배포 (24개) |

## 라이선스

이 저장소의 코드는 교육 목적으로 자유롭게 사용할 수 있습니다.
