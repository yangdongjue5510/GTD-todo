/**
 * API 서비스
 * 백엔드와의 통신을 담당
 */
export class ApiService {
  constructor() {
    this.baseURL = 'http://localhost:8080';
    this.defaultHeaders = {
      'Content-Type': 'application/json'
    };
  }

  /**
   * HTTP 요청 수행
   * @param {string} endpoint - API 엔드포인트
   * @param {Object} options - 요청 옵션
   * @returns {Promise} 응답 프로미스
   */
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    
    const config = {
      headers: { ...this.defaultHeaders, ...options.headers },
      ...options
    };

    try {
      const response = await fetch(url, config);
      
      // HTTP 에러 체크
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      // 응답이 JSON인지 확인
      const contentType = response.headers.get('content-type');
      if (contentType && contentType.includes('application/json')) {
        return await response.json();
      }
      
      return response;
    } catch (error) {
      console.error(`API 요청 실패: ${endpoint}`, error);
      throw this.handleApiError(error);
    }
  }

  /**
   * GET 요청
   * @param {string} endpoint - API 엔드포인트
   * @param {Object} params - 쿼리 파라미터
   * @returns {Promise} 응답 프로미스
   */
  async get(endpoint, params = {}) {
    const queryString = new URLSearchParams(params).toString();
    const url = queryString ? `${endpoint}?${queryString}` : endpoint;
    
    return this.request(url, { method: 'GET' });
  }

  /**
   * POST 요청
   * @param {string} endpoint - API 엔드포인트
   * @param {Object} data - 요청 데이터
   * @returns {Promise} 응답 프로미스
   */
  async post(endpoint, data = {}) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  /**
   * PUT 요청
   * @param {string} endpoint - API 엔드포인트
   * @param {Object} data - 요청 데이터
   * @returns {Promise} 응답 프로미스
   */
  async put(endpoint, data = {}) {
    return this.request(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data)
    });
  }

  /**
   * DELETE 요청
   * @param {string} endpoint - API 엔드포인트
   * @returns {Promise} 응답 프로미스
   */
  async delete(endpoint) {
    return this.request(endpoint, { method: 'DELETE' });
  }

  /**
   * API 에러 처리
   * @param {Error} error - 에러 객체
   * @returns {Error} 처리된 에러 객체
   */
  handleApiError(error) {
    if (error.name === 'TypeError' && error.message.includes('fetch')) {
      return new Error('서버에 연결할 수 없습니다. 네트워크 연결을 확인해주세요.');
    }
    
    if (error.message.includes('HTTP 404')) {
      return new Error('요청한 리소스를 찾을 수 없습니다.');
    }
    
    if (error.message.includes('HTTP 500')) {
      return new Error('서버 내부 오류가 발생했습니다.');
    }
    
    return error;
  }

  // ==========================================================================
  // Thing API 메서드들
  // ==========================================================================

  /**
   * 모든 Thing 조회
   * @returns {Promise<Array>} Thing 목록
   */
  async getThings() {
    return this.get('/things/');
  }

  /**
   * Thing 생성
   * @param {Object} thing - Thing 데이터
   * @returns {Promise<Object>} 생성된 Thing
   */
  async createThing(thing) {
    return this.post('/things/', {
      title: thing.title,
      description: thing.description || '',
      status: 0 // Pending 상태
    });
  }

  /**
   * Thing 상태 변경
   * @param {number} thingId - Thing ID
   * @param {number} status - 새로운 상태
   * @returns {Promise<Object>} 업데이트된 Thing
   */
  async updateThingStatus(thingId, status) {
    return this.put(`/things/${thingId}`, { status });
  }

  /**
   * Thing 삭제
   * @param {number} thingId - Thing ID
   * @returns {Promise} 삭제 결과
   */
  async deleteThing(thingId) {
    return this.delete(`/things/${thingId}`);
  }

  // ==========================================================================
  // Action API 메서드들 (향후 구현)
  // ==========================================================================

  /**
   * 모든 Action 조회
   * @returns {Promise<Array>} Action 목록
   */
  async getActions() {
    return this.get('/actions/');
  }

  /**
   * Action 생성
   * @param {Object} action - Action 데이터
   * @returns {Promise<Object>} 생성된 Action
   */
  async createAction(action) {
    return this.post('/actions/', action);
  }

  /**
   * Action 상태 변경
   * @param {number} actionId - Action ID
   * @param {number} status - 새로운 상태
   * @returns {Promise<Object>} 업데이트된 Action
   */
  async updateActionStatus(actionId, status) {
    return this.put(`/actions/${actionId}`, { status });
  }

  // ==========================================================================
  // Project API 메서드들 (향후 구현)
  // ==========================================================================

  /**
   * 모든 Project 조회
   * @returns {Promise<Array>} Project 목록
   */
  async getProjects() {
    return this.get('/projects/');
  }

  /**
   * Project 생성
   * @param {Object} project - Project 데이터
   * @returns {Promise<Object>} 생성된 Project
   */
  async createProject(project) {
    return this.post('/projects/', project);
  }
}