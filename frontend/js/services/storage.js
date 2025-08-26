/**
 * 로컬 스토리지 서비스
 * 클라이언트 사이드 데이터 저장 및 캐싱
 */
export class StorageService {
  constructor() {
    this.prefix = 'gtd_todo_';
    this.isSupported = this.checkSupport();
  }

  /**
   * 로컬 스토리지 지원 여부 확인
   * @returns {boolean} 지원 여부
   */
  checkSupport() {
    try {
      const test = '__storage_test__';
      localStorage.setItem(test, test);
      localStorage.removeItem(test);
      return true;
    } catch (error) {
      console.warn('로컬 스토리지가 지원되지 않습니다:', error);
      return false;
    }
  }

  /**
   * 키에 접두사 추가
   * @param {string} key - 원본 키
   * @returns {string} 접두사가 추가된 키
   */
  getKey(key) {
    return `${this.prefix}${key}`;
  }

  /**
   * 데이터 저장
   * @param {string} key - 저장 키
   * @param {any} value - 저장할 값
   * @returns {boolean} 저장 성공 여부
   */
  set(key, value) {
    if (!this.isSupported) return false;

    try {
      const serializedValue = JSON.stringify({
        data: value,
        timestamp: Date.now()
      });
      localStorage.setItem(this.getKey(key), serializedValue);
      return true;
    } catch (error) {
      console.error('데이터 저장 실패:', error);
      return false;
    }
  }

  /**
   * 데이터 조회
   * @param {string} key - 조회 키
   * @param {any} defaultValue - 기본값
   * @returns {any} 저장된 값 또는 기본값
   */
  get(key, defaultValue = null) {
    if (!this.isSupported) return defaultValue;

    try {
      const item = localStorage.getItem(this.getKey(key));
      if (item === null) return defaultValue;

      const parsed = JSON.parse(item);
      return parsed.data;
    } catch (error) {
      console.error('데이터 조회 실패:', error);
      return defaultValue;
    }
  }

  /**
   * 데이터 삭제
   * @param {string} key - 삭제 키
   * @returns {boolean} 삭제 성공 여부
   */
  remove(key) {
    if (!this.isSupported) return false;

    try {
      localStorage.removeItem(this.getKey(key));
      return true;
    } catch (error) {
      console.error('데이터 삭제 실패:', error);
      return false;
    }
  }

  /**
   * 모든 앱 데이터 삭제
   * @returns {boolean} 삭제 성공 여부
   */
  clear() {
    if (!this.isSupported) return false;

    try {
      const keys = Object.keys(localStorage);
      keys.forEach(key => {
        if (key.startsWith(this.prefix)) {
          localStorage.removeItem(key);
        }
      });
      return true;
    } catch (error) {
      console.error('데이터 전체 삭제 실패:', error);
      return false;
    }
  }

  /**
   * 캐시된 데이터 확인 (만료 시간 체크)
   * @param {string} key - 확인할 키
   * @param {number} maxAge - 최대 유효 시간 (밀리초)
   * @returns {boolean} 유효한 캐시 여부
   */
  isCacheValid(key, maxAge) {
    if (!this.isSupported) return false;

    try {
      const item = localStorage.getItem(this.getKey(key));
      if (!item) return false;

      const parsed = JSON.parse(item);
      const age = Date.now() - parsed.timestamp;
      
      return age < maxAge;
    } catch (error) {
      return false;
    }
  }

  /**
   * 만료된 캐시 데이터 정리
   * @param {number} maxAge - 최대 유효 시간 (밀리초)
   */
  cleanExpiredCache(maxAge = 24 * 60 * 60 * 1000) { // 기본 24시간
    if (!this.isSupported) return;

    try {
      const keys = Object.keys(localStorage);
      const now = Date.now();

      keys.forEach(key => {
        if (!key.startsWith(this.prefix)) return;

        try {
          const item = localStorage.getItem(key);
          const parsed = JSON.parse(item);
          const age = now - parsed.timestamp;

          if (age > maxAge) {
            localStorage.removeItem(key);
          }
        } catch (error) {
          // 파싱 실패한 항목은 삭제
          localStorage.removeItem(key);
        }
      });
    } catch (error) {
      console.error('만료된 캐시 정리 실패:', error);
    }
  }

  /**
   * 스토리지 사용량 확인 (추정치)
   * @returns {Object} 사용량 정보
   */
  getUsage() {
    if (!this.isSupported) {
      return { used: 0, total: 0, percentage: 0 };
    }

    try {
      let used = 0;
      const keys = Object.keys(localStorage);
      
      keys.forEach(key => {
        if (key.startsWith(this.prefix)) {
          used += (localStorage.getItem(key) || '').length;
        }
      });

      // LocalStorage 제한은 브라우저마다 다르지만 일반적으로 5-10MB
      const estimatedTotal = 5 * 1024 * 1024; // 5MB 추정
      const percentage = Math.round((used / estimatedTotal) * 100);

      return {
        used: used,
        total: estimatedTotal,
        percentage: Math.min(percentage, 100)
      };
    } catch (error) {
      return { used: 0, total: 0, percentage: 0 };
    }
  }

  // ==========================================================================
  // 앱별 특화 메서드들
  // ==========================================================================

  /**
   * 사용자 설정 저장
   * @param {Object} settings - 설정 객체
   */
  saveUserSettings(settings) {
    this.set('user_settings', settings);
  }

  /**
   * 사용자 설정 조회
   * @returns {Object} 사용자 설정
   */
  getUserSettings() {
    return this.get('user_settings', {
      theme: 'light',
      notifications: true,
      autoSave: true,
      keyboardShortcuts: true
    });
  }

  /**
   * 최근 활동 저장
   * @param {string} activity - 활동 내용
   */
  saveRecentActivity(activity) {
    const recent = this.get('recent_activities', []);
    recent.unshift({
      activity,
      timestamp: Date.now()
    });
    
    // 최대 50개까지만 보관
    const trimmed = recent.slice(0, 50);
    this.set('recent_activities', trimmed);
  }

  /**
   * 최근 활동 조회
   * @param {number} limit - 조회할 개수
   * @returns {Array} 최근 활동 목록
   */
  getRecentActivities(limit = 10) {
    const activities = this.get('recent_activities', []);
    return activities.slice(0, limit);
  }
}