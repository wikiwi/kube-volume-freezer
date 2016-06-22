package controllers

import (
	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/volumes"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
	"github.com/wikiwi/kube-volume-freezer/pkg/validation"
)

func NewVolume(authFilter restful.FilterFunction, m volumes.Manager) *Volume {
	return &Volume{authFilter: authFilter, manager: m}
}

// Volume is a REST resource for reporting health status.
type Volume struct {
	authFilter restful.FilterFunction
	manager    volumes.Manager
}

// Register adds this resource to the provided container.
func (r *Volume) Register(s *rest.Server) {
	authFilter := r.authFilter
	if r.authFilter == nil {
		authFilter = rest.NoOpFilter
	}
	ws := new(restful.WebService)
	ws.Path("/volumes").
		Doc("list, freeze or thaw Kubernetes Pod Volumes").
		Produces(restful.MIME_JSON).
		Filter(authFilter)
	ws.Route(ws.GET("/{podUID}").To(r.getVolumeList).
		Doc("list Volumes").
		Param(ws.PathParameter("podUID", "UID of Pod to show its Volumes.").DataType("string")).
		Writes(api.VolumeList{}).
		Returns(404, "Not Found", api.Error{}))
	ws.Route(ws.GET("/{podUID}/{name}").To(r.getVolume).
		Doc("get Volume").
		Param(ws.PathParameter("podUID", "UID of the Pod that the Volume belongs to.").DataType("string")).
		Param(ws.PathParameter("name", "Name of Volume to request info from.").DataType("string")).
		Writes(api.Volume{}).
		Returns(404, "Not Found", api.Error{}))
	ws.Route(ws.POST("/{podUID}/{name}").To(r.freezeThaw).
		Doc("freeze or thaw Volumes").
		Param(ws.PathParameter("podUID", "UID of the Pod that the Volume belongs to.").DataType("string")).
		Param(ws.PathParameter("name", "Name of Volume to request info from.").DataType("string")).
		Reads(api.FreezeThawRequest{}).
		Writes(api.Volume{}).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Returns(404, "Not Found", api.Error{}))
	s.Register(ws)
}

func (r Volume) getVolumeList(request *restful.Request, response *restful.Response) {
	uid := request.PathParameter("podUID")
	issues := validation.ValidateUIDParameter("podUID", uid)
	if len(issues) > 0 {
		rest.WriteValidationError(issues, response)
		return
	}

	list, err := r.manager.List(uid)
	rest.RespondOrDie(200, list, err, response)
}

func (r Volume) getVolume(request *restful.Request, response *restful.Response) {
	uid := request.PathParameter("podUID")
	name := request.PathParameter("name")

	issues := validation.ValidateUIDParameter("podUID", uid)
	issues = append(issues, validation.ValidateQualifiedNameParameter("name", name)...)
	if len(issues) > 0 {
		rest.WriteValidationError(issues, response)
		return
	}
	volume, err := r.manager.Get(uid, name)
	rest.RespondOrDie(200, volume, err, response)
}

func (r Volume) freezeThaw(request *restful.Request, response *restful.Response) {
	uid := request.PathParameter("podUID")
	name := request.PathParameter("name")
	issues := validation.ValidateUIDParameter("podUID", uid)
	issues = append(issues, validation.ValidateQualifiedNameParameter("name", name)...)
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
	var volume *api.Volume
	if ftr.Action == "freeze" {
		volume, err = r.manager.Freeze(uid, name)
	} else if ftr.Action == "thaw" {
		volume, err = r.manager.Thaw(uid, name)
	} else {
		panic("unknown action " + ftr.Action)
	}

	rest.RespondOrDie(200, volume, err, response)
}
