package internal

import (
	"fmt"
	"os"
)

const (
	ExitSuccess    = 0
	ExitError      = 1
	ExitAuth       = 2
	ExitNotFound   = 3
	ExitValidation = 4
)

type ZohoCliError struct {
	Message  string
	ExitCode int
}

func (e *ZohoCliError) Error() string {
	return e.Message
}

func NewError(msg string) *ZohoCliError {
	return &ZohoCliError{Message: msg, ExitCode: ExitError}
}

func NewAuthError(msg string) *ZohoCliError {
	return &ZohoCliError{Message: msg, ExitCode: ExitAuth}
}

func NewNotFoundError(msg string) *ZohoCliError {
	return &ZohoCliError{Message: msg, ExitCode: ExitNotFound}
}

func NewValidationError(msg string) *ZohoCliError {
	return &ZohoCliError{Message: msg, ExitCode: ExitValidation}
}

type ZohoAPIError struct {
	ZohoCliError
	StatusCode int
}

func NewAPIError(statusCode int, body string) *ZohoAPIError {
	exitCode := ExitError
	switch statusCode {
	case 401:
		exitCode = ExitAuth
	case 404:
		exitCode = ExitNotFound
	}
	return &ZohoAPIError{
		ZohoCliError: ZohoCliError{
			Message:  fmt.Sprintf("Zoho API error %d: %s", statusCode, body),
			ExitCode: exitCode,
		},
		StatusCode: statusCode,
	}
}

func RequireFlag(cmd interface{ String(string) string }, flag, envHint string) (string, error) {
	v := cmd.String(flag)
	if v == "" {
		return "", NewValidationError(fmt.Sprintf("--%s flag or %s env var required", flag, envHint))
	}
	return v, nil
}

func Err(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}
