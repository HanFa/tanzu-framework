// Code generated by go-swagger; DO NOT EDIT.

package aws

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/vmware-tanzu/tanzu-framework/tkg/web/server/models"
)

// GetVPCsOKCode is the HTTP code returned for type GetVPCsOK
const GetVPCsOKCode int = 200

/*
GetVPCsOK Successful retrieval of AWS VPCs

swagger:response getVPCsOK
*/
type GetVPCsOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Vpc `json:"body,omitempty"`
}

// NewGetVPCsOK creates GetVPCsOK with default headers values
func NewGetVPCsOK() *GetVPCsOK {

	return &GetVPCsOK{}
}

// WithPayload adds the payload to the get v p cs o k response
func (o *GetVPCsOK) WithPayload(payload []*models.Vpc) *GetVPCsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v p cs o k response
func (o *GetVPCsOK) SetPayload(payload []*models.Vpc) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVPCsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Vpc, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetVPCsBadRequestCode is the HTTP code returned for type GetVPCsBadRequest
const GetVPCsBadRequestCode int = 400

/*
GetVPCsBadRequest Bad request

swagger:response getVPCsBadRequest
*/
type GetVPCsBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetVPCsBadRequest creates GetVPCsBadRequest with default headers values
func NewGetVPCsBadRequest() *GetVPCsBadRequest {

	return &GetVPCsBadRequest{}
}

// WithPayload adds the payload to the get v p cs bad request response
func (o *GetVPCsBadRequest) WithPayload(payload *models.Error) *GetVPCsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v p cs bad request response
func (o *GetVPCsBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVPCsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetVPCsUnauthorizedCode is the HTTP code returned for type GetVPCsUnauthorized
const GetVPCsUnauthorizedCode int = 401

/*
GetVPCsUnauthorized Incorrect credentials

swagger:response getVPCsUnauthorized
*/
type GetVPCsUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetVPCsUnauthorized creates GetVPCsUnauthorized with default headers values
func NewGetVPCsUnauthorized() *GetVPCsUnauthorized {

	return &GetVPCsUnauthorized{}
}

// WithPayload adds the payload to the get v p cs unauthorized response
func (o *GetVPCsUnauthorized) WithPayload(payload *models.Error) *GetVPCsUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v p cs unauthorized response
func (o *GetVPCsUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVPCsUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetVPCsInternalServerErrorCode is the HTTP code returned for type GetVPCsInternalServerError
const GetVPCsInternalServerErrorCode int = 500

/*
GetVPCsInternalServerError Internal server error

swagger:response getVPCsInternalServerError
*/
type GetVPCsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetVPCsInternalServerError creates GetVPCsInternalServerError with default headers values
func NewGetVPCsInternalServerError() *GetVPCsInternalServerError {

	return &GetVPCsInternalServerError{}
}

// WithPayload adds the payload to the get v p cs internal server error response
func (o *GetVPCsInternalServerError) WithPayload(payload *models.Error) *GetVPCsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v p cs internal server error response
func (o *GetVPCsInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVPCsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
