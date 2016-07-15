package rest

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/api/errors"
	"github.com/wikiwi/kube-volume-freezer/pkg/log"
)

// RecoverFilter captures a panic and responds with an API Error.
func RecoverFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	defer func() {
		if r := recover(); r != nil {
			panicHandler(r, req, resp)
			return
		}
	}()
	chain.ProcessFilter(req, resp)
}

func panicHandler(err interface{}, req *restful.Request, resp *restful.Response) {
	var buffer bytes.Buffer
	buffer.WriteString("panic recovered\n")
	for i := 2; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buffer.WriteString(fmt.Sprintf("    %s:%d\n", file, line))
	}
	log.Instance().WithField("error", err).Error(buffer.String())
	er := errors.Unexpected(fmt.Sprintf("%v", err))
	err = resp.WriteHeaderAndEntity(er.Code, er)
	if err != nil {
		log.Instance().WithField("error", err).Error(err)
	}
}
