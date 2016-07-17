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

	"github.com/wikiwi/kube-volume-freezer/pkg/client/clienttest"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest/resttest"
)

func TestTokenAuthFilter(t *testing.T) {
	okFunc := func(request *restful.Request, response *restful.Response) {
		response.WriteEntity("ok")
	}

	filter := rest.NewTokenAuthFilter("mytoken")
	s := resttest.RunFilterTestServer(filter, okFunc)
	defer s.Close()

	client := generic.NewOrDie(s.URL, nil)
	req := client.NewRequestOrDie("GET", "/", nil)
	exp := clienttest.ResponseExpectation{
		Code: 403,
	}
	exp.DoAndValidateOrDie(t, client, req)

	req.Header.Set("Authorization", "Bearer mytoken")
	exp = clienttest.ResponseExpectation{
		Code: 200,
	}
	exp.DoAndValidateOrDie(t, client, req)
}
