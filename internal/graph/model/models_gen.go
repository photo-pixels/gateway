// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Mutation struct {
}

type Query struct {
}

type Result struct {
	Success bool `json:"success"`
}

type AuthStatus string

const (
	AuthStatusAuthStatusUnknown      AuthStatus = "AUTH_STATUS_UNKNOWN"
	AuthStatusAuthStatusSentInvite   AuthStatus = "AUTH_STATUS_SENT_INVITE"
	AuthStatusAuthStatusNotActivated AuthStatus = "AUTH_STATUS_NOT_ACTIVATED"
	AuthStatusAuthStatusActivated    AuthStatus = "AUTH_STATUS_ACTIVATED"
	AuthStatusAuthStatusBlocked      AuthStatus = "AUTH_STATUS_BLOCKED"
)

var AllAuthStatus = []AuthStatus{
	AuthStatusAuthStatusUnknown,
	AuthStatusAuthStatusSentInvite,
	AuthStatusAuthStatusNotActivated,
	AuthStatusAuthStatusActivated,
	AuthStatusAuthStatusBlocked,
}

func (e AuthStatus) IsValid() bool {
	switch e {
	case AuthStatusAuthStatusUnknown, AuthStatusAuthStatusSentInvite, AuthStatusAuthStatusNotActivated, AuthStatusAuthStatusActivated, AuthStatusAuthStatusBlocked:
		return true
	}
	return false
}

func (e AuthStatus) String() string {
	return string(e)
}

func (e *AuthStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AuthStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AuthStatus", str)
	}
	return nil
}

func (e AuthStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Бизнес ошибки
type Error string

const (
	//  Сессия не найдена
	ErrorSessionNoFound Error = "SESSION_NO_FOUND"
	//  Не авторизован
	ErrorNoAuth Error = "NO_AUTH"
)

var AllError = []Error{
	ErrorSessionNoFound,
	ErrorNoAuth,
}

func (e Error) IsValid() bool {
	switch e {
	case ErrorSessionNoFound, ErrorNoAuth:
		return true
	}
	return false
}

func (e Error) String() string {
	return string(e)
}

func (e *Error) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Error(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Error", str)
	}
	return nil
}

func (e Error) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
