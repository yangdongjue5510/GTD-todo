/**
 * Inbox 뷰
 * Thing 포착 및 목록 관리
 */
import { BaseView } from './base.js';

export class InboxView extends BaseView {
  constructor(services) {
    super(services);
    this.element = document.getElementById('inbox-view');
    this.things = [];
    this.init();
  }

  /**
   * 뷰 초기화
   */
  init() {
    if (!this.element) return;

    this.setupFormHandlers();
  }

  /**
   * 폼 이벤트 핸들러 설정
   */
  setupFormHandlers() {
    const form = this.$('#capture-form');
    const titleInput = this.$('#thing-title');
    const descriptionInput = this.$('#thing-description');

    if (form && titleInput) {
      // 폼 제출 처리
      this.addEventListener(form, 'submit', this.handleFormSubmit.bind(this));

      // 엔터키 제출 (Shift+Enter로 줄바꿈 가능하도록)
      this.addEventListener(titleInput, 'keydown', (event) => {
        if (event.key === 'Enter' && !event.shiftKey) {
          event.preventDefault();
          form.dispatchEvent(new Event('submit'));
        }
      });

      // 자동 저장 (타이핑 중단 후 2초)
      const autoSave = this.debounce(() => {
        this.saveInputState();
      }, 2000);

      this.addEventListener(titleInput, 'input', autoSave);
      this.addEventListener(descriptionInput, 'input', autoSave);

      // 페이지 로드 시 저장된 입력 복원
      this.restoreInputState();
    }
  }

  /**
   * 뷰 표시 시 호출
   */
  async onShow(params) {
    await this.loadThings();
    
    // URL 파라미터로 자동 포커스 제어
    if (params.focus === 'input') {
      const titleInput = this.$('#thing-title');
      if (titleInput) titleInput.focus();
    }
  }

  /**
   * Thing 목록 로드
   */
  async loadThings() {
    try {
      this.things = await this.loadWithState(
        () => this.api.getThings(),
        'Thing 목록을 불러오는 중...'
      );
      
      this.renderThings();
      this.storage.saveRecentActivity('Inbox 새로고침');
    } catch (error) {
      this.handleError(error, 'Thing 목록 로드');
    }
  }

  /**
   * Thing 목록 렌더링
   */
  renderThings() {
    const container = this.$('#thing-list');
    if (!container) return;

    // 빈 상태 처리
    if (this.things.length === 0) {
      container.innerHTML = this.getEmptyStateHTML();
      return;
    }

    // Thing 목록 렌더링
    const thingsHTML = this.things
      .sort((a, b) => b.id - a.id) // 최신순 정렬
      .map(thing => this.renderThingItem(thing))
      .join('');

    container.innerHTML = thingsHTML;

    // 동적으로 생성된 버튼들에 이벤트 리스너 추가
    this.setupThingItemHandlers();
  }

  /**
   * 개별 Thing 아이템 렌더링
   * @param {Object} thing - Thing 객체
   * @returns {string} HTML 문자열
   */
  renderThingItem(thing) {
    const statusText = this.getStatusText(thing.status);
    const statusClass = this.getStatusClass(thing.status);
    const description = thing.description || '설명 없음';
    
    return `
      <div class="thing-item" data-thing-id="${thing.id}" role="listitem">
        <h3 class="item-title">${this.escapeHtml(thing.title)}</h3>
        <p class="item-description">${this.escapeHtml(description)}</p>
        <div class="item-meta">
          <span class="item-status ${statusClass}" aria-label="상태: ${statusText}">
            ${statusText}
          </span>
          <div class="item-actions">
            ${thing.status === 0 ? `
              <button type="button" class="btn btn-primary btn-clarify" 
                      data-thing-id="${thing.id}"
                      aria-label="${thing.title} 명확화">
                명확화
              </button>
            ` : ''}
            <button type="button" class="btn btn-secondary btn-delete" 
                    data-thing-id="${thing.id}"
                    aria-label="${thing.title} 삭제">
              제거
            </button>
          </div>
        </div>
      </div>
    `;
  }

  /**
   * Thing 아이템 이벤트 핸들러 설정
   */
  setupThingItemHandlers() {
    const container = this.$('#thing-list');
    if (!container) return;

    // 이벤트 위임을 통한 효율적인 핸들러 등록
    this.addEventListener(container, 'click', (event) => {
      const target = event.target;
      const thingId = parseInt(target.dataset.thingId);

      if (target.classList.contains('btn-clarify')) {
        this.handleClarifyClick(thingId);
      } else if (target.classList.contains('btn-delete')) {
        this.handleDeleteClick(thingId);
      }
    });
  }

  /**
   * 폼 제출 처리
   * @param {Event} event - 제출 이벤트
   */
  async handleFormSubmit(event) {
    event.preventDefault();
    
    const form = event.target;
    const validation = this.validateForm(form);
    
    if (!validation.isValid) {
      const firstError = validation.errors[0];
      this.toast.error(firstError.message);
      firstError.element.focus();
      return;
    }

    const { title, description } = validation.data;
    
    try {
      await this.loadWithState(
        () => this.api.createThing({ title, description: description || '' }),
        '새로운 생각을 추가하는 중...'
      );

      // 폼 초기화
      form.reset();
      this.clearInputState();
      
      // 목록 새로고침
      await this.loadThings();
      
      this.toast.success('새로운 생각이 추가되었습니다.');
      this.storage.saveRecentActivity(`새 생각 추가: ${title}`);
      
      // 다음 입력을 위해 포커스
      const titleInput = this.$('#thing-title');
      if (titleInput) titleInput.focus();
      
    } catch (error) {
      this.handleError(error, '생각 추가');
    }
  }

  /**
   * 명확화 버튼 클릭 처리
   * @param {number} thingId - Thing ID
   */
  handleClarifyClick(thingId) {
    const thing = this.things.find(t => t.id === thingId);
    if (!thing) return;

    // Clarify 뷰로 이동
    this.router.navigate('clarify', { thingId });
    this.storage.saveRecentActivity(`명확화 시작: ${thing.title}`);
  }

  /**
   * 삭제 버튼 클릭 처리
   * @param {number} thingId - Thing ID
   */
  async handleDeleteClick(thingId) {
    const thing = this.things.find(t => t.id === thingId);
    if (!thing) return;

    // 확인 대화상자
    const confirmed = confirm(`\"${thing.title}\"을(를) 정말 삭제하시겠습니까?`);
    if (!confirmed) return;

    try {
      await this.loadWithState(
        () => this.api.deleteThing(thingId),
        '삭제하는 중...'
      );

      await this.loadThings();
      this.toast.success('생각이 삭제되었습니다.');
      this.storage.saveRecentActivity(`생각 삭제: ${thing.title}`);
      
    } catch (error) {
      this.handleError(error, '생각 삭제');
    }
  }

  /**
   * 입력 상태 저장 (자동 저장)
   */
  saveInputState() {
    const titleInput = this.$('#thing-title');
    const descriptionInput = this.$('#thing-description');
    
    if (titleInput && descriptionInput) {
      this.storage.set('inbox_draft', {
        title: titleInput.value,
        description: descriptionInput.value
      });
    }
  }

  /**
   * 저장된 입력 상태 복원
   */
  restoreInputState() {
    const draft = this.storage.get('inbox_draft');
    if (!draft) return;

    const titleInput = this.$('#thing-title');
    const descriptionInput = this.$('#thing-description');
    
    if (titleInput && draft.title) {
      titleInput.value = draft.title;
    }
    
    if (descriptionInput && draft.description) {
      descriptionInput.value = draft.description;
    }
  }

  /**
   * 입력 상태 정리
   */
  clearInputState() {
    this.storage.remove('inbox_draft');
  }

  /**
   * 빈 상태 HTML 반환
   * @returns {string} 빈 상태 HTML
   */
  getEmptyStateHTML() {
    return `
      <div class="empty-state">
        <h3>아직 포착된 생각이 없습니다</h3>
        <p class="text-muted">
          머릿속에 떠오르는 모든 생각을 자유롭게 입력해보세요.<br>
          나중에 명확화 과정을 통해 정리할 수 있습니다.
        </p>
        <p class="text-muted">
          <strong>팁:</strong> Ctrl+N (또는 Cmd+N)으로 빠르게 입력할 수 있습니다.
        </p>
      </div>
    `;
  }

  /**
   * 상태 텍스트 반환
   * @param {number} status - 상태 코드
   * @returns {string} 상태 텍스트
   */
  getStatusText(status) {
    const statusMap = {
      0: '대기',
      1: '연기',
      2: '완료'
    };
    return statusMap[status] || '알 수 없음';
  }

  /**
   * 상태별 CSS 클래스 반환
   * @param {number} status - 상태 코드
   * @returns {string} CSS 클래스
   */
  getStatusClass(status) {
    const classMap = {
      0: 'status-pending',
      1: 'status-someday',
      2: 'status-done'
    };
    return classMap[status] || '';
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