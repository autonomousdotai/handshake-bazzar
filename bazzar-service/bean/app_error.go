package bean

import "fmt"

type AppError struct {
	OrgError   error
	Message    string
	Code       int
	KeyMessage string
}

func (e AppError) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

func NewError(statusKey string, err error) error {
	return AppError{
		OrgError:   err,
		Code:       CodeMessage[statusKey].Code,
		Message:    CodeMessage[statusKey].Message,
		KeyMessage: statusKey,
	}
}
