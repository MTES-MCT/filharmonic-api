// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Sso is an autogenerated mock type for the Sso type
type Sso struct {
	mock.Mock
}

// ValidateTicket provides a mock function with given fields: ticket
func (_m *Sso) ValidateTicket(ticket string) (string, error) {
	ret := _m.Called(ticket)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(ticket)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ticket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
