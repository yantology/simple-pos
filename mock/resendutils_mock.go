package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/yantology/retail-pro-be/pkg/customerror"
)

type MockResendUtils struct {
	mock.Mock
}

func (m *MockResendUtils) Send(html, subject, string, to []string) *customerror.CustomError {
	args := m.Called(html, subject, to)
	return args.Get(0).(*customerror.CustomError)
}
