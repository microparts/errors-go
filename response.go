package errors

type FieldName string
type ValidationError string
type ErrorCode string
type DebugData interface{}

type Response struct {
	Error ErrorObject `json:"error,omitempty"`
}

type ErrorObject struct {
	Message    interface{}                     `json:"message"`
	Code       *ErrorCode                      `json:"code,omitempty"`
	Validation map[FieldName][]ValidationError `json:"validation,omitempty"`
	Debug      *DebugData                      `json:"debug,omitempty"`
}
