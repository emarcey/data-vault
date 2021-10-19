package common

import (
	"fmt"

	"github.com/lib/pq"
)

type ErrorWithCode interface {
	Error() string
	Code() int
}

type InitializationError struct {
	dependency  string
	message     string
	messageArgs []interface{}
}

func (e InitializationError) Error() string {
	return fmt.Sprintf("Error during server startup module, %s, with message: ", e.dependency) +
		fmt.Sprintf(e.message, e.messageArgs...)
}

func NewInitializationError(dependency, message string, messageArgs ...interface{}) InitializationError {
	return InitializationError{dependency: dependency, message: message, messageArgs: messageArgs}
}

type SecretsError struct {
	secretsManagerType string
	method             string
	message            string
	messageArgs        []interface{}
}

func (e SecretsError) Error() string {
	return fmt.Sprintf("Error in %s %s: ", e.secretsManagerType, e.method) + fmt.Sprintf(e.message, e.messageArgs...)
}

func (e SecretsError) Code() int {
	return 500
}

func NewMongoError(method, message string) SecretsError {
	return SecretsError{secretsManagerType: "mongodb", method: method, message: message}
}

func NewMongoGetOrPutSecretError(message string, messageArgs ...interface{}) SecretsError {
	return SecretsError{secretsManagerType: "mongodb", method: "GetOrPutSecret", message: message, messageArgs: messageArgs}
}

type DatabaseError struct {
	operation   string
	message     string
	code        int
	messageArgs []interface{}
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("Error during database operation, %s, with message: ", e.operation) +
		fmt.Sprintf(e.message, e.messageArgs...)
}

func (e DatabaseError) Code() int {
	return e.code
}

func NewDatabaseError(originalErr error, operation, message string, messageArgs ...interface{}) ErrorWithCode {
	pqErr, ok := originalErr.(*pq.Error)
	if ok && pqErr != nil {
		// not exhaustive, adding as I find potential errors I want to clean up
		switch pqErr.Code {
		case "23505":
			return NewResourceAlreadyExistsError(operation, pqErr.Detail)
		case "23503":
			return NewInvalidParamsError(operation, pqErr.Detail)
		case "22P02":
			return NewInvalidParamsError(operation, pqErr.Message)
		default:
		}

	}
	if message == "" {
		message = originalErr.Error()
	}
	return DatabaseError{operation: operation, message: message, messageArgs: messageArgs, code: 500}

}

type AuthorizationError struct{}

func (e AuthorizationError) Error() string {
	return fmt.Sprintf("Authorization Not Valid.")
}

func (e AuthorizationError) Code() int {
	return 401
}

func NewAuthorizationError() AuthorizationError {
	return AuthorizationError{}
}

type InvalidParamsError struct {
	functionName string
	message      string
	messageArgs  []interface{}
}

func (e InvalidParamsError) Error() string {
	return fmt.Sprintf("Error invalid params at %s: ", e.functionName) + fmt.Sprintf(e.message, e.messageArgs...)
}

func (e InvalidParamsError) Code() int {
	return 400
}

func NewInvalidParamsError(functionName string, message string, messageArgs ...interface{}) InvalidParamsError {
	return InvalidParamsError{functionName: functionName, message: message, messageArgs: messageArgs}
}

type ResourceAlreadyExistsError struct {
	operation string
	message   string
}

func (e ResourceAlreadyExistsError) Error() string {
	return fmt.Sprintf("Error for operation, %s: %s", e.operation, e.message)
}

func (e ResourceAlreadyExistsError) Code() int {
	return 409
}

func NewResourceAlreadyExistsError(operation, message string) ResourceAlreadyExistsError {
	return ResourceAlreadyExistsError{operation: operation, message: message}
}

type ResourceNotFoundError struct {
	operation string
	field     string
	value     string
}

func (e ResourceNotFoundError) Error() string {
	return fmt.Sprintf("Resource not found in operation, %s, for value, %s, at field, %s", e.operation, e.value, e.field)
}

func (e ResourceNotFoundError) Code() int {
	return 404
}

func NewResourceNotFoundError(operation, field, value string) ResourceNotFoundError {
	return ResourceNotFoundError{operation: operation, field: field, value: value}
}
