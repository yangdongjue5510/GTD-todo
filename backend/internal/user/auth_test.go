package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============ NewTokenService Tests ============

func TestNewTokenService_Success(t *testing.T) {
	// given
	secretKey := "test-secret-key-32-bytes-long!!"
	nowFunc := time.Now

	// when
	service, err := NewTokenService(secretKey, nowFunc)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, []byte(secretKey), service.secretKey)
	assert.NotNil(t, service.now)
}

func TestNewTokenService_EmptySecretKey(t *testing.T) {
	// given
	secretKey := ""
	nowFunc := time.Now

	// when
	service, err := NewTokenService(secretKey, nowFunc)

	// then
	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "secret key should not be empty")
}

func TestNewTokenService_NilNowFunction(t *testing.T) {
	// given
	secretKey := "test-secret-key-32-bytes-long!!"

	// when
	service, err := NewTokenService(secretKey, nil)

	// then
	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "now function")
}

// ============ Issue Method Tests ============

func TestIssue_Success(t *testing.T) {
	// given
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	nowFunc := func() time.Time { return fixedTime }

	service, _ := NewTokenService("test-secret-key", nowFunc)
	userID := 123
	email := "test@example.com"
	duration := 1 * time.Hour

	// when
	token, err := service.Issue(userID, email, duration)

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestIssue_VerifyTokenStructure(t *testing.T) {
	// given
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	nowFunc := func() time.Time { return fixedTime }

	service, _ := NewTokenService("test-secret-key", nowFunc)
	userID := 456
	email := "user@example.com"
	duration := 2 * time.Hour

	// when
	token, err := service.Issue(userID, email, duration)
	assert.NoError(t, err)

	// Parse the token back to verify structure
	claims, err := service.Parse(token)

	// then
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, fixedTime.Unix(), claims.IssuedAt.Unix())
	assert.Equal(t, fixedTime.Add(duration).Unix(), claims.ExpiresAt.Unix())
	assert.Equal(t, "gtd-todo-app", claims.Issuer)
}

func TestIssue_WithMockedTime(t *testing.T) {
	// given
	fixedTime := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	nowFunc := func() time.Time { return fixedTime }

	service, _ := NewTokenService("my-secret-key", nowFunc)
	duration := 30 * time.Minute

	// when
	token, err := service.Issue(1, "test@test.com", duration)
	assert.NoError(t, err)

	claims, err := service.Parse(token)

	// then
	assert.NoError(t, err)
	assert.Equal(t, fixedTime.Unix(), claims.IssuedAt.Unix())
	assert.Equal(t, fixedTime.Add(30*time.Minute).Unix(), claims.ExpiresAt.Unix())
}

// ============ Parse Method Tests ============

func TestParse_ValidToken(t *testing.T) {
	// given
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	nowFunc := func() time.Time { return fixedTime }

	service, _ := NewTokenService("test-secret-key", nowFunc)
	userID := 100
	email := "valid@example.com"
	duration := 1 * time.Hour

	token, _ := service.Issue(userID, email, duration)

	// when
	claims, err := service.Parse(token)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestParse_VerifyClaimsContent(t *testing.T) {
	// given
	nowFunc := time.Now
	service, _ := NewTokenService("test-secret-key", nowFunc)

	expectedUserID := 789
	expectedEmail := "verify@example.com"

	token, _ := service.Issue(expectedUserID, expectedEmail, 1*time.Hour)

	// when
	claims, err := service.Parse(token)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, claims.UserID)
	assert.Equal(t, expectedEmail, claims.Email)
}

func TestParse_InvalidSignature(t *testing.T) {
	// given
	nowFunc := time.Now
	service1, _ := NewTokenService("secret-key-1", nowFunc)
	service2, _ := NewTokenService("secret-key-2", nowFunc)

	// service1으로 토큰 생성
	token, _ := service1.Issue(1, "test@example.com", 1*time.Hour)

	// when - service2로 파싱 시도
	claims, err := service2.Parse(token)

	// then
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "failed to parse token")
}

func TestParse_ExpiredToken(t *testing.T) {
	// given
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currentTime := fixedTime
	nowFunc := func() time.Time { return currentTime }

	service, _ := NewTokenService("test-secret-key", nowFunc)

	// 1시간 duration으로 토큰 생성
	token, _ := service.Issue(1, "test@example.com", 1*time.Hour)

	// 시간을 2시간 후로 이동 (토큰 만료)
	currentTime = fixedTime.Add(2 * time.Hour)

	// when - 만료된 토큰 파싱
	claims, err := service.Parse(token)

	// then
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "failed to parse token")
}

func TestParse_MalformedToken(t *testing.T) {
	// given
	nowFunc := time.Now
	service, _ := NewTokenService("test-secret-key", nowFunc)

	malformedToken := "invalid.token.string"

	// when
	claims, err := service.Parse(malformedToken)

	// then
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "failed to parse token")
}

func TestParse_EmptyToken(t *testing.T) {
	// given
	nowFunc := time.Now
	service, _ := NewTokenService("test-secret-key", nowFunc)

	// when
	claims, err := service.Parse("")

	// then
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "failed to parse token")
}

func TestParse_TokenWithinValidTimeRange(t *testing.T) {
	// given
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currentTime := fixedTime
	nowFunc := func() time.Time { return currentTime }

	service, _ := NewTokenService("test-secret-key", nowFunc)

	// 1시간 duration으로 토큰 생성
	token, _ := service.Issue(1, "test@example.com", 1*time.Hour)

	// 30분 후로 이동 (여전히 유효)
	currentTime = fixedTime.Add(30 * time.Minute)

	// when
	claims, err := service.Parse(token)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, 1, claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestParse_TokenExactlyAtExpiration(t *testing.T) {
	// given
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currentTime := fixedTime
	nowFunc := func() time.Time { return currentTime }

	service, _ := NewTokenService("test-secret-key", nowFunc)

	duration := 1 * time.Hour
	token, _ := service.Issue(1, "test@example.com", duration)

	// 정확히 만료 시간으로 이동
	currentTime = fixedTime.Add(duration)

	// when
	claims, err := service.Parse(token)

	// then
	// ExpiresAt와 같은 시간은 만료로 간주
	assert.Error(t, err)
	assert.Nil(t, claims)
}
