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
