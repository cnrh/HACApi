package repository

import "errors"

// The error thrown when the body parameters are invalid.
var ErrorBadBodyParams = errors.New("bad body params")

// The error thrown when there is an invalid authentication field.
var ErrorInvalidAuthentication = errors.New("invalid username/password/base")

// The error thrown when there is an internal error.
var ErrorInternalError = errors.New("resource not found. possibly an internal error")
