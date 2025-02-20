// Code generated by go-swagger; DO NOT EDIT.

package features

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetFeatureFlagsHandlerFunc turns a function with the right signature into a get feature flags handler
type GetFeatureFlagsHandlerFunc func(GetFeatureFlagsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetFeatureFlagsHandlerFunc) Handle(params GetFeatureFlagsParams) middleware.Responder {
	return fn(params)
}

// GetFeatureFlagsHandler interface for that can handle valid get feature flags params
type GetFeatureFlagsHandler interface {
	Handle(GetFeatureFlagsParams) middleware.Responder
}

// NewGetFeatureFlags creates a new http.Handler for the get feature flags operation
func NewGetFeatureFlags(ctx *middleware.Context, handler GetFeatureFlagsHandler) *GetFeatureFlags {
	return &GetFeatureFlags{Context: ctx, Handler: handler}
}

/*
GetFeatureFlags swagger:route GET /api/features features getFeatureFlags

Retrieve list of features
*/
type GetFeatureFlags struct {
	Context *middleware.Context
	Handler GetFeatureFlagsHandler
}

func (o *GetFeatureFlags) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetFeatureFlagsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
