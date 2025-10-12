# CLAUDE.md

이 파일은 Claude Code (claude.ai/code)가 이 저장소에서 코드 작업을 할 때 참고할 가이드를 제공합니다.
클로드 코드로 코드 베이스에 작업이 진행되고 나면 반드시 해당 작업에 해당하는 CLAUDE.md 파일에 해당 내용을 최신화해야 합니다.

## 프로젝트 개요

GTD-todo는 Getting Things Done (GTD) 방법론을 기반으로 한 웹 서비스입니다. 사용자가 부담 없이 생각을 포착하고, 이를 실행 가능한 항목으로 명확화하여 관리 가능한 시스템으로 구성하는 것을 돕습니다.

## 아키텍처

GTD-todo는 다음과 같은 구조로 구성됩니다:

- **backend/**: Go 기반 REST API 서버 (포트 8080)
- **frontend/**: 순수 HTML/CSS/JavaScript 클라이언트
- **docs/**: 프로젝트 설계 문서 및 가이드

## 핵심 도메인 개념

### GTD 워크플로우
```
포착(Capture) → 명확화(Clarify) → 칸반 관리(Organize) → 실행(Execute) → 검토(Review)
```

### 전체 시스템 구조
```
[Inbox] → Clarify → [Kanban Board] + [Projects] + [Routines]
```

### 주요 엔티티

#### 1. Thing (Inbox)
사용자가 자유롭게 기입하는 구조화되지 않은 생각/아이디어

**상태:**
- 대기: 생성되어 분류되기를 대기
- 처리됨: Clarify 과정을 통해 Action/Project/Routine으로 변환됨

#### 2. Action (칸반 보드)
Thing에서 도출된 구체적이고 실행 가능한 작업을 칸반 형태로 관리

**칸반 컬럼 (5개 상태):**
- **Delay**: 당장 하지 않아도 나중에 고려할 일
- **Next Action**: 당장 진행 가능한 구체적 행동  
- **In Progress**: 현재 수행 중인 작업
- **Waiting For**: 다른 사람 위임/특정 응답 대기
- **Done**: 완료된 작업

#### 3. Project
여러 Action으로 세분화되는 작업 집합
- 프로젝트 하위 Action들은 칸반에 프로젝트 라벨 표시
- 프로젝트별 진행률 추적

#### 4. Routine  
주기적으로 반복되는 작업 관리
- 주기마다 Action 생성하여 칸반에 루틴 라벨 표시
- 습관 형성 지원

### 핵심 프로세스

#### Capture
사용자가 생각을 자유롭게 Inbox에 포착하는 과정

#### Clarify (Thing → Action 변환)
Thing을 드래그하여 6가지 선택지 중 선택:

1. **지금 당장하기!** → In Progress 상태로 Action 생성
2. **나중에 하기** → Delay 상태로 Action 생성  
3. **할일로 기록하기** → Next Action 상태로 Action 생성
4. **달력에 기록하기** → Calendar 상태로 Action 생성
5. **프로젝트 시작하기!** → Project 생성
6. **루틴 생성하기** → Routine 생성

#### Organize & Execute
칸반 보드를 통한 시각적 작업 관리
- 드래그 드롭으로 상태 변경
- 프로젝트/루틴 라벨을 통한 분류

#### Review
일간, 주간, 월간 리뷰를 통한 전체 현황 점검

## 개발 환경 설정

### 전체 애플리케이션 실행
```bash
# 백엔드 서버 시작 (터미널 1)
cd backend && go run main.go

# 프론트엔드 서빙 (터미널 2)
cd frontend && python -m http.server 3000
```

브라우저에서 `http://localhost:3000`으로 접속

## 디렉토리별 가이드

각 디렉토리에는 해당 영역에 특화된 `CLAUDE.md` 파일이 있습니다:

- `backend/CLAUDE.md`: Go 백엔드 개발 가이드
- `frontend/CLAUDE.md`: 프론트엔드 개발 가이드

## 중요한 설계 원칙

프로젝트 문서에서:
- 사용자 입력 과정이 절대 부담스럽게 느껴져서는 안 됩니다
- 명확화 과정은 가벼우며 방해가 되지 않아야 합니다
- 큰 작업을 분해하기 위한 GTD 프로젝트 개념을 지원합니다
- 명확화 과정에서 항목을 지연하거나 제거할 수 있습니다

## 프로젝트 현황

### 구현 완료
- ✅ Thing 포착 및 목록 조회 API
- ✅ 기본 인박스 UI
- ✅ 백엔드-프론트엔드 연동

### 구현 예정 (GTD 핵심 기능)
- ❌ Thing → Action 드래그 Clarify 인터페이스 (6가지 선택지)
- ❌ Action 칸반 보드 (Delay/Next Action/In Progress/Waiting For/Done)
- ❌ Project 생성 및 관리 (하위 Action 라벨링)
- ❌ Routine 생성 및 주기적 Action 생성
- ❌ 검토(Review) 대시보드
- ❌ 데이터 영속성 (현재 인메모리)

### 현재 구현 범위
- ✅ 포함: Thing, Action 칸반, Clarify, Project, Routine
- ❌ 제외: GTD Context(@Calls 등), Archive

## 개발 시 주의사항

- 각 기능은 GTD 철학에 맞게 "부담 없이" 사용할 수 있어야 함
- UI는 최소한의 마찰로 빠른 입력이 가능해야 함
- 상세한 구현 가이드는 각 디렉토리의 `CLAUDE.md` 참조