# 📘 GTD-TODO 설계 문서 (v1.0)

> 작성자: Head of Product & Design (30년 경력 기준)  
> 목적: 개발/디자인 팀 전체가 공유할 수 있는 전사 기준 문서

---

## 1. 단계별 구현 전략 (Execution Plan)

### ✅ Phase 0 — 기술/구조 셋업 (Day 0~2)

- Monorepo 구조: `/frontend`, `/backend`, `/docs`
- Go Gin 기반 REST API 서버 구축
- Vanilla JS 기반 프론트엔드 초기화
- SQLite 또는 PostgreSQL 로컬 DB 연결
- Task Runner: makefile 또는 pnpm 기반 Dev Toolchain 구성

---

### ✅ Phase 1 — 핵심 개념 3종 도입 (Week 1~2)

#### 1.1 Capture 기능 (Thing 등록)

- 자유 입력 필드 + Enter → `thing` 등록 (status: `대기`)
- `/things` POST, GET API
- Inbox 화면 구성 (등록된 Thing 리스트)

#### 1.2 Clarify & Organize 기능 (Thing → Action/Project)

- Thing → Action 변환, 또는 Project 생성 흐름 구현
- Clarify 도중 제거, 연기 처리 가능하도록 분기
- `/things/:id/clarify` PATCH API 설계
- Action 구조 및 상태 흐름 정의 (해야됨/완료/위임 등)

#### 1.3 Action 화면 구성

- `/actions` GET (필터: context, status)
- 완료 체크, 연기 버튼, 우선순위 정렬 적용

#### 1.4 Project 구조 구현

- `/projects` GET/POST/DELETE
- 프로젝트별 Action 묶음 구조 및 진행률 표시

---

### ✅ Phase 2 — Reflect 기능 및 정기 리뷰 (Week 3)

- `/review/daily`, `/weekly`, `/monthly` 뷰
- 오래된 Thing/Action 탐지 및 정리 제안
- 미처리 Task 수, 연기 횟수, 완료율 등 KPI 통계 시각화

---

### ✅ Phase 3 — UI 개선 및 보조 기능 (Week 4)

- Drag & Drop 기반 Clarify UX
- 전역 Capture 단축키 (`Cmd+K`) / 플로팅 Quick Add
- localStorage + 서버 Sync 캐시 전략
- 반응형 모바일 최적화 (모바일 Inbox 우선)

---

## 2. 전체 사용자 여정 및 앱 플로우

### 🔁 사용자 핵심 플로우 요약

```
[Thing Capture] → [Clarify & Organize] → [Action 실행] ↔ [Project 정리] → [Reflect]
```

---

### 📍 주요 메뉴/페이지

| 메뉴      | 경로              | 기능 요약                                  |
|-----------|-------------------|---------------------------------------------|
| Inbox     | `/inbox`          | 모든 Thing 모음. 등록 & Clarify 시작 지점   |
| Clarify   | `/clarify/:id`    | Thing → Action/Project 변환, 제거/연기 처리 |
| Actions   | `/actions`        | 현재 해야 할 일 목록                         |
| Projects  | `/projects`       | 프로젝트별 할 일 목록 및 진행률 추적        |
| Review    | `/review`         | 일간/주간/월간 리뷰                          |

---

### 📦 데이터 구조 단위

- **Thing**  
  - 필드: id, title, memo?, createdAt  
  - 상태: `대기`, `제거`, `완료`, `연기`

- **Action**  
  - 필드: id, title, dueDate?, context?, projectId?, createdAt  
  - 상태: `해야됨`, `완료`, `수행중`, `딜레이됨`, `연기함`, `위임함`, `계획됨`, `제거`

- **Project**  
  - 필드: id, title, description?, createdAt  
  - 상태: `해야됨`, `진행 중`, `완료`, `삭제`

---

## 3. 제품 디자인 가이드 (Design System v1.0)

### 🎨 색상 시스템

| 토큰                | HEX        | 용도              |
|---------------------|------------|-------------------|
| `--color-primary`   | `#3478F6`  | 주요 버튼, 링크     |
| `--color-bg`        | `#F9FAFB`  | 배경색             |
| `--color-border`    | `#E5E7EB`  | 카드, 입력 경계선   |
| `--color-text`      | `#1F2937`  | 본문 텍스트        |
| `--color-muted`     | `#9CA3AF`  | 날짜, 보조 정보     |
| `--color-danger`    | `#EF4444`  | 제거, 삭제         |

---

### 🖋️ 폰트 및 타이포그래피

- 기본 폰트: `Inter`, `sans-serif`
- 제목 (H1): 28px, Bold
- 제목 (H2): 20px, Semibold
- 본문: 16px, Regular
- 캡션: 12px, Muted

---

### 📏 마진 & 패딩

- 카드 내부 padding: `16px`
- 열 간격: `24px`
- 컴포넌트 간 간격: `12px`
- 입력창 높이: `48px` (모바일: `40px`)

---

### 🌫️ 그림자 (Shadow)

```css
--shadow-card: 0 2px 6px rgba(0, 0, 0, 0.04);
--shadow-modal: 0 12px 24px rgba(0, 0, 0, 0.12);
```

---

### ✨ 인터랙션 & 애니메이션

- Drag & Drop: `transform: scale(1.05)` + shadow 강조
- 버튼 hover: `background-color` 부드럽게 전환 (`150ms ease-in-out`)
- 모달 열기: fade-in + scale 효과 (200ms)
- Quick Add: 슬라이드 업 애니메이션

---

### 📱 반응형 대응

- 모바일 뷰: Inbox + Quick Add 우선 표시
- 너비 ≤ 768px: 하단 탭 메뉴 + 플로팅 버튼
- 드래그 UX: long-press 대체 UI (모바일 대응)

---

## 🔚 마무리

> 본 문서는 GTD-TODO 제품 구조와 철학의 기준점입니다.  
> 기획 변경이 발생하면 반드시 본 문서를 기준으로 구조를 재검토해야 합니다.
