# Tech Stack

---

## Backend

| 구분 | 기술 | 버전 |
|------|------|------|
| 언어 | Go | 1.21+ |
| 프레임워크 | Gin | 1.9+ |
| DB | PostgreSQL | 15+ |
| 마이그레이션 | golang-migrate | 4.16+ |
| 인증 | JWT | - |
| 비밀번호 | bcrypt | - |

---

## Frontend

| 구분 | 기술 |
|------|------|
| 언어 | Vanilla JavaScript (ES6+) |
| CSS | Tailwind CSS |
| 드래그앤드롭 | HTML5 Drag & Drop API |
| HTTP | Fetch API |

---

## Infrastructure

**개발**: Docker Compose (PostgreSQL)

**배포**:
- Backend: Fly.io / Railway
- Frontend: Vercel / Netlify
- DB: Supabase (무료 500MB)
