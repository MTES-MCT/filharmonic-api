// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import emails "github.com/MTES-MCT/filharmonic-api/emails"
import mock "github.com/stretchr/testify/mock"

// EmailService is an autogenerated mock type for the EmailService type
type EmailService struct {
	mock.Mock
}

// Send provides a mock function with given fields: email
func (_m *EmailService) Send(email emails.Email) error {
	ret := _m.Called(email)

	var r0 error
	if rf, ok := ret.Get(0).(func(emails.Email) error); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
