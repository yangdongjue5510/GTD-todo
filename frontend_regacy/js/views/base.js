/**
 * 베이스 뷰 클래스
 * 모든 뷰 컴포넌트의 공통 기능 제공
 */
export class BaseView {
  constructor(services) {
    this.api = services.api;
    this.storage = services.storage;
    this.toast = services.toast;
    this.loading = services.loading;
    this.router = services.router;
    
    this.element = null;
    this.isVisible = false;
    this.listeners = [];
  }

  /**
   * 뷰 표시
   * @param {Object} params - 라우트 파라미터
   */
  async show(params = {}) {
    if (!this.element) {
      console.warn(`${this.constructor.name}: 뷰 엘리먼트가 설정되지 않았습니다.`);
      return;
    }

    this.isVisible = true;
    this.element.classList.remove('hidden');
    
    // 뷰별 표시 로직 실행
    await this.onShow(params);
    
    // 접근성: 포커스 관리
    this.manageFocus();
  }

  /**
   * 뷰 숨김
   */
  async hide() {
    if (!this.element || !this.isVisible) return;

    this.isVisible = false;
    this.element.classList.add('hidden');
    
    // 뷰별 숨김 로직 실행
    await this.onHide();
    
    // 이벤트 리스너 정리
    this.cleanupListeners();
  }

  /**
   * 뷰 표시 시 호출되는 메서드 (오버라이드 필요)
   * @param {Object} params - 라우트 파라미터
   */
  async onShow(params) {
    // 하위 클래스에서 구현
  }

  /**
   * 뷰 숨김 시 호출되는 메서드 (오버라이드 필요)
   */
  async onHide() {
    // 하위 클래스에서 구현
  }

  /**
   * 이벤트 리스너 추가 (자동 정리를 위해)
   * @param {Element} element - 이벤트 대상 엘리먼트
   * @param {string} event - 이벤트 타입
   * @param {Function} handler - 이벤트 핸들러
   * @param {Object} options - 이벤트 옵션
   */
  addEventListener(element, event, handler, options = {}) {
    element.addEventListener(event, handler, options);
    
    // 나중에 정리를 위해 저장
    this.listeners.push({
      element,
      event,
      handler,
      options
    });
  }

  /**
   * 등록된 모든 이벤트 리스너 정리
   */
  cleanupListeners() {
    this.listeners.forEach(({ element, event, handler }) => {
      element.removeEventListener(event, handler);
    });
    this.listeners = [];
  }

  /**
   * 포커스 관리 (접근성)
   */
  manageFocus() {
    // 뷰의 첫 번째 포커스 가능한 엘리먼트에 포커스
    const focusableElements = this.element.querySelectorAll(
      'button:not([disabled]), [href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex=\"-1\"])'
    );

    if (focusableElements.length > 0) {
      focusableElements[0].focus();
    }
  }

  /**
   * 폼 검증
   * @param {HTMLFormElement} form - 검증할 폼
   * @returns {Object} 검증 결과 { isValid, errors }
   */
  validateForm(form) {
    const errors = [];
    const formData = new FormData(form);
    
    // HTML5 기본 검증
    if (!form.checkValidity()) {
      const invalidElements = form.querySelectorAll(':invalid');
      invalidElements.forEach(element => {
        errors.push({
          element,
          message: element.validationMessage
        });
      });
    }

    return {
      isValid: errors.length === 0,
      errors,
      data: Object.fromEntries(formData)
    };
  }

  /**
   * 에러 처리 및 사용자에게 표시
   * @param {Error} error - 에러 객체
   * @param {string} context - 에러 발생 컨텍스트
   */
  handleError(error, context = '') {
    console.error(`${this.constructor.name} 에러 ${context}:`, error);
    
    const message = error.message || '알 수 없는 오류가 발생했습니다.';
    this.toast.error(message);
  }

  /**
   * 데이터 로딩 상태 처리
   * @param {Function} loadFunction - 로딩할 함수
   * @param {string} loadingMessage - 로딩 메시지
   * @returns {Promise} 로딩 결과
   */
  async loadWithState(loadFunction, loadingMessage = '데이터를 로드하는 중...') {
    try {
      return await this.loading.wrap(loadFunction(), loadingMessage);
    } catch (error) {
      this.handleError(error, '데이터 로딩 중');
      throw error;
    }
  }

  /**
   * 엘리먼트 찾기 (뷰 내에서)
   * @param {string} selector - CSS 선택자
   * @returns {Element|null} 찾은 엘리먼트
   */
  $(selector) {
    return this.element ? this.element.querySelector(selector) : null;
  }

  /**
   * 엘리먼트들 찾기 (뷰 내에서)
   * @param {string} selector - CSS 선택자
   * @returns {NodeList} 찾은 엘리먼트들
   */
  $$(selector) {
    return this.element ? this.element.querySelectorAll(selector) : [];
  }

  /**
   * HTML 문자열을 DOM 엘리먼트로 변환
   * @param {string} html - HTML 문자열
   * @returns {DocumentFragment} DOM 프래그먼트
   */
  createElementFromHTML(html) {
    const template = document.createElement('template');
    template.innerHTML = html.trim();
    return template.content;
  }

  /**
   * 디바운스된 함수 생성
   * @param {Function} func - 디바운스할 함수
   * @param {number} delay - 지연 시간 (밀리초)
   * @returns {Function} 디바운스된 함수
   */
  debounce(func, delay) {
    let timeoutId;
    return (...args) => {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => func.apply(this, args), delay);
    };
  }

  /**
   * 쓰로틀된 함수 생성
   * @param {Function} func - 쓰로틀할 함수
   * @param {number} limit - 제한 시간 (밀리초)
   * @returns {Function} 쓰로틀된 함수
   */
  throttle(func, limit) {
    let inThrottle;
    return (...args) => {
      if (!inThrottle) {
        func.apply(this, args);
        inThrottle = true;
        setTimeout(() => inThrottle = false, limit);
      }
    };
  }

  /**
   * 날짜 포맷팅
   * @param {Date|string|number} date - 포맷할 날짜
   * @param {string} format - 포맷 타입 ('short', 'long', 'time')
   * @returns {string} 포맷된 날짜 문자열
   */
  formatDate(date, format = 'short') {
    const dateObj = new Date(date);
    
    const options = {
      short: { year: 'numeric', month: 'short', day: 'numeric' },
      long: { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' },
      time: { hour: '2-digit', minute: '2-digit' }
    };

    return dateObj.toLocaleDateString('ko-KR', options[format] || options.short);
  }

  /**
   * 상대적 시간 표시 (예: "2분 전", "1시간 전")
   * @param {Date|string|number} date - 기준 날짜
   * @returns {string} 상대적 시간 문자열
   */
  getRelativeTime(date) {
    const now = new Date();
    const diff = now - new Date(date);
    
    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (seconds < 60) return '방금 전';
    if (minutes < 60) return `${minutes}분 전`;
    if (hours < 24) return `${hours}시간 전`;
    if (days < 7) return `${days}일 전`;
    
    return this.formatDate(date);
  }
}