package validation

import (
	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/api/issues"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/validation"
)

// ValidateFreezeThawRequest validates a FreezeThawRequest.
func ValidateFreezeThawRequest(ftr *api.FreezeThawRequest) api.IssueList {
	if ftr.Action != "freeze" && ftr.Action != "thaw" {
		return api.IssueList{
			issues.InvalidField("action", "action must be either freeze or thaw but was %q", ftr.Action)}
	}
	return nil
}

func ValidateUIDParameter(param string, uid string) (list api.IssueList) {
	valIssues := validation.ValidateUID(uid)
	for _, p := range valIssues {
		list = append(list, issues.InvalidParameter(param, p))
	}
	return list
}

func ValidateQualifiedNameParameter(param string, name string) (list api.IssueList) {
	valIssues := validation.ValidateQualitfiedName(name)
	for _, p := range valIssues {
		list = append(list, issues.InvalidParameter(param, p))
	}
	return list
}
