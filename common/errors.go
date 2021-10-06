package common

import (
	"fmt"
)


type InitializationError struct {
	dependency string
	message string
}

func (i InitializationError) Error() string {
	return fmt.Sprintf("Error during server startup module, %s, with message: %s", i.dependency, i.message)
}

func NewInitializationError(dependency, message string) InitializationError {
	return InitializationError{dependency: dependency, message: message}
}