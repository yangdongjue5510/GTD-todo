# CLAUDE.md - 백엔드

이 파일은 GTD-todo Go 애플리케이션의 백엔드 관련 가이드를 제공합니다.

## 아키텍처

### Handler-Service-Model 패턴
```
HTTP Request → Handler → Service → Model
```

- **Handlers** (`*_handler.go`): HTTP 요청/응답 처리, JSON 바인딩, 상태 코드
- **Services** (`*_service.go`): 비즈니스 로직, 데이터 검증, 워크플로우 구현
- **Models** (`*_model.go`): 데이터 구조, 열거형, 도메인 엔티티

### 현재 패키지
- `thing/`: Pending/Someday/Done 상태를 가진 Thing (생각) 관리
- `action/`: 더 세분화된 상태 추적을 가진 Action (작업) 관리
- 메인 패키지: 서버 설정, 라우팅, CORS 구성

## 개발 명령어

```bash
# 개발 서버 실행
go run main.go

# 바이너리 빌드
go build -o gtd-todo

# 모든 테스트 실행
go test ./...

# 상세한 테스트 출력
go test -v ./...

# 의존성 업데이트
go mod tidy
```

## 코딩 컨벤션

### 오류 처리
```go
if err := c.ShouldBindJSON(&thing); err != nil {
    c.JSON(400, gin.H{"error": "Invalid input"})
    return
}
```

### 서비스 인터페이스 패턴
```go
type ThingService interface {
    AddThing(thing Thing)
    GetThings() []Thing
    Clarify(thing Thing)
}
```

### 상태 열거형
String() 메서드와 함께 `iota`를 사용한 상태 열거형:
```go
type Status int

const (
    Pending Status = iota
    Someday
    Done
)

func (s Status) String() string { ... }
```

### JSON 구조체 태그
API 모델에는 항상 JSON 태그를 포함:
```go
type Thing struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      Status `json:"status"`
}
```

## API 설계

### REST 컨벤션
- `POST /things/` - 새 Thing 생성
- `GET /things/` - 모든 Thing 목록 조회
- 적절한 HTTP 상태 코드 사용 (200, 201, 400, 404, 500)
- 일관된 구조의 JSON 응답 반환

### CORS 구성
현재 개발을 위해 모든 origin을 허용. 운영환경에서는 업데이트 필요:
```go
AllowOrigins: []string{"*"} // 운영환경에서는 변경
```

## 테스팅

### 테스트 파일 명명
- 같은 패키지에 `*_test.go` 파일
- 테스트 함수 이름: `TestFunctionName`
- 여러 시나리오에는 테이블 기반 테스트 사용

### 테스트 구조 예시
```go
func TestAddThing(t *testing.T) {
    service := &InmemoryThingService{}
    thing := Thing{Title: "Test", Description: "Test desc"}
    
    service.AddThing(thing)
    
    things := service.GetThings()
    assert.Equal(t, 1, len(things))
}
```

## 현재 제한사항 및 TODO

1. **인메모리 저장소**: 재시작 시 데이터 손실
2. **Action 엔드포인트 누락**: Thing API만 구현됨
3. **인증 없음**: 개방형 API
4. **Clarify 메서드**: ThingService에서 구현되지 않음
5. **Project 엔티티**: 정의되었지만 완전히 구현되지 않음

## GTD 구현 노트

### 핵심 워크플로우
1. **포착**: POST /things/를 통한 Thing 생성
2. **명확화**: Thing을 Action으로 변환 (구현 예정)
3. **정리**: 상태 기반 필터링 및 Project 그룹화
4. **검토**: 오래되거나 정체된 항목의 주기적 정리

### 상태 전환
- **Thing**: Pending → Someday/Done/Removed
- **Action**: 위임, 지연, 계획 상태를 포함한 더 복잡한 워크플로우