/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

// Package validation provides validation capabilites.
package validation

import (
	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/api/issues"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/validation"
)

// ValidateFreezeThawRequest validates a FreezeThawRequest and return a list of issues.
func ValidateFreezeThawRequest(ftr *api.FreezeThawRequest) api.IssueList {
	if ftr.Action != "freeze" && ftr.Action != "thaw" {
		return api.IssueList{
			issues.InvalidField("action", "action must be either freeze or thaw but was %q", ftr.Action)}
	}
	return nil
}

// ValidateUIDParameter validates given uid and return a list of issues.
// Argument param is the name of the Parameter being included in the list of issues.
func ValidateUIDParameter(param string, uid string) (list api.IssueList) {
	valIssues := validation.ValidateUID(uid)
	for _, p := range valIssues {
		list = append(list, issues.InvalidParameter(param, p))
	}
	return list
}

// ValidateQualifiedNameParameter validates given qualified name and return a list of issues.
// Argument param is the name of the Parameter being included in the list of issues.
func ValidateQualifiedNameParameter(param string, name string) (list api.IssueList) {
	valIssues := validation.ValidateQualitfiedName(name)
	for _, p := range valIssues {
		list = append(list, issues.InvalidParameter(param, p))
	}
	return list
}
