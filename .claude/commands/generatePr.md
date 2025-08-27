---
description: 변경사항을 커밋하고 푸시한 후 GitHub PR을 생성
allowed-tools: Bash(git *, gh *)
---

모든 커밋과 푸시가 완료되면:
1. 전체 커밋 히스토리를 분석하여 PR 제목과 설명을 작성
2. GitHub CLI(`gh pr create`)를 사용하여 PR 생성
3. PR 제목은 주요 변경사항을 요약
4. PR 설명은 다음 형식을 따름:
   ```
   ## Summary
   - 주요 변경사항 1
   - 주요 변경사항 2
   - 주요 변경사항 3
   
   ## Changes
   - 구체적인 변경 내용들을 나열
   
   ## Test plan
   - [ ] 코드 컴파일 확인
   - [ ] 기존 기능 동작 확인
   - [ ] 새로운 기능 테스트
   
   🤖 Generated with [Claude Code](https://claude.ai/code)
   ```

PR 생성 후 PR URL을 반환한다.