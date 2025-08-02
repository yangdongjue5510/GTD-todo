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
