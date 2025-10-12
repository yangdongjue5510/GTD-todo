# CLAUDE.md

이 파일은 Claudecode 작업 지시를 다룸.

## 프로젝트 작업 방식

GTD-TODO 프로젝트는 단계별 설계 프롬프트를 기반으로 작업.

### 작업 프롬프트 종류

1. **Product Design** (`/docs/prompts/product_design_prompt.md`)
   - 제품 요구사항 및 전체 비전 정의
   
2. **System Design** (`/docs/prompts/system_architect_design_prompt.md`)
   - 시스템 아키텍처 및 구조 설계
   
3. **Feature Design** (`/docs/prompts/feature_design_prompt.md`)
   - 개별 기능의 상세 설계
   
4. **Feature Implement** (`/docs/prompts/feature_implement_prompt.md`)
   - 설계된 기능의 실제 구현

### 작업 진행 방법

1. **작업 시작 시**: 사용자가 작업 유형을 명시하지 않으면, 작업 내용을 분석하여 적절한 프롬프트 선택
2. **프롬프트 읽기**: 해당 프롬프트 파일을 Read 도구로 읽고 지시사항 확인
3. **단계별 진행**: 프롬프트에 명시된 단계와 체크리스트를 따라 작업 수행
4. **TodoWrite 활용**: 각 프롬프트의 단계를 todo 항목으로 변환하여 진행 상황 추적

### 주의사항

- 각 프롬프트의 지시사항은 **반드시** 준수해야 함
- 단계를 건너뛰지 말고 순차적으로 진행
- 설계 단계에서는 구현 코드를 작성하지 않음