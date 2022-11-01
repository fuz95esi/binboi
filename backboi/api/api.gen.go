// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.2 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Address Bin collection point address - as referenced in https://api.reading.gov.uk/v0.json#/components/schemas/Address'
type Address struct {
	AccountSiteId     *string `json:"AccountSiteId,omitempty"`
	AccountSiteUprn   *string `json:"AccountSiteUprn,omitempty"`
	SiteAddress2      *string `json:"SiteAddress2,omitempty"`
	SiteAddressPrefix *string `json:"SiteAddressPrefix,omitempty"`
	SiteEasting       *string `json:"SiteEasting,omitempty"`
	SiteId            *string `json:"SiteId,omitempty"`
	SiteLatitude      *string `json:"SiteLatitude,omitempty"`
	SiteLongitude     *string `json:"SiteLongitude,omitempty"`
	SiteNorthing      *string `json:"SiteNorthing,omitempty"`
	SiteShortAddress  *string `json:"SiteShortAddress,omitempty"`
}

// Addresses List of addresses - as referenced in https://api.reading.gov.uk/v0.json#/components/schemas/Addresses'
type Addresses = []Address

// Collection defines model for Collection.
type Collection struct {
	Date     *string `json:"date,omitempty"`
	Day      *string `json:"day,omitempty"`
	ReadDate *string `json:"read_date,omitempty"`
	Round    *string `json:"round,omitempty"`
	Schedule *string `json:"schedule,omitempty"`
	Service  *string `json:"service,omitempty"`
}

// Collections defines model for Collections.
type Collections struct {
	// Collections List of collections for associated UPRN
	Collections      *[]Collection `json:"collections,omitempty"`
	ErrorCode        *int          `json:"error_code,omitempty"`
	ErrorDescription *string       `json:"error_description,omitempty"`
	Success          *bool         `json:"success,omitempty"`
	Uprn             *string       `json:"uprn,omitempty"`
}

// Status defines model for Status.
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetCollectionDatesParams defines parameters for GetCollectionDates.
type GetCollectionDatesParams struct {
	// FromDate Start date of when collection date list will be generated from (in YYYY-MM-DD format)
	FromDate *string `form:"from_date,omitempty" json:"from_date,omitempty"`

	// ToDate End date of when collection date list will be generated from (in YYYY-MM-DD format)
	ToDate *string `form:"to_date,omitempty" json:"to_date,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /addresses/{postcode})
	GetAddresses(ctx echo.Context, postcode string) error

	// (GET /collections/{uprn})
	GetCollectionDates(ctx echo.Context, uprn string, params GetCollectionDatesParams) error

	// (GET /health)
	Health(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAddresses converts echo context to params.
func (w *ServerInterfaceWrapper) GetAddresses(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "postcode" -------------
	var postcode string

	err = runtime.BindStyledParameterWithLocation("simple", false, "postcode", runtime.ParamLocationPath, ctx.Param("postcode"), &postcode)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter postcode: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetAddresses(ctx, postcode)
	return err
}

// GetCollectionDates converts echo context to params.
func (w *ServerInterfaceWrapper) GetCollectionDates(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "uprn" -------------
	var uprn string

	err = runtime.BindStyledParameterWithLocation("simple", false, "uprn", runtime.ParamLocationPath, ctx.Param("uprn"), &uprn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uprn: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetCollectionDatesParams
	// ------------- Optional query parameter "from_date" -------------

	err = runtime.BindQueryParameter("form", true, false, "from_date", ctx.QueryParams(), &params.FromDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter from_date: %s", err))
	}

	// ------------- Optional query parameter "to_date" -------------

	err = runtime.BindQueryParameter("form", true, false, "to_date", ctx.QueryParams(), &params.ToDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter to_date: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCollectionDates(ctx, uprn, params)
	return err
}

// Health converts echo context to params.
func (w *ServerInterfaceWrapper) Health(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.Health(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/addresses/:postcode", wrapper.GetAddresses)
	router.GET(baseURL+"/collections/:uprn", wrapper.GetCollectionDates)
	router.GET(baseURL+"/health", wrapper.Health)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xXUW/bOAz+K4TugN0Bbmz3kmbNW7bmugFZVyTtQzEMgyozsTZH8iQ6XVDkvx8ku3FS",
	"29cOuB72Flv8yI/kR8q5Z0Kvcq1QkWWje2ZFiivuf46TxKD1PxO0wsicpFZsxN5IBUJnGQr3AnItFQEv",
	"reEIuAWDCzSoBCYgFaREuR2FIc9lzyBPpFr2lnrdK76F66j31Wr1W1izCCsKYRX/FQtYbnSOhiSWxITQ",
	"haK5JHyfuBf4g6/yDNmIDeP+oM8CRpvcPVoyUi3ZNtjHXOdGHaL+iqNoeNwfDNqQDlJROT6Eza/g/P10",
	"Moe304/zyRPYS4ML+ePQwYVW2AWbcEvu8QDQH8av47gXdYEel+OkP4zjLuMpJ0lFgoeQQdzrD+IoOo1O",
	"+v1OqFbLFuxR1Ds9GZ6e9F+fRvFg2IW+0IbSRm7x8Ph00O/ObZ5qQ3uarJHj6dlkBu8+Xs8nAfw9HV9B",
	"HMBhcwKYncdwPJ41vW93b/TtVxTk1VKGwRbxT6Ul0IsHveMLKB695iXhysf/3eCCjdi/IVidBDeGb9zz",
	"292EOieHE5Rw8q1rFDrhm9b3LokvnSijC5W0njieSZG1wyyatRRtZ209qfOxzYTE4WF7y/aMYKENcGu1",
	"kJwwgevL2cVza75X2JayozHafBE62c9LKsIlmvr8gGBbaQohKp1XZ7daZ8h9yKJaYM8o2pw4Fa316uK3",
	"Qmv5sqMpBr8X0mDCRp9KF7X950Z0B5BqoZsNGSsYX773PVgiuTUHt4d3SoLEZWZhYfQKZuUEgdvgQmZu",
	"hCX5yX8j1a2WzhkL2BqNLf1HPbcjtwHTOSqeS7fie+XazDmlvgLhboDD+1xbculs3cESqcl4hlQYZSFr",
	"TD+lnGDFSaRAKcJSrtFdiaVD5ikY7ry43czOkerd4tgYvkJCY9no0+OYl5UTF8+53oupwSAZiWsXQTpj",
	"lxcLmOIrV5e9+HXLyBQYVBd8W3s/O2Oba1XtveMoKqWiCJWvCc/zTAqfTeiWWP3B8Mw9hbaUxSM57BLz",
	"dXSt3iXgjRe8yOg/41LNRAuRszIS+BmFStlVtQuj3C+RSf+ptPXwcG+lhPduLrsldI60k8++0jlhuY6q",
	"HdQQTL1vzpztU7K5StG7ghbtwF2qLTbD3+Erg0Bm42r/pLr8/vkZZQWPOc6JG/KxHcG7FNVjTmWp7mSW",
	"wS3CEpUrCiblQvhDKri5ubk5+vDh6OzMFW/F6c8Hst8LNJuarUOUV9dPUZyo5H8iSPppei85m/sXa8tQ",
	"TJ/U7C8+oynyjNLOuXznj0GkKL4BqsT/nWnMYWnFXrAN3Sm7u9J9K6GB6rtgUWTZBkyh1MN3269a/+32",
	"nwAAAP//ATR6Y2IOAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
