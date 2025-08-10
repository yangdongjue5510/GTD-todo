/**
 * 클라이언트 사이드 라우터
 * 해시 기반 라우팅 구현
 */
export class Router {
  constructor() {
    this.routes = new Map();
    this.currentRoute = '';
    this.listeners = new Map();
    
    // 브라우저 뒤로가기/앞으로가기 처리
    window.addEventListener('hashchange', this.handleHashChange.bind(this));
  }

  /**
   * 이벤트 리스너 등록
   * @param {string} event - 이벤트 이름
   * @param {Function} callback - 콜백 함수
   */
  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event).push(callback);
  }

  /**
   * 이벤트 발생
   * @param {string} event - 이벤트 이름
   * @param {...any} args - 전달할 인자들
   */
  emit(event, ...args) {
    const callbacks = this.listeners.get(event);
    if (callbacks) {
      callbacks.forEach(callback => callback(...args));
    }
  }

  /**
   * 라우트 등록
   * @param {string} path - 라우트 경로
   * @param {Function} handler - 핸들러 함수
   */
  register(path, handler) {
    this.routes.set(path, handler);
  }

  /**
   * 라우트로 이동
   * @param {string} path - 이동할 경로
   * @param {Object} params - 파라미터
   */
  navigate(path, params = {}) {
    const fullPath = params && Object.keys(params).length > 0 
      ? `${path}?${new URLSearchParams(params).toString()}` 
      : path;
    
    window.location.hash = fullPath;
  }

  /**
   * 해시 변경 처리
   */
  handleHashChange() {
    const hash = window.location.hash.slice(1);
    const [path, queryString] = hash.split('?');
    const params = this.parseQueryString(queryString);
    
    this.currentRoute = path || 'inbox';
    this.emit('routeChange', this.currentRoute, params);
  }

  /**
   * 쿼리 스트링 파싱
   * @param {string} queryString - 쿼리 스트링
   * @returns {Object} 파싱된 파라미터 객체
   */
  parseQueryString(queryString) {
    if (!queryString) return {};
    
    const params = {};
    const pairs = queryString.split('&');
    
    pairs.forEach(pair => {
      const [key, value] = pair.split('=');
      if (key) {
        params[decodeURIComponent(key)] = decodeURIComponent(value || '');
      }
    });
    
    return params;
  }

  /**
   * 현재 라우트 반환
   * @returns {string} 현재 라우트
   */
  getCurrentRoute() {
    return this.currentRoute;
  }

  /**
   * 뒤로 가기
   */
  back() {
    window.history.back();
  }

  /**
   * 앞으로 가기
   */
  forward() {
    window.history.forward();
  }
}