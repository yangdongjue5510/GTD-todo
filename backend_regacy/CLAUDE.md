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

### Domain Driven Design
backend 디렉토리 하위는 도메인 별로 디렉토리가 나뉘고 각 디렉토리의 `CLAUDE.md` 파일에 규정된 도메인 로직을 기반으로 해당 디렉토리에 구현.
각 도메인에 해당하는 API, 상태 관련 스키마 등은 해당 도메인 디렉토리의 CLAUDE.md에 정해놓음.

#### 서브 도메인 구조
Domain Driven Design에 의해 서브 도메인은 아래와 설계 되었음.

**Workflow Context (Core Domain)**

- Action과 Project의 생명 주기를 관리
- GTD의 organize, engage에 따라 Action과 Project의 생성 및 상태 변경을 처리.

**Capture Context (Support Domain)**

- 외부의 다양한 데이터 혹은 입력을 기반으로 inbox에 Thing으로 변환하여 모으는 기능 관리
- Thing을 Action으로 변환하는 clarify도 지원

**Reflect Context (Support Domain)**

- 일정한 주기(일간, 주간)으로 Thing, Action, Project 진행 현황을 정리해서 사용자에게 보고하는 기능

**Notification Context (Generic Domain)**

- 사용자에게 알림을 보내는 기능을 다루는 영역

**Authorization Context (Generic Domain)**

- 사용자 인증/인가 관련된 기능을 다루는 영역.

> 주요 도메인의 의존성 방향은 아래와 같음.

Capture → Workflow  
Reflect → Capture + Workflow

이때 서로 다른 도메인의 모델을 직접 참조해서 사용하지는 않도록 하고 서비스만 호출하도록 한다.

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


## API 설계
모든 도메인 디렉토리의 HTTP 핸들러는 아래와 같은 방식에 따라 API를 설계해야 함.
### REST 컨벤션
- RESTful API에 맞는 HTTP Method 사용.
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

## GTD 구현 노트

### 핵심 워크플로우
1. **포착**: POST /things/를 통한 Thing 생성
2. **명확화**: Thing을 Action으로 변환 (구현 예정)
3. **정리**: 상태 기반 필터링 및 Project 그룹화
4. **검토**: 오래되거나 정체된 항목의 주기적 정리

### 상태 전환
- **Thing**: Pending → Someday/Done/Removed
- **Action**: 위임, 지연, 계획 상태를 포함한 더 복잡한 워크플로우