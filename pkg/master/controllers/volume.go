package controllers

import (
	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/master/volumes"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
	"github.com/wikiwi/kube-volume-freezer/pkg/validation"
)

// NewVolume creates a Volume Controller.
func NewVolume(authFilter restful.FilterFunction, manager volumes.Manager) *Volume {
	return &Volume{authFilter: authFilter, manager: manager}
}

// Volume is a REST Controller for the Volume Resource.
type Volume struct {
	authFilter restful.FilterFunction
	manager    volumes.Manager
}

// Register adds this resource to the provided REST Server.
func (r Volume) Register(s *rest.Server) {
	authFilter := r.authFilter
	if r.authFilter == nil {
		authFilter = rest.NoOpFilter
	}
	ws := new(restful.WebService)
	ws.Path("/volumes").
		Doc("list, freeze or thaw Kubernetes Pod Volumes").
		Produces(restful.MIME_JSON).
		Filter(authFilter)
	ws.Route(ws.GET("/{namespace}/{podName}").To(r.getVolumeList).
		Doc("list Volumes").
		Param(ws.PathParameter("namespace", "Namespace of Pod.").DataType("string")).
		Param(ws.PathParameter("podName", "Name of Pod.").DataType("string")).
		Writes(api.VolumeList{}).
		Returns(404, "Not Found", api.Error{}))
	ws.Route(ws.GET("/{namespace}/{podName}/{volName}").To(r.getVolume).
		Doc("get Volume").
		Param(ws.PathParameter("namespace", "Namespace of Pod.").DataType("string")).
		Param(ws.PathParameter("podName", "Name of Pod.").DataType("string")).
		Param(ws.PathParameter("volName", "Name of Volume.").DataType("string")).
		Writes(api.VolumeList{}).
		Returns(404, "Not Found", api.Error{}))
	ws.Route(ws.POST("/{namespace}/{podName}/{volName}").To(r.freezeThaw).
		Doc("freeze or thaw Volume").
		Param(ws.PathParameter("namespace", "Namespace of Pod.").DataType("string")).
		Param(ws.PathParameter("podName", "Name of Pod.").DataType("string")).
		Param(ws.PathParameter("volName", "Name of Volume.").DataType("string")).
		Writes(api.VolumeList{}).
		Returns(404, "Not Found", api.Error{}))
	s.Register(ws)
}

func (r Volume) getVolumeList(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("podName")

	issues := validation.ValidateQualifiedNameParameter("namespace", namespace)
	issues = append(issues, validation.ValidateQualifiedNameParameter("podName", podName)...)
	if len(issues) > 0 {
		rest.WriteValidationError(issues, response)
		return
	}

	list, err := r.manager.List(namespace, podName)
	rest.RespondOrDie(200, list, err, response)
}

func (r Volume) getVolume(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("podName")
	volName := request.PathParameter("volName")

	issues := validation.ValidateQualifiedNameParameter("namespace", namespace)
	issues = append(issues, validation.ValidateQualifiedNameParameter("podName", podName)...)
	issues = append(issues, validation.ValidateQualifiedNameParameter("volName", volName)...)
	if len(issues) > 0 {
		rest.WriteValidationError(issues, response)
		return
	}

	vol, err := r.manager.Get(namespace, podName, volName)
	rest.RespondOrDie(200, vol, err, response)
}

func (r Volume) freezeThaw(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	podName := request.PathParameter("podName")
	volName := request.PathParameter("volName")

	issues := validation.ValidateQualifiedNameParameter("namespace", namespace)
	issues = append(issues, validation.ValidateQualifiedNameParameter("podName", podName)...)
	issues = append(issues, validation.ValidateQualifiedNameParameter("volName", volName)...)
	if len(issues) > 0 {
		rest.WriteValidationError(issues, response)
		return
	}

	var ftr api.FreezeThawRequest
	if success := rest.ReadEntityOrBadRequest(&ftr, request, response); !success {
		return
	}

	issues = validation.ValidateFreezeThawRequest(&ftr)
	if len(issues) > 0 {
		rest.WriteValidationError(issues, response)
		return
	}

	var err error
	var vol *api.Volume
	if ftr.Action == "freeze" {
		vol, err = r.manager.Freeze(namespace, podName, volName)
	} else if ftr.Action == "thaw" {
		vol, err = r.manager.Thaw(namespace, podName, volName)
	} else {
		panic("unknown action " + ftr.Action)
	}

	rest.RespondOrDie(200, vol, err, response)
}
