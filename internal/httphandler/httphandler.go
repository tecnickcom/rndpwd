// Package httphandler handles the inbound service requests.
package httphandler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Vonage/gosrvlib/pkg/httpserver"
	"github.com/Vonage/gosrvlib/pkg/httputil"
	"github.com/Vonage/gosrvlib/pkg/httputil/jsendx"
	"github.com/Vonage/gosrvlib/pkg/uidc"
	"github.com/tecnickcom/rndpwd/internal/metrics"
	"github.com/tecnickcom/rndpwd/internal/password"
	"github.com/tecnickcom/rndpwd/internal/validator"
)

// Service is the interface representing the business logic of the service.
type Service interface {
	// NOTE
	// This is a sample Service interface.
	// It is meant to demonstrate where the business logic of a service should reside.
	// It adds the capability of mocking the HTTP Handler independently from the rest of the code.
	// Add service functions here.
}

// HTTPHandler is the struct containing all the http handlers.
type HTTPHandler struct {
	appInfo *jsendx.AppInfo
	metric  metrics.Metrics
	val     validator.Validator
	rndpwd  *password.Password
}

// New creates a new instance of the HTTP handler.
func New(appInfo *jsendx.AppInfo, metric metrics.Metrics, val validator.Validator, rndpwd *password.Password) *HTTPHandler {
	return &HTTPHandler{
		appInfo: appInfo,
		metric:  metric,
		val:     val,
		rndpwd:  rndpwd,
	}
}

// BindHTTP implements the function to bind the handler to a server.
func (h *HTTPHandler) BindHTTP(_ context.Context) []httpserver.Route {
	return []httpserver.Route{
		{
			Method:      http.MethodGet,
			Path:        "/password",
			Handler:     h.handlePassword,
			Description: "Returns random passwords; charset, length and quantity can be specified as query parameters",
		},
		{
			Method:      http.MethodGet,
			Path:        "/uid",
			Handler:     h.handleGenUID,
			Description: "Generates a random UID",
		},
	}
}

func (h *HTTPHandler) handleGenUID(w http.ResponseWriter, r *http.Request) {
	httputil.SendJSON(r.Context(), w, http.StatusOK, uidc.NewID128())
}

func (h *HTTPHandler) handlePassword(w http.ResponseWriter, r *http.Request) {
	validParam := map[string]bool{
		"charset":  true,
		"length":   true,
		"quantity": true,
	}

	query := r.URL.Query()

	for param := range query {
		if !validParam[param] || (len(query[param]) > 1) || (query.Get(param) == "") || (param != "charset" && !isInteger(query.Get(param))) {
			httputil.SendJSON(r.Context(), w, http.StatusBadRequest, "invalid query parameter")
			return
		}
	}

	// URL query parameters can override the config settings
	p := password.New(
		httputil.QueryStringOrDefault(query, "charset", h.rndpwd.Charset),
		httputil.QueryIntOrDefault(query, "length", h.rndpwd.Length),
		httputil.QueryIntOrDefault(query, "quantity", h.rndpwd.Quantity),
	)

	err := h.val.ValidateStruct(p)
	if err != nil {
		httputil.SendJSON(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.SendJSON(r.Context(), w, http.StatusOK, p.Generate())
}

func isInteger(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}
