/**
 * Projects 뷰
 * 프로젝트 목록 관리
 */
import { BaseView } from './base.js';

export class ProjectsView extends BaseView {
  constructor(services) {
    super(services);
    this.element = document.getElementById('projects-view');
    this.projects = [];
    this.init();
  }

  /**
   * 뷰 초기화
   */
  init() {
    if (!this.element) return;
    // 현재는 기본 구현만
  }

  /**
   * 뷰 표시 시 호출
   */
  async onShow(params) {
    await this.loadProjects();
  }

  /**
   * Project 목록 로드
   */
  async loadProjects() {
    try {
      // 현재는 Project API가 없으므로 빈 배열
      this.projects = [];
      
      this.renderProjects();
      this.storage.saveRecentActivity('Projects 새로고침');
      
    } catch (error) {
      this.handleError(error, 'Project 목록 로드');
    }
  }

  /**
   * Project 목록 렌더링
   */
  renderProjects() {
    const container = this.$('#project-list');
    if (!container) return;

    // 현재는 빈 상태만 표시
    container.innerHTML = `
      <div class="empty-state">
        <h3>아직 프로젝트가 없습니다</h3>
        <p class="text-muted">
          복잡한 작업을 Inbox에서 명확화할 때 프로젝트로 만들 수 있습니다.
        </p>
        <p class="text-muted">
          프로젝트는 여러 개의 관련된 Action들을 그룹화하여 관리할 수 있게 해줍니다.
        </p>
      </div>
    `;
  }
}