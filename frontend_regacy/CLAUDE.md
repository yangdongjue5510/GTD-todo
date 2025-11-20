# CLAUDE.md - 프론트엔드

이 파일은 GTD-todo 클라이언트 애플리케이션의 프론트엔드 관련 가이드를 제공합니다.

## 아키텍처

### 기술 스택
- **HTML**: 최소한의 DOM 요소로 구성된 의미론적 구조
- **CSS**: 테마 설정을 위한 CSS 커스텀 속성(변수)
- **JavaScript**: API 호출을 위한 순수 ES6+ async/await 사용
- **프레임워크 없음**: "부담 없음"이라는 GTD 철학을 유지하기 위해 의도적으로 단순화

### 파일 구조
- `index.html`: 인박스 인터페이스가 포함된 메인 애플리케이션 구조
- `styles.css`: 디자인 시스템 변수가 포함된 CSS
- `app.js`: DOM 조작 및 API 통합

## 개발

### 로컬 개발
```bash
# 파일 서빙 (로컬 서버 사용)
python -m http.server 3000
# 또는
npx serve .
# 또는 브라우저에서 index.html 직접 열기
```

### 백엔드 통합
- API 베이스 URL: `http://localhost:8080`
- Content-Type: `application/json`
- 로컬 개발을 위한 CORS 활성화

## 디자인 시스템

### CSS 변수
```css
:root {
  --bg: #f8f9fa;        /* 앱 배경 */
  --card: #ffffff;       /* 카드/컨테이너 배경 */  
  --primary: #4f46e5;    /* 주요 액션 */
  --text: #111827;       /* 메인 텍스트 */
  --border: #e5e7eb;     /* 경계선 및 구분선 */
}
```

### 컴포넌트 패턴

#### 입력 행
```html
<div class="input-row">
  <input type="text" placeholder="What's on your mind?">
  <input type="text" id="description" placeholder="Add a description (optional)">
  <button>Add</button>
</div>
```

#### Thing 항목 표시
```html
<div class="thing-item">
  <strong>제목</strong><br>
  <em>설명 또는 "설명 없음"</em><br>
  상태: [상태_값]
</div>
```

## JavaScript 패턴

### API 호출
적절한 오류 처리와 함께 async/await 사용:
```javascript
async function fetchThings() {
    try {
        const response = await fetch("http://localhost:8080/things/");
        if (!response.ok) throw new Error("Failed to fetch things");
        const things = await response.json();
        // 응답 처리...
    } catch (error) {
        console.error("Error fetching things:", error);
    }
}
```

### DOM 조작
- 단일 요소에는 `querySelector` 사용
- 이벤트 처리에는 `addEventListener` 사용
- 복잡한 상태 관리보다는 목록 지우기 및 재구축

### 이벤트 처리
```javascript
// 버튼 클릭
addButton.addEventListener("click", handleAdd);

// 엔터 키 지원
inputField.addEventListener("keypress", (event) => {
    if (event.key === "Enter") {
        addButton.click();
    }
});
```

## GTD UI 원칙

### 포착 인터페이스
- **최소한의 마찰**: 선택적 설명이 포함된 단일 입력 필드
- **즉각적인 피드백**: 성공적인 제출 후 입력 필드 지우기
- **키보드 지원**: 엔터 키로 제출
- **시각적 단순함**: 깔끔하고 정리된 인박스 뷰

### 상태 표시
- 상태를 숫자가 아닌 읽을 수 있는 텍스트로 표시
- 콘텐츠 계층을 위한 의미론적 HTML (`<strong>`, `<em>`) 사용
- 일관된 간격 및 타이포그래피

## 현재 구현 상태

### 구현된 기능
- ✅ 모듈화된 ES6+ JavaScript 아키텍처
- ✅ 클라이언트 사이드 라우팅 (해시 기반)
- ✅ 다중 뷰 지원 (Inbox, Clarify, Actions, Projects, Review)
- ✅ Thing 포착 및 목록 관리
- ✅ 명확화 인터페이스 (Thing → Action 변환)
- ✅ Action 목록 뷰 (필터링 지원)
- ✅ Review 대시보드 (통계 및 활동 내역)
- ✅ 토스트 알림 시스템
- ✅ 로딩 상태 관리
- ✅ 로컬 스토리지 캐싱
- ✅ 키보드 단축키 (Ctrl/Cmd + 1-4, Ctrl/Cmd + N)
- ✅ 반응형 디자인
- ✅ 접근성 (ARIA, 포커스 관리)
- ✅ 에러 처리 및 네트워크 상태 감지

### 부분 구현된 기능
- 🔄 Action 상태 변경 (UI는 구현, API 대기)
- 🔄 Project 관리 (기본 구조만 구현)
- 🔄 자동 저장 (입력 상태 복원)

## 향후 개발 가이드라인

## 코드 구조

```
js/
├── app.js              # 메인 애플리케이션 진입점
├── components/         # 재사용 가능한 UI 컴포넌트
│   └── navigation.js   # 네비게이션 컴포넌트
├── services/          # 비즈니스 로직 및 외부 통신
│   ├── api.js         # REST API 통신
│   ├── storage.js     # 로컬 스토리지 관리
│   ├── toast.js       # 알림 메시지 서비스
│   └── loading.js     # 로딩 상태 관리
├── utils/             # 유틸리티 함수
│   └── router.js      # 클라이언트 사이드 라우터
└── views/             # 페이지별 뷰 컴포넌트
    ├── base.js        # 베이스 뷰 클래스
    ├── inbox.js       # Inbox 뷰
    ├── clarify.js     # 명확화 뷰
    ├── actions.js     # Actions 뷰
    ├── projects.js    # Projects 뷰
    └── review.js      # Review 뷰
```

## 아키텍처 패턴

### MVC/MVP 패턴
- **Model**: API 서비스를 통한 데이터 관리
- **View**: 각 화면별 뷰 클래스
- **Controller**: 라우터가 뷰 간 전환 제어

### 서비스 레이어
- 각 뷰는 공통 서비스에 의존성 주입
- 서비스 간 느슨한 결합
- 단일 책임 원칙 준수

### 이벤트 기반 아키텍처
- 라우터 이벤트를 통한 뷰 전환
- 사용자 상호작용에 대한 이벤트 위임
- 서비스 간 이벤트 통신

## 새 기능 추가 시 가이드라인

### 새 뷰 추가
1. `views/` 디렉토리에 새 뷰 클래스 생성
2. `BaseView`를 상속받아 기본 기능 활용
3. `app.js`에서 뷰 인스턴스 등록
4. HTML에 해당 섹션 추가

### 새 서비스 추가
1. `services/` 디렉토리에 서비스 클래스 생성
2. 필요시 싱글톤 패턴 적용
3. `app.js`에서 의존성 주입
4. 에러 처리 및 로깅 포함

### UI 컴포넌트 추가
1. `components/` 디렉토리에 컴포넌트 생성
2. CSS 변수를 활용한 일관된 스타일
3. 접근성 고려 (ARIA, 키보드 네비게이션)
4. 재사용성을 위한 설정 가능한 옵션

## 개발 모범 사례

- ✅ ES6+ 문법 사용 (모듈, 클래스, async/await)
- ✅ JSDoc 주석으로 문서화
- ✅ 에러 바운더리 및 graceful degradation
- ✅ 디바운싱/쓰로틀링으로 성능 최적화
- ✅ 메모리 누수 방지 (이벤트 리스너 정리)
- ✅ 시맨틱 HTML과 ARIA 레이블
- ✅ 반응형 디자인 우선 고려