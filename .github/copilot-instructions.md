---
applyTo: "frontend/**"
---
# GTD-TODO 서비스 컨셉 문서.

## 주요 문제 상황

업무 중 할 일이 순차적으로 오지 않고 동시 다발적으로 주어져서 당장 어떤 일을 실행해야 할 지 난감한 상황이 존재.
혹은 어떤 일이나 생각, 관심을 나중에 해야지 하고 머릿 속으로만 뒀다가 결국 잊혀지게 되어버리는 상황이 존재.
TODO 앱을 활용하여도 할 일을 등록하는 과정 자체가 또 일처럼 느껴지는 상황이 존재.
머릿 속에 해야 할 일이 너무 많으면 그 자체로 너무 스트레스. 혹은 해야 할 일이 너무 막연해서 아무 실행도 못해서 스트레스를 받기도.

## 주요 컨셉

GTD를 통해서 사용자가 머릿속에 있는 모든 아이디어나 할 일을 부담 없이 서비스에 등록할 수 있게 하자.
해야 할 일을 머리 속이 아닌 시스템에 모두 맡기고 사용자는 지금 당장 해야 할 일에만 집중할 수 있게 하자.

그리고 일간, 주간, 월간으로 리포트를 생성하여 오래된 작업이나 보충되어야 할 부분을 정리해서 사용자에게 제공하자.

이때 우리 서비스의 주요한 특징은 아래와 같다. 

- 사용자가 할 일 등록하는 프로세스가 절대 일처럼 느껴지지 않도록 부담 없이 할 수 있어야 한다.
- 그리고 capture된 내용들을 clarify하는 과정에서도 부담스럽지 않게 해야 된다.
- 또한 GTD의 project 개념도 도입해서 큰 단위의 일을 여러 단계로 나눠서 진행할 수 있게도 제공해야 한다.
- 또한 capture되었던 일도 clarify하는 과정에서 delay나 remove 될 수 있다.

## 주요 컴포넌트 및 상태

### 생각(Thing)

사용자가 inbox에 자유롭게 기입하는 내용들.  
순간 떠오른 생각이나 갑자기 주어진 업무 등 아직 구조화 되지 않은 내용을 다룸.

#### 상태

- 대기 : 생성되어 분류되기를 대기
- 제거 : 분류 단계에서 삭제된 생각
- 완료 : 생각이 분류되어 할 일이나 프로젝트 등으로 파생
- 연기 : 당장 분류되지 않아도 되는 생각

### 할 일(Action)

사용자가 생각으로부터 당장 수행 가능올 인식된 단일 작업을 의미.  
생각과 다르게 구체적인 일정이나 실행에 필요한 내용 등 구조화 된 내용을 가짐.  

#### 상태

- 해야됨 : 생성되어 당장 수행 가능 상태
- 완료 : 할 일이 완료됨
- 수행중 : 할 일이 진행중
- 딜레이됨 : 예상 일자보다 늦어짐
- 연기함 : 수행은 가능하나 당장 수행이 아닌 할 일
- 위임함 : 수행 당사자가 사용자가 아님.
- 계획됨 : 특정 일자에 계획된 할 일임. 당장 수행할 일은 아님.
- 제거 : 삭제됨

### 프로젝트(Project)

사용자가 할 일로부터 당장 수행은 가능하다 단일 작업이 아닌 여러개의 할 일로 구성된 집합을 의미.
프로젝트 하위에는 여러 관련있는 할 일로 구성되어 있음.

- 해야됨 : 생성되어 당장 수행 가능함.
- 진행 중 : 프로젝트 하위 할일 중 일부가 수행중이거나 완료된 프로젝트
- 완료 : 모든 소속 할 일이 완료되고 완료처리함.
- 삭제 : 삭제됨.

## 주요 동작

### Capture

사용자가 생각을 자유롭게 inbox에 모으는 행위.
입력을 자유롭게 받는다.

### Clarify & Organize

Inbox에 수집된 생각을 당장 수행 가능한 할 일로 분류하여 ToDoList로 기록됨
이때 만약 2분내로 수행 가능하면 바로 시작을 하도록함.
이때 생각이 불필요한 내용이면 제거 처리하고 생각이 당장 수행 가능한 것은 아니지만 아카이빙할 내용이면 Archive에 저장. 혹은 나중에 다시 처리하고 싶으면 연기 처리.
만약 생각이 단일 할 일이 아닌 여러 할일로 구성될 것 같으면 프로젝트가 생성됨.

### Reflect

일간, 주간, 월간으로 진행된 일들, 오래동안 방치된 것들, 추가되지 않은 것 등 전체를 리뷰 해준다.

## 미래 전략 (아직 설계에는 포함 되어선 안됨)

아직은 구현 계획에는 없으나 아이디어 관점에서 서비스의 발전 방향을 설계.

- 할 일을 수집하는 capture 과정을 사용자가 직접 입력하지 않아도, 서비스가 사용자의 이메일, 메신저 대화를 기반으로 할 일을 추출하는 기능
- 사용자가 done한 내용을 학습하여, 사용자가 관심있어할 할 일을 서비스가 파이프라인을 통해 capture
# 🎨 GTD-TODO UI/UX 디자인 가이드 (v1.0)

> 작성자: Head of Product & Design (30년 경력 기준)  
> 목적: 프론트엔드, 디자이너, 마크업 개발자가 일관성 있게 적용할 수 있도록 시각/인터랙션 시스템 정의

---

## 1. 색상 시스템 (Color System)

| 토큰 | HEX | 용도 |
|------|-----|------|
| `--color-primary` | `#3478F6` | 주요 버튼, 강조 색상 |
| `--color-secondary` | `#F3F4F6` | 배경 카드, 라이트 존 |
| `--color-bg` | `#F9FAFB` | 전체 앱 배경 |
| `--color-text` | `#1F2937` | 본문 텍스트 |
| `--color-muted` | `#9CA3AF` | 보조 텍스트, 날짜, 설명 |
| `--color-border` | `#E5E7EB` | 입력창, 카드 경계선 |
| `--color-danger` | `#EF4444` | 제거, 삭제 관련 버튼 등 |
| `--color-success` | `#10B981` | 완료 상태 강조 |

---

## 2. 타이포그래피 (Typography)

- 폰트: `Inter`, `sans-serif`
- 제목 H1: 28px / Bold / `--color-text`
- 제목 H2: 20px / Semibold
- 본문: 16px / Regular
- 캡션: 12px / Muted

---

## 3. 간격 시스템 (Spacing)

- Grid: 12 column system, gutter 24px 기준
- Section padding: `24px` top-bottom
- 카드 내부 padding: `16px`
- 컴포넌트 간 margin: 기본 `12px`
- 입력창 height: 데스크탑 `48px`, 모바일 `40px`

---

## 4. 컴포넌트 스타일 가이드

### 📌 카드 (Card)
- border-radius: `12px`
- padding: `16px`
- 배경색: `--color-secondary`
- 그림자: `--shadow-card`

### 🔘 버튼 (Button)
- 기본 버튼: Primary color, radius `8px`, 높이 `40px`
- Hover: 밝은 음영 처리 (`opacity 90%`)
- Disabled: 색상 약화 (`--color-muted`)

### ✅ 체크박스/토글
- 라운드 스타일 체크 (스위치형)
- 상태 변화에 따라 color transition 적용 (`150ms ease-in-out`)

### 📥 입력창 (Input Field)
- 기본 padding: `12px`
- border: 1px solid `--color-border`
- Focus 상태: border-color `--color-primary`, 그림자 강조

---

## 5. 그림자 및 효과 (Shadows & Effects)

```css
--shadow-card: 0 2px 6px rgba(0,0,0,0.04);
--shadow-modal: 0 12px 24px rgba(0,0,0,0.12);
```

- 카드: 미묘한 입체감 강조 (flat하지 않도록)
- 모달: 배경과 분리되도록 대비 강조

---

## 6. 애니메이션 및 트랜지션

| 요소 | 효과 |
|------|------|
| 버튼 Hover | background subtle transition (150ms ease) |
| Modal 열기 | fade-in + scale (200ms) |
| Drag & Drop | `scale(1.05)` + 그림자 강조 |
| Quick Add 진입 | 슬라이드업 또는 fade (200ms) |

---

## 7. 반응형 디자인 기준 (Responsiveness)

- Mobile breakpoint: `max-width: 768px`
- 모바일 전용 하단 탭 메뉴 사용 (`Inbox`, `Add`, `Actions`, `Review`)
- 드래그 UX는 long-press 또는 버튼 기반 대체 제공
- 입력창, 카드 크기는 뷰포트에 따라 유동 배치

---

## ✅ 디자인 원칙 요약

- 구조는 단순하되 정보 계층은 명확하게 구분할 것
- 사용자가 입력하거나 조작할 때 즉각적 피드백 제공
- 가능한 한 클릭 없이도 흐름을 이어갈 수 있도록 구성
- 시각적으로는 "가볍고, 기능적으로 집중된" UI 지향
# 🧭 GTD-TODO 사용자 여정 및 앱 플로우 문서 (v1.0)

> 목적: 변경된 개념(Thing, Action, Project)을 바탕으로 사용자의 실제 흐름과 앱 구조를 정의

---

## 🎯 핵심 사용자 플로우 요약

```
[Thing Capture] → [Clarify & Organize] → [Action 실행] → [Project 정리] → [Reflect]
```

사용자는 머릿속 생각을 그대로 Thing으로 입력 → Clarify 화면에서 Action/Project로 변환 → Action 중심으로 실행 → Project 기반으로 구조화 → Reflect를 통해 전체 흐름 점검

---

## 👤 대표 사용자 시나리오

### 💼 직장인 시나리오
1. 회의 중 떠오른 업무를 Inbox에 빠르게 입력
2. 퇴근 전 Inbox에서 Clarify 수행 → 일부는 Action, 일부는 삭제 또는 연기
3. 다음날 출근 후 Next Actions에서 오늘 할 일 확인
4. 금요일 Review 탭에서 미완료 Action 정리, 오래된 프로젝트 확인

### 📚 학생 시나리오
1. 강의 중 과제나 아이디어를 Thing으로 캡처
2. 주말에 Clarify로 Action/Project로 전환
3. 학기별 프로젝트에 할 일 정리하여 프로젝트 관리
4. 월말 Review에서 연기된 과제 확인 및 일정 조정

---

## 📦 주요 엔터티 및 상태 요약

### ✅ Thing (생각)
- 상태: `대기`, `완료`, `연기`, `제거`
- 특성: 자유롭게 입력, 구조화되지 않음

### ✅ Action (할 일)
- 상태: `해야됨`, `완료`, `수행중`, `딜레이됨`, `연기함`, `위임함`, `계획됨`, `제거`
- 특성: 구체적 실행 단위, context 및 프로젝트 연결 가능

### ✅ Project (프로젝트)
- 상태: `해야됨`, `진행 중`, `완료`, `삭제`
- 특성: 여러 Action의 집합, 상태별 Action 수 기반 진행률 계산

---

## 🧩 전체 앱 메뉴 및 화면 구성

| 메뉴       | 경로             | 설명 |
|------------|------------------|------|
| Inbox      | `/inbox`         | Thing 등록 및 목록. 입력만 있고 구조화되지 않은 상태 |
| Clarify    | `/clarify/:id`   | Thing을 Action/Project로 명확화하거나 제거/연기 |
| Actions    | `/actions`       | 현재 수행해야 할 Action 리스트. 필터/정렬 기능 포함 |
| Projects   | `/projects`      | 프로젝트 목록 및 상세 화면. 하위 Action 정리 |
| Review     | `/review`        | 일간/주간/월간 기준으로 미완료, 정체, 연기 항목 리뷰 |

---

## 🔄 유저 플로우 예시

### 캡처 → 명확화 → 실행 흐름
```
사용자 입력 → Inbox → Clarify 진입 → Action 생성 or 제거/연기 → Actions → 완료 처리
```

### 프로젝트 생성 흐름
```
Thing → Clarify에서 Project로 분기 생성 → Project 상세에서 Action 등록 → Actions와 병렬로 실행
```

### 리뷰 흐름
```
금요일 → Review → 오래된 Action/Thing 확인 → 필요한 건 다시 Clarify로 이동 or 삭제
```

---

## 🧠 UX 설계 포인트

- 모든 흐름은 마우스 및 키보드로 자연스럽게 처리
- Clarify는 빠르게 끝낼 수 있도록 드래그, 토글, 단축키 중심
- Review는 사용자의 행동 피드백을 명확히 제공해야 함 (예: 완료율, 미완료 항목 강조)
- Inbox는 어디서든 진입 가능해야 하며, 전역 단축키로 Capture 유도
# Copilot Development Guide for GTD-TODO

## Overview
This document provides essential knowledge for AI coding agents to be productive in the GTD-TODO project. It supplements the existing conceptual and design documentation with actionable development workflows and codebase-specific patterns.

## Codebase Structure
- **Backend** (`backend/`):
  - Written in Go, it handles the core business logic and API endpoints.
  - Key files:
    - `main.go`: Entry point of the application.
    - `thing/handler.go`: HTTP handlers for "Thing" resources.
    - `thing/service.go`: Business logic for "Thing" management.
    - `thing/model.go`: Data models for "Thing" entities.
  - Dependencies managed via Go modules (`go.mod`, `go.sum`).

- **Frontend** (`frontend/`):
  - Contains a basic HTML file (`html.html`) for user interaction.

- **Documentation** (`docs/`):
  - Includes design guides (`gtd_todo_design_guide.md`) and user flows (`userflow.md`).

## Development Workflows

### Building the Backend
To compile the Go application:
```bash
cd backend && go build
```

### Running the Backend
To start the backend server:
```bash
cd backend && go run main.go
```

### Testing
Add test files in the `*_test.go` format and run:
```bash
cd backend && go test ./...
```

### Debugging
Use Go's built-in debugging tools or an IDE like VS Code with Go extensions.

## Conventions and Patterns
- **Handler-Service-Model Pattern**:
  - Handlers (`handler.go`) manage HTTP requests and responses.
  - Services (`service.go`) encapsulate business logic.
  - Models (`model.go`) define data structures.

- **Error Handling**:
  - Follow Go's idiomatic `if err != nil` pattern.

## Integration Points
- The backend serves as the API provider for the frontend.
- No external APIs or databases are currently integrated.

## Notes for AI Agents
- Ensure Go modules are up-to-date by running `go mod tidy` in the `backend/` directory.
- The frontend is minimal and may require further development for a complete user experience.
- Refer to `docs/` for additional context on design and user flows.


## frontend Development
- Use HTML, CSS, and Vanila JavaScript for the frontend.
- Using the directory `frontend/` for all frontend related files.
- frontend uses backend API for data operations.
- check which API endpoints are available in the backend and how to use them.(`backend` directory)