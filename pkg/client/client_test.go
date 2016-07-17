/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package client

import (
	"fmt"
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/version"
)

func ExampleClient() {
	c := NewOrDie("http://kube-volume-freezer:8080", &Options{Token: "token"})

	list, err := c.Volumes().List("default", "podName")
	if err != nil {
		panic(err)
	}

	for _, volume := range list.Items {
		fmt.Println(volume)
	}
}

func TestUserAgent(t *testing.T) {
	c := NewOrDie("http://localhost", nil)
	r, err := c.NewRequest("GET", "relative", nil)
	if err != nil {
		t.Fatal(err)
	}
	expUA := "kvf/" + version.Version
	if r.UserAgent() != expUA {
		t.Fatalf("expected user-agent header to be %q but was %q", expUA, r.UserAgent())
	}
}

func TestAuthorizationHeader(t *testing.T) {
	c := NewOrDie("http://localhost", &Options{Token: "token"})
	r, err := c.NewRequest("GET", "relative", nil)
	if err != nil {
		t.Fatal(err)
	}
	expAuth := "Bearer token"
	if auth := r.Header.Get("Authorization"); auth != expAuth {
		t.Fatalf("expected authorization header to be %q but was %q", expAuth, auth)
	}
}
