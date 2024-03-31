package models

import "errors"

var (
	ErrInvalidInput        = errors.New("INVALID_INPUT")
	ErrInternalServer      = errors.New("INTERNAL_SERVER_ERROR")
	ErrCompanyDoesntExists = errors.New("COMPANY_DOES_NOT_EXIST")
	ErrUsernameExists      = errors.New("USERNAME_EXISTS")
	ErrUserNotFound        = errors.New("USER_NOT_FOUND")
	ErrPermissionDenied    = errors.New("PERMISSION_DENIED")
	ErrCompanyNotFound     = errors.New("COMPANY_NOT_FOUND")
	ErrInterviewNotFound   = errors.New("INTERVIEW_NOT_FOUND")
	ErrQuestionNotFound    = errors.New("QUESTION_NOT_FOUND")
)
