package blame

import (
	"fmt"
	"runtime"
	"strings"
)

type Blame[T any] interface {
	error
	GetErrCode() string
	GetMessage() string
	GetFields() map[string]interface{}
	// GetSource() string
	GetComponent() string
	GetResponseType() string
	GetCauses() []error
	WithField(key string, value interface{}) *Error[T]
	WithCause(err error) *Error[T]
	SetComponent(component string) *Error[T]
	SetResponseType(responseType string) *Error[T]
}

type Error[T any] struct {
	statusCode   string                 `json:"statusCode"`
	errCode      string                 `json:"errCode"`
	component    string                 `json:"component"`
	responseType string                 `json:"responseType"`
	message      string                 `json:"message"`
	description  string                 `json:"description"`
	fields       map[string]interface{} `json:"fields"`
	causes       []error                `json:"causes"`
	// source       string                 `json:"source"`
}

func NewBlame[T any](
	statusCode, errCode, message string,
) Blame[T] {
	return NewError[T](statusCode, errCode, message)
}

func NewError[T any](statusCode, errorCode, message string) *Error[T] {
	return &Error[T]{
		statusCode:  statusCode,
		errCode:     errorCode,
		message:     message,
		description: message,
		fields:      map[string]interface{}{},
		//	source:       getSource(),
	}
}

func (e *Error[T]) GetErrCode() string {
	return e.errCode
}

func (e *Error[T]) GetMessage() string {
	return e.message
}

func (e *Error[T]) GetFields() map[string]interface{} {
	return e.fields
}

// func (e *Error[T]) GetSource() string {
// 	return e.source
// }

func (e *Error[T]) GetComponent() string {
	return e.component
}

func (e *Error[T]) GetResponseType() string {
	return e.responseType
}

func (e *Error[T]) GetCauses() []error {
	return e.causes
}

func (e *Error[T]) WithField(key string, value interface{}) *Error[T] {
	e.fields[key] = value
	return e
}

func (e *Error[T]) WithCause(err error) *Error[T] {
	e.causes = append(e.causes, err)
	return e
}

func (e *Error[T]) SetComponent(component string) *Error[T] {
	e.component = component
	return e
}

func (e *Error[T]) SetResponseType(responseType string) *Error[T] {
	e.responseType = responseType
	return e
}
func (e *Error[T]) Error() string {
	return e.errCode
}

func getSource() string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", strings.TrimPrefix(file, runtime.GOROOT()+"/src/"), line)
}

func FindErrorDefinition[T any](errors []*Error[T], errorCode string) *Error[T] {
	for _, err := range errors {
		if err.errCode == errorCode {
			return err
		}
	}
	return nil
}

func ReplaceDynamicValues[T any](err *Error[T], data map[string]interface{}) *Error[T] {
	for key, value := range data {
		err.message = strings.Replace(err.message, "{{."+key+"}}", fmt.Sprintf("%v", value), -1)
		err.description = strings.Replace(err.description, "{{."+key+"}}", fmt.Sprintf("%v", value), -1)
		if err.fields != nil {
			for fieldKey, fieldValue := range err.fields {
				if fieldValueStr, ok := fieldValue.(string); ok {
					err.fields[fieldKey] = strings.Replace(fieldValueStr, "{{."+key+"}}", fmt.Sprintf("%v", value), -1)
				}
			}
		}
	}
	return err
}
