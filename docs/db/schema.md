# Database Schema

PostgreSQL 기반 GTD-TODO 스키마 (Hard Delete)

---

## 1. users

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);
```

- `id`: 사용자 ID
- `email`: 로그인 ID (UNIQUE)
- `password_hash`: bcrypt 해시

---

## 2. projects

```sql
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color VARCHAR(7) DEFAULT '#3B82F6',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_projects_user_id ON projects(user_id);
```

- `user_id` → `users(id)` CASCADE (사용자 삭제 시 프로젝트도 삭제)
- `color`: HEX 색상 코드

---

## 3. todos

```sql
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'inbox'
        CHECK (status IN ('inbox', 'next_actions', 'in_progress', 'done', 'someday', 'waiting_for')),
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_project_id ON todos(project_id);
CREATE INDEX idx_todos_user_status ON todos(user_id, status);
```

- `user_id` → `users(id)` CASCADE
- `project_id` → `projects(id)` SET NULL (프로젝트 삭제 시 TODO 유지)
- `status`: GTD 6단계 (VARCHAR + CHECK 제약조건)
  - 값: inbox, next_actions, in_progress, done, someday, waiting_for
  - **ENUM 대신 VARCHAR + CHECK 선택 이유**: 애자일 방식에서 상태 추가/변경/삭제 유연성 확보
- `position`: 드래그앤드롭 순서 (동일 status 내)

---

## ERD

```
users (1) ──┬─< projects (N)  [CASCADE]
            └─< todos (N)     [CASCADE]
                  └──< projects (0..1)  [SET NULL]
```
