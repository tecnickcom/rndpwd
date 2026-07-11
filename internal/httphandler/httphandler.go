// Package httphandler handles the inbound service requests.
package httphandler

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tecnickcom/nurago/pkg/httpserver"
	"github.com/tecnickcom/nurago/pkg/httputil"
	"github.com/tecnickcom/nurago/pkg/httputil/jsendx"
	"github.com/tecnickcom/nurago/pkg/random"
	"github.com/tecnickcom/rndpwd/internal/metrics"
	"github.com/tecnickcom/rndpwd/internal/password"
	"github.com/tecnickcom/rndpwd/internal/validator"
)

// generator produces random passwords.
type generator interface {
	Generate() ([]string, error)
}

// HTTPHandler is the struct containing all the http handlers.
type HTTPHandler struct {
	httpres     *httputil.HTTPResp
	appInfo     *jsendx.AppInfo
	metric      metrics.Metrics
	val         validator.Validator
	rndpwd      *password.Password
	rnd         *random.Rnd
	newPassword func(charset string, length, quantity int) generator
}

// New creates a new instance of the HTTP handler.
func New(l *slog.Logger, appInfo *jsendx.AppInfo, metric metrics.Metrics, val validator.Validator, rndpwd *password.Password) *HTTPHandler {
	return &HTTPHandler{
		httpres: httputil.NewHTTPResp(l),
		appInfo: appInfo,
		metric:  metric,
		val:     val,
		rndpwd:  rndpwd,
		rnd:     random.New(nil),
		newPassword: func(charset string, length, quantity int) generator {
			return password.New(charset, length, quantity)
		},
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
	h.httpres.SendJSON(r.Context(), w, http.StatusOK, h.rnd.UUIDv7().String())
}

func (h *HTTPHandler) handlePassword(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if !validQueryParams(query) {
		h.httpres.SendJSON(r.Context(), w, http.StatusBadRequest, "invalid query parameter")
		return
	}

	// URL query parameters can override the config settings
	p := h.newPassword(
		httputil.QueryStringOrDefault(query, "charset", h.rndpwd.Charset),
		httputil.QueryIntOrDefault(query, "length", h.rndpwd.Length),
		httputil.QueryIntOrDefault(query, "quantity", h.rndpwd.Quantity),
	)

	err := h.val.ValidateStruct(p)
	if err != nil {
		h.httpres.SendJSON(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	pwds, err := p.Generate()
	if err != nil {
		h.httpres.SendJSON(r.Context(), w, http.StatusInternalServerError, "failed generating passwords")
		return
	}

	h.httpres.SendJSON(r.Context(), w, http.StatusOK, pwds)
}

// validQueryParams reports whether the request query only contains the allowed
// single-valued parameters, with integer values where required.
func validQueryParams(query url.Values) bool {
	validParam := map[string]bool{
		"charset":  true,
		"length":   true,
		"quantity": true,
	}

	for param := range query {
		if !validParam[param] || len(query[param]) > 1 || query.Get(param) == "" || (param != "charset" && !isInteger(query.Get(param))) {
			return false
		}
	}

	return true
}

func isInteger(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}
