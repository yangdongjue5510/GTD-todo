/**
 * 로딩 인디케이터 서비스
 * 비동기 작업 중 로딩 상태 관리
 */
export class LoadingService {
  constructor() {
    this.container = document.getElementById('loading');
    this.isLoading = false;
    this.loadingCount = 0;
  }

  /**
   * 로딩 상태 표시
   */
  show() {
    this.loadingCount++;
    
    if (!this.isLoading && this.container) {
      this.isLoading = true;
      this.container.setAttribute('aria-hidden', 'false');
      // 접근성을 위한 포커스 트랩
      this.container.setAttribute('role', 'dialog');
      this.container.setAttribute('aria-label', '로딩 중');
    }
  }

  /**
   * 로딩 상태 숨김
   */
  hide() {
    this.loadingCount = Math.max(0, this.loadingCount - 1);
    
    if (this.loadingCount === 0 && this.isLoading && this.container) {
      this.isLoading = false;
      this.container.setAttribute('aria-hidden', 'true');
      this.container.removeAttribute('role');
      this.container.removeAttribute('aria-label');
    }
  }

  /**
   * 비동기 작업을 로딩 상태로 감싸기
   * @param {Promise} promise - 비동기 작업 프로미스
   * @param {string} message - 로딩 메시지 (선택사항)
   * @returns {Promise} 원본 프로미스
   */
  async wrap(promise, message) {
    this.show();
    
    if (message) {
      // 로딩 메시지가 있다면 스피너 아래에 표시
      const spinner = this.container?.querySelector('.loading-spinner');
      if (spinner) {
        let messageElement = this.container.querySelector('.loading-message');
        if (!messageElement) {
          messageElement = document.createElement('div');
          messageElement.className = 'loading-message';
          messageElement.style.marginTop = '1rem';
          messageElement.style.color = 'var(--color-text-secondary)';
          messageElement.style.fontSize = 'var(--font-size-sm)';
          this.container.appendChild(messageElement);
        }
        messageElement.textContent = message;
      }
    }

    try {
      const result = await promise;
      return result;
    } finally {
      this.hide();
      // 메시지 제거
      const messageElement = this.container?.querySelector('.loading-message');
      if (messageElement) {
        messageElement.remove();
      }
    }
  }

  /**
   * 현재 로딩 상태 확인
   * @returns {boolean} 로딩 중 여부
   */
  isCurrentlyLoading() {
    return this.isLoading;
  }

  /**
   * 모든 로딩 상태 강제 종료
   */
  forceHide() {
    this.loadingCount = 0;
    this.hide();
  }
}