/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package rest

import (
	"encoding/json"
	"net/http/httptest"

	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/log"
)

// LogFilter enriches the log with requenst/response data.
func LogFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Instance().Debugf("Processing %s %s", req.Request.Method, req.Request.RequestURI)
	origWriter := resp.ResponseWriter
	recorder := httptest.NewRecorder()
	resp.ResponseWriter = recorder
	chain.ProcessFilter(req, resp)

	var errDescription string
	if recorder.Code > 399 {
		var apiError api.Error
		err := json.Unmarshal(recorder.Body.Bytes(), &apiError)
		if err == nil {
			errDescription = apiError.Message
		}
	}

	for header, items := range recorder.HeaderMap {
		for _, item := range items {
			origWriter.Header().Add(header, item)
		}
	}

	origWriter.WriteHeader(recorder.Code)
	_, err := recorder.Body.WriteTo(origWriter)
	if err != nil {
		panic(err)
	}

	if recorder.Code >= 200 && recorder.Code < 500 {
		if errDescription != "" {
			log.Instance().Debugf("Handler for %s %s returned %d %s", req.Request.Method, req.Request.RequestURI, recorder.Code, errDescription)
		} else {
			log.Instance().Debugf("Handler for %s %s returned %d", req.Request.Method, req.Request.RequestURI, recorder.Code)
		}
	} else {
		if errDescription != "" {
			log.Instance().Errorf("Handler for %s %s returned %d %s", req.Request.Method, req.Request.RequestURI, recorder.Code, errDescription)
		} else {
			log.Instance().Errorf("Handler for %s %s returned %d", req.Request.Method, req.Request.RequestURI, recorder.Code)
		}
	}
}
