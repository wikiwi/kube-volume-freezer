package validation

import (
	"github.com/wikiwi/kube-volume-freezer/pkg/util/diff"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/reflect"
	"testing"
)

func TestValidateUID(t *testing.T) {
	testScenarios := []struct {
		uid    string
		issues []string
	}{
		{uid: "11111111-1111-1111-1111-111111111111", issues: nil},
		{uid: "invalid", issues: []string{
			"UID must be 36 characters long",
		}},
	}
	for _, x := range testScenarios {
		issues := ValidateUID(x.uid)
		if x.issues == nil && issues != nil {
			t.Errorf("uid: %q, expected to have no issues, but got %v", x.uid, issues)
		}
		if !reflect.DeepDerivative(x.issues, issues) {
			t.Errorf("uid: %q, %v", x.uid, diff.ObjectDiffDerivative(x.issues, issues))
		}
	}
}
