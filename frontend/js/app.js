/**
 * GTD-TODO 메인 애플리케이션
 * ES6 모듈과 최신 JavaScript 베스트 프랙티스 사용
 */

import { Router } from './utils/router.js';
import { ApiService } from './services/api.js';
import { StorageService } from './services/storage.js';
import { ToastService } from './services/toast.js';
import { LoadingService } from './services/loading.js';

import { InboxView } from './views/inbox.js';
import { ClarifyView } from './views/clarify.js';
import { ActionsView } from './views/actions.js';
import { ProjectsView } from './views/projects.js';
import { ReviewView } from './views/review.js';

import { Navigation } from './components/navigation.js';

/**
 * 메인 애플리케이션 클래스
 */
class GTDApp {
  constructor() {
    this.initializeServices();
    this.initializeViews();
    this.initializeComponents();
    this.setupEventListeners();
    this.handleInitialRoute();
  }

  /**
   * 서비스 초기화
   */
  initializeServices() {
    this.apiService = new ApiService();
    this.storageService = new StorageService();
    this.toastService = new ToastService();
    this.loadingService = new LoadingService();
    this.router = new Router();
  }

  /**
   * 뷰 초기화
   */
  initializeViews() {
    const services = {
      api: this.apiService,
      storage: this.storageService,
      toast: this.toastService,
      loading: this.loadingService,
      router: this.router
    };

    this.views = {
      inbox: new InboxView(services),
      clarify: new ClarifyView(services),
      actions: new ActionsView(services),
      projects: new ProjectsView(services),
      review: new ReviewView(services)
    };
  }

  /**
   * 컴포넌트 초기화
   */
  initializeComponents() {
    this.navigation = new Navigation(this.router);
  }

  /**
   * 이벤트 리스너 설정
   */
  setupEventListeners() {
    // 라우터 변경 이벤트
    this.router.on('routeChange', (route, params) => {
      this.handleRouteChange(route, params);
    });

    // 전역 에러 처리
    window.addEventListener('error', (event) => {
      console.error('Global error:', event.error);
      this.toastService.show('오류가 발생했습니다. 페이지를 새로고침해주세요.', 'error');
    });

    // 네트워크 상태 변경
    window.addEventListener('online', () => {
      this.toastService.show('인터넷 연결이 복구되었습니다.', 'success');
    });

    window.addEventListener('offline', () => {
      this.toastService.show('인터넷 연결이 끊어졌습니다.', 'warning');
    });

    // 키보드 단축키
    document.addEventListener('keydown', this.handleKeyboardShortcuts.bind(this));
  }

  /**
   * 라우트 변경 처리
   * @param {string} route - 새로운 라우트
   * @param {Object} params - 라우트 파라미터
   */
  handleRouteChange(route, params) {
    // 모든 뷰 숨기기
    Object.values(this.views).forEach(view => view.hide());

    // 해당 뷰 표시
    if (this.views[route]) {
      this.views[route].show(params);
    } else {
      // 기본값으로 inbox 표시
      this.router.navigate('inbox');
    }
  }

  /**
   * 초기 라우트 처리
   */
  handleInitialRoute() {
    const currentPath = window.location.hash.slice(1) || 'inbox';
    this.router.navigate(currentPath);
  }

  /**
   * 키보드 단축키 처리
   * @param {KeyboardEvent} event - 키보드 이벤트
   */
  handleKeyboardShortcuts(event) {
    // Ctrl/Cmd 키와 함께 눌러진 경우만 처리
    if (!event.ctrlKey && !event.metaKey) return;

    switch (event.key) {
      case '1':
        event.preventDefault();
        this.router.navigate('inbox');
        break;
      case '2':
        event.preventDefault();
        this.router.navigate('actions');
        break;
      case '3':
        event.preventDefault();
        this.router.navigate('projects');
        break;
      case '4':
        event.preventDefault();
        this.router.navigate('review');
        break;
      case 'n':
        event.preventDefault();
        // Inbox로 이동하고 입력 필드에 포커스
        this.router.navigate('inbox');
        setTimeout(() => {
          const titleInput = document.getElementById('thing-title');
          if (titleInput) titleInput.focus();
        }, 100);
        break;
    }
  }
}

/**
 * DOM 로드 완료 시 앱 초기화
 */
document.addEventListener('DOMContentLoaded', () => {
  console.log('DOM 로드 완료, GTD-TODO 앱 초기화 시작...');
  
  try {
    console.log('GTDApp 클래스 생성 중...');
    const app = new GTDApp();
    console.log('GTD-TODO 앱이 성공적으로 초기화되었습니다.');
    window.gtdApp = app; // 디버깅용 전역 접근
  } catch (error) {
    console.error('앱 초기화 중 오류가 발생했습니다:', error);
    console.error('Error stack:', error.stack);
    
    // 에러 정보를 화면에 표시
    const errorDiv = document.createElement('div');
    errorDiv.innerHTML = `
      <div style="padding: 2rem; text-align: center; color: #ef4444; background: #fee; border: 1px solid #fcc; margin: 1rem; border-radius: 8px;">
        <h1>🚨 앱 초기화 오류</h1>
        <p><strong>오류:</strong> ${error.message}</p>
        <details>
          <summary>상세 정보</summary>
          <pre style="text-align: left; background: #f8f8f8; padding: 1rem; margin: 1rem 0; border-radius: 4px;">${error.stack}</pre>
        </details>
        <p>브라우저 개발자 도구 콘솔을 확인해주세요.</p>
      </div>
    `;
    document.body.appendChild(errorDiv);
  }
});

// 개발 환경에서 디버깅을 위한 전역 접근
if (typeof window !== 'undefined') {
  window.GTDApp = GTDApp;
}