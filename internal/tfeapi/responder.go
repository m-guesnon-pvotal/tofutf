package tfeapi

import (
	"encoding/json"
	"fmt"
	tfetypes "github.com/hashicorp/go-tfe"
	tfejsonapi "github.com/hashicorp/jsonapi"
	"github.com/tofutf/tofutf/internal"
	"net/http"

	"github.com/DataDog/jsonapi"
	"github.com/tofutf/tofutf/internal/resource"
)

const mediaType = "application/vnd.api+json"

// Responder handles responding to API requests.
type Responder struct {
	*includer
}

func NewResponder() *Responder {
	return &Responder{
		includer: &includer{
			registrations: make(map[IncludeName][]IncludeFunc),
		},
	}
}

func (res *Responder) RespondWithPage(w http.ResponseWriter, r *http.Request, items any, pagination *resource.Pagination) {
	meta := jsonapi.MarshalMeta(map[string]*resource.Pagination{
		"pagination": pagination,
	})
	res.Respond(w, r, items, http.StatusOK, meta)
}

func (res *Responder) Respond(w http.ResponseWriter, r *http.Request, payload any, status int, opts ...jsonapi.MarshalOption) {
	includes, err := res.addIncludes(r, payload)
	if err != nil {
		Error(w, err)
		return
	}
	if len(includes) > 0 {
		opts = append(opts, jsonapi.MarshalInclude(includes...))
	}
	b, err := jsonapi.Marshal(payload, opts...)
	if err != nil {
		Error(w, err)
		return
	}
	res.RespondRaw(w, r, b, status)
}

func (res *Responder) JsonAPIResponse(w http.ResponseWriter, r *http.Request, status int, payload any, pagination *resource.Pagination) {
	p, err := tfejsonapi.Marshal(payload)
	if err != nil {
		Error(w, err)
		return
	}
	var b []byte
	switch typ := p.(type) {
	case *tfejsonapi.OnePayload:
		b, err = json.Marshal(typ)
		if err != nil {
			Error(w, &internal.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to marshal: %s", err.Error()),
			})
			return
		}
	case *tfejsonapi.ManyPayload:
		meta := make(map[string]interface{})
		tfePagination := &tfetypes.Pagination{}
		if pagination != nil {
			tfePagination.TotalCount = pagination.TotalCount
			tfePagination.TotalPages = pagination.TotalPages
			tfePagination.CurrentPage = pagination.CurrentPage
			if pagination.NextPage != nil {
				tfePagination.NextPage = *pagination.NextPage
			} else {
				tfePagination.NextPage = pagination.CurrentPage
			}
			if pagination.PreviousPage != nil {
				tfePagination.PreviousPage = *pagination.PreviousPage
			} else {
				tfePagination.PreviousPage = pagination.CurrentPage
			}
		}
		meta["pagination"] = tfePagination
		computedMeta := tfejsonapi.Meta(meta)
		typ.Meta = &computedMeta
		b, err = json.Marshal(typ)
		if err != nil {
			Error(w, &internal.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to marshal: %s", err.Error()),
			})
			return
		}
	default:
		Error(w, &internal.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to marshal: unsupported payload type: %T", typ),
		})
		return
	}

	res.RespondRaw(w, r, b, status)
}

func (res *Responder) RespondRaw(w http.ResponseWriter, r *http.Request, payload []byte, status int) {
	w.Header().Set("Content-type", mediaType)
	w.WriteHeader(status)
	w.Write(payload) //nolint:errcheck
}
