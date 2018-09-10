package utils

import (
	"github.com/ansel1/merry"
	"net/http"
)

// Input Errors
var (
	// InputError wraps data input errors
	InputError = merry.New("Input error")

	// ArgumentError error with one or more input arguments
	ArgumentError = merry.WithMessage(InputError, "Invalid argument error").WithHTTPCode(http.StatusBadRequest)

	// RequestBodyError error with the request body
	RequestBodyError = merry.WithMessage(InputError, "Invalid request body errors").WithHTTPCode(http.StatusBadRequest)

	// StatusConflictError error with conflicting statuses
	StatusConflictError = merry.WithMessage(InputError, "Conflicting status error").WithHTTPCode(http.StatusConflict)
)

// HTTP Errors
var (

	// HTTPError is the base error for any errors related to http client calls
	HTTPError = merry.New("HTTP error")

	// HTTPBadRequestError is for 400 error codes
	HTTPBadRequestError = merry.WithMessage(HTTPError, "400 Bad Request").WithHTTPCode(http.StatusBadRequest)

	// HTTPUnauthorizedError is for 401 error codes
	HTTPUnauthorizedError = merry.WithMessage(HTTPBadRequestError, "401 Unauthorized").WithHTTPCode(http.StatusUnauthorized)

	// HTTPForbiddenError is for 403 error codes
	HTTPForbiddenError = merry.WithMessage(HTTPBadRequestError, "403 Forbidden").WithHTTPCode(http.StatusForbidden)

	// HTTPNotFoundError for 404 error codes
	HTTPNotFoundError = HTTPError.WithMessage("404 Not Found").WithHTTPCode(http.StatusNotFound)
)

// Marshalling Errors
var (
	// MarshalingError wraps errors that occur during marshaling
	MarshalingError = merry.New("Marshaling error")

	// JSONMarshalingError indicates an issue with masrhaling a struct to json
	JSONMarshalingError = merry.WithMessage(MarshalingError, "JSON marshaling error")
)
