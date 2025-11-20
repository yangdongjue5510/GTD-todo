/**
 * Review 뷰
 * 진행 상황 검토 및 리포트
 */
import { BaseView } from './base.js';

export class ReviewView extends BaseView {
  constructor(services) {
    super(services);
    this.element = document.getElementById('review-view');
    this.reviewData = null;
    this.init();
  }

  /**
   * 뷰 초기화
   */
  init() {
    if (!this.element) return;
    // 현재는 기본 구현만
  }

  /**
   * 뷰 표시 시 호출
   */
  async onShow(params) {
    await this.loadReviewData();
  }

  /**
   * 리뷰 데이터 로드
   */
  async loadReviewData() {
    try {
      // 기본 통계 데이터 수집
      const things = await this.loadWithState(
        () => this.api.getThings(),
        '리뷰 데이터를 준비하는 중...'
      );

      this.reviewData = this.generateReviewData(things);
      this.renderReviewContent();
      this.storage.saveRecentActivity('Review 확인');
      
    } catch (error) {
      this.handleError(error, '리뷰 데이터 로드');
    }
  }

  /**
   * 리뷰 데이터 생성
   * @param {Array} things - Thing 목록
   * @returns {Object} 리뷰 데이터
   */
  generateReviewData(things) {
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const weekAgo = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);

    return {
      total: things.length,
      pending: things.filter(t => t.status === 0).length,
      someday: things.filter(t => t.status === 1).length,
      done: things.filter(t => t.status === 2).length,
      recentActivities: this.storage.getRecentActivities(5),
      weeklyStats: {
        thisWeek: things.filter(t => new Date(t.createdAt || now) >= weekAgo).length,
        lastWeek: 0 // 실제로는 지난주 데이터가 필요
      }
    };
  }

  /**
   * 리뷰 콘텐츠 렌더링
   */
  renderReviewContent() {
    const container = this.$('#review-content');
    if (!container || !this.reviewData) return;

    container.innerHTML = `
      <div class="review-summary">
        <h3>📊 현재 상황</h3>
        
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-number">${this.reviewData.total}</div>
            <div class="stat-label">전체 Thing</div>
          </div>
          
          <div class="stat-card stat-pending">
            <div class="stat-number">${this.reviewData.pending}</div>
            <div class="stat-label">대기 중</div>
          </div>
          
          <div class="stat-card stat-someday">
            <div class="stat-number">${this.reviewData.someday}</div>
            <div class="stat-label">연기됨</div>
          </div>
          
          <div class="stat-card stat-done">
            <div class="stat-number">${this.reviewData.done}</div>
            <div class="stat-label">완료됨</div>
          </div>
        </div>
      </div>

      <div class="review-section">
        <h3>📅 최근 활동</h3>
        <div class="recent-activities">
          ${this.renderRecentActivities()}
        </div>
      </div>

      <div class="review-section">
        <h3>💡 개선 제안</h3>
        <div class="improvement-suggestions">
          ${this.renderSuggestions()}
        </div>
      </div>
    `;
  }

  /**
   * 최근 활동 렌더링
   * @returns {string} 최근 활동 HTML
   */
  renderRecentActivities() {
    if (!this.reviewData.recentActivities.length) {
      return '<p class="text-muted">최근 활동이 없습니다.</p>';
    }

    return this.reviewData.recentActivities
      .map(activity => `
        <div class="activity-item">
          <span class="activity-text">${this.escapeHtml(activity.activity)}</span>
          <span class="activity-time">${this.getRelativeTime(activity.timestamp)}</span>
        </div>
      `)
      .join('');
  }

  /**
   * 개선 제안 렌더링
   * @returns {string} 개선 제안 HTML
   */
  renderSuggestions() {
    const suggestions = [];

    if (this.reviewData.pending > 5) {
      suggestions.push('대기 중인 항목이 많습니다. 명확화를 진행해보세요.');
    }

    if (this.reviewData.someday > 10) {
      suggestions.push('연기된 항목이 많습니다. 정말 필요한 것들인지 검토해보세요.');
    }

    if (this.reviewData.done === 0) {
      suggestions.push('아직 완료된 항목이 없습니다. 작은 것부터 실행해보세요.');
    }

    if (suggestions.length === 0) {
      suggestions.push('좋습니다! GTD 시스템을 잘 활용하고 있습니다.');
    }

    return suggestions
      .map(suggestion => `<div class="suggestion-item">💡 ${suggestion}</div>`)
      .join('');
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