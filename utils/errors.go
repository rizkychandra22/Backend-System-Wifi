package utils

type AppError struct {
	StatusCode int
	Message    string
	Data       interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{StatusCode: code, Message: msg}
}

func NewAppErrorWithData(code int, msg string, data interface{}) *AppError {
	return &AppError{StatusCode: code, Message: msg, Data: data}
}
