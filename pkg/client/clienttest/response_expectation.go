/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package clienttest

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/diff"
	reflectutil "github.com/wikiwi/kube-volume-freezer/pkg/util/reflect"
)

// ResponseExpectation executes a HTTP request and validates its response.
type ResponseExpectation struct {
	Code    int
	Entity  interface{}
	Headers http.Header
}

// DoAndValidate executes a HTTP request and validates its response.
func (exp *ResponseExpectation) DoAndValidate(client *generic.Client, req *http.Request) error {
	var store interface{}
	if exp.Entity != nil {
		store = reflect.New(reflect.Indirect(reflect.ValueOf(exp.Entity)).Type()).Interface()
	}
	resp, err := client.Do(req, store)
	if c := exp.Code; c >= 200 && c <= 299 {
		if err != nil {
			return fmt.Errorf("error getting response %v", err)
		}

		if resp.StatusCode != exp.Code {
			return fmt.Errorf("expected status code to be %d, was %d",
				exp.Code, resp.StatusCode)
		}

		if store != nil && !reflectutil.DeepDerivative(exp.Entity, store) {
			return fmt.Errorf("%v", diff.ObjectDiffDerivative(exp.Entity, store))
		}
	} else {
		if err == nil {
			return fmt.Errorf("expected error but was nil")
		}

		if resp.StatusCode != exp.Code {
			return fmt.Errorf("expected status code to be %d, was %d",
				exp.Code, resp.StatusCode)
		}

		if apiErr, ok := err.(*api.Error); ok {
			if exp.Code != apiErr.Code {
				return fmt.Errorf("expected error code to be %d, was %d",
					exp.Code, apiErr.Code)
			}
		} else {
			return fmt.Errorf("expected *api.Error but was %T", err)
		}
	}

	for key, items := range exp.Headers {
		if !reflect.DeepEqual(items, resp.Header[key]) {
			return fmt.Errorf("expected Header %q to be %q, was %q",
				key, items, resp.Header[key])
		}
	}

	return nil
}

// DoAndValidateOrDie is for convenience failing the test on error.
func (exp *ResponseExpectation) DoAndValidateOrDie(t *testing.T, client *generic.Client, req *http.Request) {
	err := exp.DoAndValidate(client, req)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
