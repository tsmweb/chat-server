package common

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

// MockJWT injects mock dependency into Controller layer.
type MockJWT struct {
	mock.Mock
}

// GenerateToken represents the simulated method for to generate token feature in the JWT.
func (m *MockJWT) GenerateToken(payload map[string]interface{}, exp int) (string, error) {
	args := m.Called(payload, exp)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

// ExtractToken represents the simulated method for the extract token feature in the JWT.
func (m *MockJWT) ExtractToken(r *http.Request) (string, error) {
	args := m.Called(r)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

// GetDataToken represents the simulated method for the resource to get token data in JWT.
func (m *MockJWT) GetDataToken(r *http.Request, key string) (interface{}, error) {
	args := m.Called(r, key)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}
