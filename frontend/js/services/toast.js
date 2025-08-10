/**
 * 토스트 알림 서비스
 * 사용자에게 피드백 메시지 표시
 */
export class ToastService {
  constructor() {
    this.container = document.getElementById('toast');
    this.hideTimeout = null;
  }

  /**
   * 토스트 메시지 표시
   * @param {string} message - 표시할 메시지
   * @param {string} type - 메시지 타입 ('success', 'error', 'warning', 'info')
   * @param {number} duration - 표시 시간 (ms), 기본값 3초
   */
  show(message, type = 'info', duration = 3000) {
    if (!this.container) {
      console.warn('토스트 컨테이너를 찾을 수 없습니다.');
      return;
    }

    // 이전 타이머 정리
    if (this.hideTimeout) {
      clearTimeout(this.hideTimeout);
      this.hideTimeout = null;
    }

    // 메시지 설정
    this.container.textContent = message;
    this.container.className = `toast toast-${type}`;

    // 표시
    this.container.classList.add('show');
    this.container.setAttribute('aria-hidden', 'false');

    // 자동 숨김 설정
    if (duration > 0) {
      this.hideTimeout = setTimeout(() => {
        this.hide();
      }, duration);
    }

    // 클릭시 숨김
    this.container.onclick = () => this.hide();
  }

  /**
   * 토스트 메시지 숨김
   */
  hide() {
    if (!this.container) return;

    this.container.classList.remove('show');
    this.container.setAttribute('aria-hidden', 'true');
    
    if (this.hideTimeout) {
      clearTimeout(this.hideTimeout);
      this.hideTimeout = null;
    }
  }

  /**
   * 성공 메시지 표시
   * @param {string} message - 메시지
   * @param {number} duration - 표시 시간
   */
  success(message, duration = 3000) {
    this.show(message, 'success', duration);
  }

  /**
   * 에러 메시지 표시
   * @param {string} message - 메시지
   * @param {number} duration - 표시 시간 (에러는 기본적으로 더 오래 표시)
   */
  error(message, duration = 5000) {
    this.show(message, 'error', duration);
  }

  /**
   * 경고 메시지 표시
   * @param {string} message - 메시지
   * @param {number} duration - 표시 시간
   */
  warning(message, duration = 4000) {
    this.show(message, 'warning', duration);
  }

  /**
   * 정보 메시지 표시
   * @param {string} message - 메시지
   * @param {number} duration - 표시 시간
   */
  info(message, duration = 3000) {
    this.show(message, 'info', duration);
  }
}