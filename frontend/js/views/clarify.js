/**
 * Clarify 뷰
 * Thing을 Action으로 명확화하는 화면
 */
import { BaseView } from './base.js';

export class ClarifyView extends BaseView {
  constructor(services) {
    super(services);
    this.element = document.getElementById('clarify-view');
    this.currentThing = null;
    this.init();
  }

  /**
   * 뷰 초기화
   */
  init() {
    if (!this.element) return;

    const backButton = this.$('#clarify-back');
    if (backButton) {
      this.addEventListener(backButton, 'click', () => {
        this.router.navigate('inbox');
      });
    }
  }

  /**
   * 뷰 표시 시 호출
   * @param {Object} params - 라우트 파라미터
   */
  async onShow(params) {
    const thingId = params.thingId;
    
    if (!thingId) {
      this.toast.error('명확화할 항목을 찾을 수 없습니다.');
      this.router.navigate('inbox');
      return;
    }

    await this.loadThing(parseInt(thingId));
  }

  /**
   * Thing 데이터 로드
   * @param {number} thingId - Thing ID
   */
  async loadThing(thingId) {
    try {
      // 현재는 단일 Thing API가 없으므로 전체 목록에서 찾기
      const things = await this.loadWithState(
        () => this.api.getThings(),
        'Thing 정보를 불러오는 중...'
      );

      this.currentThing = things.find(t => t.id === thingId);
      
      if (!this.currentThing) {
        this.toast.error('해당 항목을 찾을 수 없습니다.');
        this.router.navigate('inbox');
        return;
      }

      this.renderClarifyContent();
      
    } catch (error) {
      this.handleError(error, 'Thing 로드');
      this.router.navigate('inbox');
    }
  }

  /**
   * 명확화 콘텐츠 렌더링
   */
  renderClarifyContent() {
    const container = this.$('#clarify-content');
    if (!container || !this.currentThing) return;

    container.innerHTML = `
      <div class="clarify-thing">
        <div class="thing-preview">
          <h3 class="thing-title">${this.escapeHtml(this.currentThing.title)}</h3>
          <p class="thing-description">${this.escapeHtml(this.currentThing.description || '설명 없음')}</p>
        </div>

        <div class="clarify-options">
          <h4>이것을 어떻게 처리하시겠습니까?</h4>
          
          <div class="clarify-option">
            <button type="button" class="btn btn-primary option-btn" data-action="quick">
              <span class="option-icon">⚡</span>
              <span class="option-text">
                <strong>2분 안에 할 수 있는 일</strong>
                <small>바로 완료하기</small>
              </span>
            </button>
          </div>

          <div class="clarify-option">
            <button type="button" class="btn btn-primary option-btn" data-action="action">
              <span class="option-icon">✅</span>
              <span class="option-text">
                <strong>실행 가능한 할 일</strong>
                <small>Action 목록에 추가</small>
              </span>
            </button>
          </div>

          <div class="clarify-option">
            <button type="button" class="btn btn-primary option-btn" data-action="project">
              <span class="option-icon">📁</span>
              <span class="option-text">
                <strong>여러 단계가 필요한 일</strong>
                <small>Project로 만들기</small>
              </span>
            </button>
          </div>

          <div class="clarify-option">
            <button type="button" class="btn btn-secondary option-btn" data-action="someday">
              <span class="option-icon">📅</span>
              <span class="option-text">
                <strong>나중에 할 일</strong>
                <small>Someday/Maybe 목록으로</small>
              </span>
            </button>
          </div>

          <div class="clarify-option">
            <button type="button" class="btn btn-secondary option-btn" data-action="delete">
              <span class="option-icon">🗑️</span>
              <span class="option-text">
                <strong>불필요한 일</strong>
                <small>삭제하기</small>
              </span>
            </button>
          </div>
        </div>
      </div>
    `;

    this.setupClarifyHandlers();
  }

  /**
   * 명확화 옵션 이벤트 핸들러 설정
   */
  setupClarifyHandlers() {
    const container = this.$('#clarify-content');
    if (!container) return;

    this.addEventListener(container, 'click', (event) => {
      const optionBtn = event.target.closest('.option-btn');
      if (!optionBtn) return;

      const action = optionBtn.dataset.action;
      this.handleClarifyAction(action);
    });
  }

  /**
   * 명확화 액션 처리
   * @param {string} action - 선택된 액션
   */
  async handleClarifyAction(action) {
    if (!this.currentThing) return;

    try {
      switch (action) {
        case 'quick':
          await this.handleQuickAction();
          break;
        case 'action':
          await this.handleCreateAction();
          break;
        case 'project':
          await this.handleCreateProject();
          break;
        case 'someday':
          await this.handleSomedayAction();
          break;
        case 'delete':
          await this.handleDeleteAction();
          break;
      }
    } catch (error) {
      this.handleError(error, '명확화 처리');
    }
  }

  /**
   * 빠른 완료 처리
   */
  async handleQuickAction() {
    const confirmed = confirm('이 작업을 지금 바로 완료하시겠습니까?');
    if (!confirmed) return;

    await this.loadWithState(
      () => this.api.updateThingStatus(this.currentThing.id, 2), // Done 상태
      '완료 처리하는 중...'
    );

    this.toast.success('작업이 완료로 처리되었습니다.');
    this.storage.saveRecentActivity(`빠른 완료: ${this.currentThing.title}`);
    this.router.navigate('inbox');
  }

  /**
   * Action 생성 처리
   */
  async handleCreateAction() {
    // 현재는 Action API가 구현되지 않았으므로 임시로 완료 상태로 변경
    this.toast.info('Action 기능은 곧 구현될 예정입니다.');
    
    await this.loadWithState(
      () => this.api.updateThingStatus(this.currentThing.id, 2),
      'Action으로 변환하는 중...'
    );

    this.storage.saveRecentActivity(`Action 생성: ${this.currentThing.title}`);
    this.router.navigate('actions');
  }

  /**
   * Project 생성 처리
   */
  async handleCreateProject() {
    this.toast.info('Project 기능은 곧 구현될 예정입니다.');
    
    await this.loadWithState(
      () => this.api.updateThingStatus(this.currentThing.id, 2),
      'Project로 변환하는 중...'
    );

    this.storage.saveRecentActivity(`Project 생성: ${this.currentThing.title}`);
    this.router.navigate('projects');
  }

  /**
   * Someday 처리
   */
  async handleSomedayAction() {
    await this.loadWithState(
      () => this.api.updateThingStatus(this.currentThing.id, 1), // Someday 상태
      'Someday로 이동하는 중...'
    );

    this.toast.success('Someday/Maybe 목록으로 이동되었습니다.');
    this.storage.saveRecentActivity(`Someday 이동: ${this.currentThing.title}`);
    this.router.navigate('inbox');
  }

  /**
   * 삭제 처리
   */
  async handleDeleteAction() {
    const confirmed = confirm(`\"${this.currentThing.title}\"을(를) 정말 삭제하시겠습니까?`);
    if (!confirmed) return;

    await this.loadWithState(
      () => this.api.deleteThing(this.currentThing.id),
      '삭제하는 중...'
    );

    this.toast.success('항목이 삭제되었습니다.');
    this.storage.saveRecentActivity(`삭제: ${this.currentThing.title}`);
    this.router.navigate('inbox');
  }

  /**
   * HTML 이스케이프 처리
   * @param {string} text - 이스케이프할 텍스트
   * @returns {string} 이스케이프된 텍스트
   */
  escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }
}