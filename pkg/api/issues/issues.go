/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

// Package issues contains all issues that can be included in API errors.
package issues

import (
	"fmt"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
)

// MinionNotFound describes a missing Minion.
func MinionNotFound(message string, a ...interface{}) *api.Issue {
	return &api.Issue{Reason: "MinionNotFound", Message: fmt.Sprintf(message, a...)}
}

// PodNotFound describes a missing Pod.
func PodNotFound(message string, a ...interface{}) *api.Issue {
	return &api.Issue{Reason: "PodNotFound", Message: fmt.Sprintf(message, a...)}
}

// VolumeNotFound describes a missing Pod Volume.
func VolumeNotFound(message string, a ...interface{}) *api.Issue {
	return &api.Issue{Reason: "VolumeNotFound", Message: fmt.Sprintf(message, a...)}
}

// InvalidParameter describes a validation error in a request path parameter.
func InvalidParameter(param string, message string, a ...interface{}) *api.Issue {
	return &api.Issue{Reason: "InvalidParameter", Message: fmt.Sprintf(message, a...), LocationType: "parameter", Location: param}
}

// InvalidField describes a validation error in the request body.
func InvalidField(field string, message string, a ...interface{}) *api.Issue {
	return &api.Issue{Reason: "InvalidField", Message: fmt.Sprintf(message, a...), LocationType: "field", Location: field}
}
