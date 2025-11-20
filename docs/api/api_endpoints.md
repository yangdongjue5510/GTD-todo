# API Endpoints

RESTful API 엔드포인트 설계 (Phase 1 기준)

---

## 인증 (Authentication)

| Method | Endpoint | Request | Response |
|--------|----------|---------|----------|
| POST | `/api/auth/signup` | `{email, password}` | `{user_id, token}` |
| POST | `/api/auth/login` | `{email, password}` | `{user_id, token}` |
| POST | `/api/auth/logout` | - | `{message}` |

**인증**: JWT Bearer Token

---

## TODO 관리

### CRUD
| Method | Endpoint | Request | Response |
|--------|----------|---------|----------|
| POST | `/api/todos` | `{title, description?, project_id?, status?}` | `{todo}` |
| GET | `/api/todos` | Query: `status?, project_id?, sort?, order?` | `{todos: [], total}` |
| GET | `/api/todos/:id` | - | `{todo}` |
| PATCH | `/api/todos/:id` | `{title?, description?, project_id?, status?, position?}` | `{todo}` |
| DELETE | `/api/todos/:id` | - | `{message}` |

**Query Parameters** (GET `/api/todos`):
- `status`: inbox, next_actions, in_progress, done, someday, waiting_for
- `project_id`: 프로젝트 필터
- `sort`: created_at, updated_at, position (default: position)
- `order`: asc, desc (default: asc)

### 상태 변경
| Method | Endpoint | Request | Response |
|--------|----------|---------|----------|
| PATCH | `/api/todos/:id/status` | `{status}` | `{todo}` |

### 순서 변경
| Method | Endpoint | Request | Response |
|--------|----------|---------|----------|
| PATCH | `/api/todos/:id/position` | `{new_position, status?}` | `{todo}` |

---

## 프로젝트 관리

| Method | Endpoint | Request | Response |
|--------|----------|---------|----------|
| POST | `/api/projects` | `{name, description?, color?}` | `{project}` |
| GET | `/api/projects` | - | `{projects: []}` |
| GET | `/api/projects/:id` | - | `{project, todo_count}` |
| PATCH | `/api/projects/:id` | `{name?, description?, color?}` | `{project}` |
| DELETE | `/api/projects/:id` | - | `{message}` |

---

## 대시보드

| Method | Endpoint | Response |
|--------|----------|----------|
| GET | `/api/dashboard/stats` | `{inbox_count, next_actions_count, in_progress_count, done_count, someday_count, waiting_for_count, total_count, projects_count}` |

---

## 헬스체크

| Method | Endpoint | Response |
|--------|----------|----------|
| GET | `/api/health` | `{status}` |


## 응답 형식

### 에러
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

### HTTP Status
- `200 OK`: 성공
- `201 Created`: 생성 성공
- `400 Bad Request`: 잘못된 요청
- `401 Unauthorized`: 인증 실패
- `404 Not Found`: 리소스 없음
- `500 Internal Server Error`: 서버 오류

---

## Phase 1 우선순위

**Must Have** (Phase 1.1-1.3):
1. `/api/auth/*` (회원가입, 로그인)
2. `/api/todos` (CRUD)
3. `/api/todos` (필터링: status, project_id)
4. `/api/todos/:id/status` (상태 전환)
5. `/api/projects` (CRUD)
6. `/api/health` (헬스체크)

**Should Have** (Phase 1.4):
6. `/api/todos/:id/position` (순서 변경)
7. `/api/dashboard/stats` (통계)
