/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

// Package errors contains all API errors.
package errors

import (
	"github.com/wikiwi/kube-volume-freezer/pkg/api"
)

// Unexpected returns an error with a 500 error code.
func Unexpected(message string) *api.Error {
	return &api.Error{Code: 500, Message: message}
}

// IsUnexpected returns true when err is an Unexpected error.
func IsUnexpected(err error) bool {
	if apiErr, ok := err.(*api.Error); ok && apiErr.Code == 500 {
		return true
	}
	return false
}

// NotFound returns an error with a 404 error code.
func NotFound(message string) *api.Error {
	return &api.Error{Code: 404, Message: message}
}

// IsNotFound returns true when err is a NotFound error.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*api.Error); ok && apiErr.Code == 404 {
		return true
	}
	return false
}

// BadRequest returns an error with a 400 error code.
// This is used when a request could not be parsed.
func BadRequest(message string) *api.Error {
	return &api.Error{Code: 400, Message: message}
}

// IsBadRequest returns true when err is a BadRequest error.
func IsBadRequest(err error) bool {
	if apiErr, ok := err.(*api.Error); ok && apiErr.Code == 400 {
		return true
	}
	return false
}

// UnprocessableEntity returns an error with a 422 error code.
// This is used when a request contains invalid data, but could be parsed
// syntaxtically.
func UnprocessableEntity(message string) *api.Error {
	return &api.Error{Code: 422, Message: message}
}

// IsUnprocessableEntity returns true when err is a UnprocessableEntity error.
func IsUnprocessableEntity(err error) bool {
	if apiErr, ok := err.(*api.Error); ok && apiErr.Code == 422 {
		return true
	}
	return false
}

// Forbidden returns an error with a 403 error code.
func Forbidden(message string) *api.Error {
	return &api.Error{Code: 403, Message: message}
}

// IsForbidden returns true when err is a Forbidden error.
func IsForbidden(err error) bool {
	if apiErr, ok := err.(*api.Error); ok && apiErr.Code == 403 {
		return true
	}
	return false
}
