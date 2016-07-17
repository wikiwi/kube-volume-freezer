/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package validation

import (
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/diff"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/reflect"
)

func TestValidateFreezeThawRequest(t *testing.T) {
	testScenarios := []struct {
		action string
		issues api.IssueList
	}{
		{action: "freeze", issues: nil},
		{action: "thaw", issues: nil},
		{action: "invalid", issues: api.IssueList{
			&api.Issue{
				Reason:       "InvalidField",
				Location:     "action",
				LocationType: "field",
			}},
		},
	}
	for _, x := range testScenarios {
		req := &api.FreezeThawRequest{Action: x.action}
		issues := ValidateFreezeThawRequest(req)
		if x.issues == nil && issues != nil {
			t.Errorf("action: %q, expected to have no issues, but got %v", x.action, issues)
		}
		if !reflect.DeepDerivative(x.issues, issues) {
			t.Errorf("action: %q, %v", x.action, diff.ObjectDiffDerivative(x.issues, issues))
		}
	}
}

func TestValidateUIDParameter(t *testing.T) {
	testScenarios := []struct {
		param  string
		uid    string
		issues api.IssueList
	}{
		{param: "paramName", uid: "11111111-1111-1111-1111-111111111111", issues: nil},
		{param: "paramName", uid: "", issues: api.IssueList{
			&api.Issue{
				Reason:       "InvalidParameter",
				Location:     "paramName",
				LocationType: "parameter",
			}},
		},
	}
	for _, x := range testScenarios {
		issues := ValidateUIDParameter(x.param, x.uid)
		if x.issues == nil && issues != nil {
			t.Errorf("param: %q, uid: %q, expected to have no issues, but got %v", x.param, x.uid, issues)
		}
		if !reflect.DeepDerivative(x.issues, issues) {
			t.Errorf("param: %q, uid: %q, %v", x.param, x.uid, diff.ObjectDiffDerivative(x.issues, issues))
		}
	}
}

func TestValidateNameParameter(t *testing.T) {
	testScenarios := []struct {
		param  string
		name   string
		issues api.IssueList
	}{
		{param: "paramName", name: "valid-name", issues: nil},
		{param: "paramName", name: "$Invalid0", issues: api.IssueList{
			&api.Issue{
				Reason:       "InvalidParameter",
				Location:     "paramName",
				LocationType: "parameter",
			}},
		},
	}
	for _, x := range testScenarios {
		issues := ValidateQualifiedNameParameter(x.param, x.name)
		if x.issues == nil && issues != nil {
			t.Errorf("param: %q, name: %q, expected to have no issues, but got %v", x.param, x.name, issues)
		}
		if !reflect.DeepDerivative(x.issues, issues) {
			t.Errorf("param: %q, name: %q, %v", x.param, x.name, diff.ObjectDiffDerivative(x.issues, issues))
		}
	}
}
