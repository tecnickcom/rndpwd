package httphandler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tecnickcom/gogen/pkg/testutil"
	"github.com/tecnickcom/rndpwd/internal/password"
	"github.com/tecnickcom/rndpwd/internal/validator"
)

func TestNew(t *testing.T) {
	t.Parallel()

	hh := New(nil, nil, nil, nil)
	require.NotNil(t, hh)
}

func TestHTTPHandler_BindHTTP(t *testing.T) {
	t.Parallel()

	h := &HTTPHandler{}
	got := h.BindHTTP(testutil.Context())
	require.Len(t, got, 2)
}

func TestHTTPHandler_handleGenUID(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(testutil.Context(), http.MethodGet, "/", nil)

	(&HTTPHandler{}).handleGenUID(rr, req)

	resp := rr.Result()
	require.NotNil(t, resp)

	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err, "error closing resp.Body")
	}()

	body, _ := io.ReadAll(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	require.NotEmpty(t, string(body))
}

func TestHTTPHandler_handlePassword(t *testing.T) {
	t.Parallel()

	val, _ := validator.New("json")

	h := New(
		nil,
		nil,
		val,
		password.New("0123456789abcdefghijklmnopqrstuvwxyz", 16, 3),
	)

	tests := []struct {
		name    string
		params  string
		wantErr bool
	}{
		{
			name:    "valid empty",
			params:  "",
			wantErr: false,
		},
		{
			name:    "valid length",
			params:  "?length=16",
			wantErr: false,
		},
		{
			name:    "valid quantity",
			params:  "?quantity=3",
			wantErr: false,
		},
		{
			name:    "valid charset",
			params:  "?charset=0123456789abcdefghijklmnopqrstuvwxyz",
			wantErr: false,
		},
		{
			name:    "valid all params",
			params:  "?charset=0123456789&length=8&quantity=1",
			wantErr: false,
		},
		{
			name:    "invalid charset",
			params:  "?charset=in va lid",
			wantErr: true,
		},
		{
			name:    "invalid parameter",
			params:  "?invalid=3",
			wantErr: true,
		},
		{
			name:    "empty parameter",
			params:  "?quantity=",
			wantErr: true,
		},
		{
			name:    "not integer length",
			params:  "?length=abc",
			wantErr: true,
		},
		{
			name:    "not integer quantity",
			params:  "?quantity=abc",
			wantErr: true,
		},
		{
			name:    "duplicate parameter",
			params:  "?quantity=1&quatity=2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(testutil.Context(), http.MethodGet, "/"+tt.params, nil)

			h.handlePassword(rr, req)

			resp := rr.Result()
			require.NotNil(t, resp)

			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err, "error closing resp.Body")
			}()

			body, _ := io.ReadAll(resp.Body)

			if tt.wantErr {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			} else {
				require.Equal(t, http.StatusOK, resp.StatusCode)
				require.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
				require.NotEmpty(t, string(body))
			}
		})
	}
}
