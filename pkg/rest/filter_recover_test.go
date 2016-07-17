/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package rest_test

import (
	"testing"

	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/clienttest"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest/resttest"
)

func TestRecoverFilter(t *testing.T) {
	panicingFunc := func(request *restful.Request, response *restful.Response) {
		panic("panic")
	}

	s := resttest.RunFilterTestServer(rest.RecoverFilter, panicingFunc)
	defer s.Close()

	client := generic.NewOrDie(s.URL, nil)
	req := client.NewRequestOrDie("GET", "/", nil)
	exp := clienttest.ResponseExpectation{
		Code:   500,
		Entity: &api.Error{Code: 500},
	}
	exp.DoAndValidateOrDie(t, client, req)
}
