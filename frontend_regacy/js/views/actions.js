/**
 * Actions 뷰
 * 실행 가능한 Action 목록 관리
 */
import { BaseView } from './base.js';

export class ActionsView extends BaseView {
  constructor(services) {
    super(services);
    this.element = document.getElementById('actions-view');
    this.actions = [];
    this.currentFilter = 'all';
    this.init();
  }

  /**
   * 뷰 초기화
   */
  init() {
    if (!this.element) return;

    this.setupFilterHandlers();
  }

  /**
   * 필터 이벤트 핸들러 설정
   */
  setupFilterHandlers() {
    const filterBar = this.$('.filter-bar');
    if (!filterBar) return;

    this.addEventListener(filterBar, 'click', (event) => {
      const filterBtn = event.target.closest('.filter-btn');
      if (!filterBtn) return;

      const filter = filterBtn.dataset.filter;
      this.setActiveFilter(filter);
      this.renderFilteredActions();
    });
  }

  /**
   * 뷰 표시 시 호출
   */
  async onShow(params) {
    await this.loadActions();
    
    // URL 파라미터로 필터 설정
    if (params.filter) {
      this.setActiveFilter(params.filter);
    }
  }

  /**
   * Action 목록 로드
   */
  async loadActions() {
    try {
      // 현재는 Action API가 없으므로 완료된 Thing들을 Action으로 표시
      const things = await this.loadWithState(
        () => this.api.getThings(),
        'Action 목록을 불러오는 중...'
      );
      
      // 완료된 Thing들을 Action으로 변환
      this.actions = things
        .filter(thing => thing.status === 2) // Done 상태인 것들
        .map(thing => ({
          id: thing.id,
          title: thing.title,
          description: thing.description,
          status: 'done', // 임시로 모두 완료 상태
          createdAt: new Date(), // 임시 날짜
          completedAt: new Date()
        }));

      this.renderFilteredActions();
      this.updateFilterBadges();
      this.storage.saveRecentActivity('Actions 새로고침');
      
    } catch (error) {
      this.handleError(error, 'Action 목록 로드');
    }
  }

  /**
   * 필터링된 Action 목록 렌더링
   */
  renderFilteredActions() {
    const container = this.$('#action-list');
    if (!container) return;

    const filteredActions = this.getFilteredActions();

    // 빈 상태 처리
    if (filteredActions.length === 0) {
      container.innerHTML = this.getEmptyStateHTML();
      return;
    }

    // Action 목록 렌더링
    const actionsHTML = filteredActions
      .map(action => this.renderActionItem(action))
      .join('');

    container.innerHTML = actionsHTML;
    this.setupActionHandlers();
  }

  /**
   * 필터에 따른 Action 목록 반환
   * @returns {Array} 필터링된 Action 목록
   */
  getFilteredActions() {
    switch (this.currentFilter) {
      case 'todo':
        return this.actions.filter(action => action.status === 'todo');
      case 'doing':
        return this.actions.filter(action => action.status === 'doing');
      case 'done':
        return this.actions.filter(action => action.status === 'done');
      default:
        return this.actions;
    }
  }

  /**
   * 개별 Action 아이템 렌더링
   * @param {Object} action - Action 객체
   * @returns {string} HTML 문자열
   */
  renderActionItem(action) {
    const statusText = this.getStatusText(action.status);
    const statusClass = this.getStatusClass(action.status);
    const description = action.description || '설명 없음';
    
    return `
      <div class="action-item ${statusClass}" data-action-id="${action.id}" role="listitem">
        <div class="item-content">
          <h3 class="item-title">${this.escapeHtml(action.title)}</h3>
          <p class="item-description">${this.escapeHtml(description)}</p>
          <div class="item-meta">
            <span class="item-status" aria-label="상태: ${statusText}">
              ${statusText}
            </span>
            ${action.completedAt ? `
              <span class="completion-time">
                ${this.getRelativeTime(action.completedAt)} 완료
              </span>
            ` : ''}
          </div>
        </div>
        <div class="item-actions">
          ${this.renderActionButtons(action)}
        </div>
      </div>
    `;
  }

  /**
   * Action 버튼들 렌더링
   * @param {Object} action - Action 객체
   * @returns {string} 버튼 HTML
   */
  renderActionButtons(action) {
    switch (action.status) {
      case 'todo':
        return `
          <button type="button" class="btn btn-primary btn-start" 
                  data-action-id="${action.id}" aria-label="시작하기">
            시작
          </button>
          <button type="button" class="btn btn-secondary btn-defer" 
                  data-action-id="${action.id}" aria-label="연기하기">
            연기
          </button>
        `;
      case 'doing':
        return `
          <button type="button" class="btn btn-primary btn-complete" 
                  data-action-id="${action.id}" aria-label="완료하기">
            완료
          </button>
          <button type="button" class="btn btn-secondary btn-pause" 
                  data-action-id="${action.id}" aria-label="일시정지">
            일시정지
          </button>
        `;
      case 'done':
        return `
          <button type="button" class="btn btn-secondary btn-reopen" 
                  data-action-id="${action.id}" aria-label="다시 열기">
            다시 열기
          </button>
        `;
      default:
        return '';
    }
  }

  /**
   * Action 아이템 이벤트 핸들러 설정
   */
  setupActionHandlers() {
    const container = this.$('#action-list');
    if (!container) return;

    this.addEventListener(container, 'click', (event) => {
      const target = event.target;
      const actionId = parseInt(target.dataset.actionId);

      if (target.classList.contains('btn-start')) {
        this.handleActionStart(actionId);
      } else if (target.classList.contains('btn-complete')) {
        this.handleActionComplete(actionId);
      } else if (target.classList.contains('btn-pause')) {
        this.handleActionPause(actionId);
      } else if (target.classList.contains('btn-defer')) {
        this.handleActionDefer(actionId);
      } else if (target.classList.contains('btn-reopen')) {
        this.handleActionReopen(actionId);
      }
    });
  }

  /**
   * Action 시작 처리
   * @param {number} actionId - Action ID
   */
  async handleActionStart(actionId) {
    // 임시 구현
    this.toast.info('Action 시작 기능은 곧 구현될 예정입니다.');
  }

  /**
   * Action 완료 처리
   * @param {number} actionId - Action ID
   */
  async handleActionComplete(actionId) {
    // 임시 구현
    this.toast.info('Action 완료 기능은 곧 구현될 예정입니다.');
  }

  /**
   * Action 일시정지 처리
   * @param {number} actionId - Action ID
   */
  async handleActionPause(actionId) {
    // 임시 구현
    this.toast.info('Action 일시정지 기능은 곧 구현될 예정입니다.');
  }

  /**
   * Action 연기 처리
   * @param {number} actionId - Action ID
   */
  async handleActionDefer(actionId) {
    // 임시 구현
    this.toast.info('Action 연기 기능은 곧 구현될 예정입니다.');
  }

  /**
   * Action 다시 열기 처리
   * @param {number} actionId - Action ID
   */
  async handleActionReopen(actionId) {
    // 임시 구현
    this.toast.info('Action 다시 열기 기능은 곧 구현될 예정입니다.');
  }

  /**
   * 활성 필터 설정
   * @param {string} filter - 필터 타입
   */
  setActiveFilter(filter) {
    this.currentFilter = filter;
    
    // 필터 버튼 활성 상태 업데이트
    this.$$('.filter-btn').forEach(btn => {
      btn.classList.remove('active');
      if (btn.dataset.filter === filter) {
        btn.classList.add('active');
      }
    });
  }

  /**
   * 필터 뱃지 업데이트
   */
  updateFilterBadges() {
    const counts = {
      all: this.actions.length,
      todo: this.actions.filter(a => a.status === 'todo').length,
      doing: this.actions.filter(a => a.status === 'doing').length,
      done: this.actions.filter(a => a.status === 'done').length
    };

    Object.entries(counts).forEach(([filter, count]) => {
      const btn = this.$(`.filter-btn[data-filter=\"${filter}\"]`);
      if (btn) {
        const badge = btn.querySelector('.filter-badge');
        if (badge) badge.remove();
        
        if (count > 0) {
          const badgeEl = document.createElement('span');
          badgeEl.className = 'filter-badge';
          badgeEl.textContent = count;
          btn.appendChild(badgeEl);
        }
      }
    });
  }

  /**
   * 빈 상태 HTML 반환
   * @returns {string} 빈 상태 HTML
   */
  getEmptyStateHTML() {
    const filterText = {
      all: '아직 Action이 없습니다',
      todo: '해야 할 Action이 없습니다',
      doing: '진행 중인 Action이 없습니다',
      done: '완료된 Action이 없습니다'
    };

    return `
      <div class="empty-state">
        <h3>${filterText[this.currentFilter]}</h3>
        <p class="text-muted">
          ${this.currentFilter === 'all' 
            ? 'Inbox에서 생각을 명확화하여 Action을 만들어보세요.' 
            : '다른 필터를 선택해보세요.'}
        </p>
      </div>
    `;
  }

  /**
   * 상태 텍스트 반환
   * @param {string} status - 상태 코드
   * @returns {string} 상태 텍스트
   */
  getStatusText(status) {
    const statusMap = {
      todo: '해야됨',
      doing: '수행중',
      done: '완료',
      deferred: '연기됨',
      delegated: '위임됨'
    };
    return statusMap[status] || '알 수 없음';
  }

  /**
   * 상태별 CSS 클래스 반환
   * @param {string} status - 상태 코드
   * @returns {string} CSS 클래스
   */
  getStatusClass(status) {
    return `action-status-${status}`;
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