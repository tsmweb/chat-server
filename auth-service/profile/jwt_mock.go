package profile

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

// mockJWT injects mock dependency into Controller layer.
type mockJWT struct {
	mock.Mock
}

func (m *mockJWT) GenerateToken(id string, exp int) (string, error) {
	args := m.Called(id, exp)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

func (m *mockJWT) ExtractToken(r *http.Request) (string, error) {
	args := m.Called(r)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

func (m *mockJWT) GetDataToken(r *http.Request, key string) (interface{}, error) {
	args := m.Called(r, key)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}
