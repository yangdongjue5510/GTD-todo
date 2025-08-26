/**
 * 네비게이션 컴포넌트
 * 메인 네비게이션 UI 및 상태 관리
 */
export class Navigation {
  constructor(router) {
    this.router = router;
    this.container = document.querySelector('.main-nav');
    this.currentActive = null;
    
    this.init();
  }

  /**
   * 네비게이션 초기화
   */
  init() {
    if (!this.container) {
      console.warn('네비게이션 컨테이너를 찾을 수 없습니다.');
      return;
    }

    this.setupEventListeners();
    this.setupRouterListeners();
  }

  /**
   * 이벤트 리스너 설정
   */
  setupEventListeners() {
    // 네비게이션 버튼 클릭 처리
    this.container.addEventListener('click', (event) => {
      const navItem = event.target.closest('.nav-item');
      if (!navItem) return;

      event.preventDefault();
      
      const view = navItem.dataset.view;
      if (view) {
        this.router.navigate(view);
      }
    });

    // 키보드 네비게이션 지원
    this.container.addEventListener('keydown', (event) => {
      this.handleKeyboardNavigation(event);
    });
  }

  /**
   * 라우터 이벤트 리스너 설정
   */
  setupRouterListeners() {
    this.router.on('routeChange', (route) => {
      this.setActiveItem(route);
    });
  }

  /**
   * 키보드 네비게이션 처리
   * @param {KeyboardEvent} event - 키보드 이벤트
   */
  handleKeyboardNavigation(event) {
    const navItems = Array.from(this.container.querySelectorAll('.nav-item'));
    const currentIndex = navItems.indexOf(document.activeElement);

    switch (event.key) {
      case 'ArrowLeft':
      case 'ArrowUp':
        event.preventDefault();
        const prevIndex = currentIndex > 0 ? currentIndex - 1 : navItems.length - 1;
        navItems[prevIndex].focus();
        break;

      case 'ArrowRight':
      case 'ArrowDown':
        event.preventDefault();
        const nextIndex = currentIndex < navItems.length - 1 ? currentIndex + 1 : 0;
        navItems[nextIndex].focus();
        break;

      case 'Home':
        event.preventDefault();
        navItems[0].focus();
        break;

      case 'End':
        event.preventDefault();
        navItems[navItems.length - 1].focus();
        break;

      case 'Enter':
      case ' ':
        event.preventDefault();
        if (document.activeElement && document.activeElement.classList.contains('nav-item')) {
          document.activeElement.click();
        }
        break;
    }
  }

  /**
   * 활성 네비게이션 아이템 설정
   * @param {string} route - 활성화할 라우트
   */
  setActiveItem(route) {
    // 이전 활성 아이템 비활성화
    if (this.currentActive) {
      this.currentActive.classList.remove('active');
      this.currentActive.setAttribute('aria-current', 'false');
    }

    // 새로운 활성 아이템 찾기 및 활성화
    const newActive = this.container.querySelector(`[data-view=\"${route}\"]`);
    if (newActive) {
      newActive.classList.add('active');
      newActive.setAttribute('aria-current', 'page');
      this.currentActive = newActive;
    }

    // 접근성을 위한 알림
    this.announceNavigation(route);
  }

  /**
   * 네비게이션 변경 알림 (스크린 리더용)
   * @param {string} route - 변경된 라우트
   */
  announceNavigation(route) {
    const routeNames = {
      inbox: 'Inbox - 생각 포착',
      actions: 'Actions - 할 일 관리',
      projects: 'Projects - 프로젝트 관리',
      review: 'Review - 진행 상황 검토'
    };

    const routeName = routeNames[route] || route;
    
    // aria-live 영역에 알림
    let announcer = document.getElementById('nav-announcer');
    if (!announcer) {
      announcer = document.createElement('div');
      announcer.id = 'nav-announcer';
      announcer.className = 'sr-only';
      announcer.setAttribute('aria-live', 'polite');
      document.body.appendChild(announcer);
    }

    announcer.textContent = `${routeName} 페이지로 이동했습니다.`;
    
    // 알림 후 내용 정리
    setTimeout(() => {
      announcer.textContent = '';
    }, 1000);
  }

  /**
   * 네비게이션 아이템에 뱃지 추가
   * @param {string} route - 대상 라우트
   * @param {number|string} badge - 뱃지 내용
   */
  setBadge(route, badge) {
    const navItem = this.container.querySelector(`[data-view=\"${route}\"]`);
    if (!navItem) return;

    // 기존 뱃지 제거
    const existingBadge = navItem.querySelector('.nav-badge');
    if (existingBadge) {
      existingBadge.remove();
    }

    // 새 뱃지 추가 (값이 0이 아닌 경우에만)
    if (badge && badge !== 0) {
      const badgeElement = document.createElement('span');
      badgeElement.className = 'nav-badge';
      badgeElement.textContent = badge;
      badgeElement.setAttribute('aria-label', `${badge}개의 새로운 항목`);
      
      // 뱃지 스타일 (CSS에서 정의되어야 함)
      badgeElement.style.cssText = `
        position: absolute;
        top: -4px;
        right: -4px;
        background: var(--color-danger);
        color: white;
        border-radius: 50%;
        width: 18px;
        height: 18px;
        font-size: 11px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: var(--font-weight-semibold);
      `;
      
      navItem.style.position = 'relative';
      navItem.appendChild(badgeElement);
    }
  }

  /**
   * 특정 네비게이션 아이템 비활성화/활성화
   * @param {string} route - 대상 라우트
   * @param {boolean} disabled - 비활성화 여부
   */
  setDisabled(route, disabled) {
    const navItem = this.container.querySelector(`[data-view=\"${route}\"]`);
    if (!navItem) return;

    if (disabled) {
      navItem.setAttribute('disabled', 'true');
      navItem.setAttribute('aria-disabled', 'true');
      navItem.style.opacity = '0.5';
      navItem.style.pointerEvents = 'none';
    } else {
      navItem.removeAttribute('disabled');
      navItem.removeAttribute('aria-disabled');
      navItem.style.opacity = '';
      navItem.style.pointerEvents = '';
    }
  }

  /**
   * 현재 활성 라우트 반환
   * @returns {string|null} 활성 라우트
   */
  getActiveRoute() {
    return this.currentActive ? this.currentActive.dataset.view : null;
  }
}