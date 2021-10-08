package common

import (
	"fmt"
)

type InitializationError struct {
	dependency string
	message    string
}

func (e InitializationError) Error() string {
	return fmt.Sprintf("Error during server startup module, %s, with message: "+e.message, e.dependency)
}

func NewInitializationError(dependency, message string) InitializationError {
	return InitializationError{dependency: dependency, message: message}
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

func NewMongoError(method, message string) SecretsError {
	return SecretsError{secretsManagerType: "mongodb", method: method, message: message}
}

func NewMongoGetOrPutSecretError(message string, messageArgs ...interface{}) SecretsError {
	return SecretsError{secretsManagerType: "mongodb", method: "GetOrPutSecret", message: message, messageArgs: messageArgs}
}
