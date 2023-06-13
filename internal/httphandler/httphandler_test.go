package httphandler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vonage/gosrvlib/pkg/testutil"
	"github.com/stretchr/testify/require"
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
	require.Equal(t, 2, len(got))
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
		charset string
		wantErr bool
	}{
		{
			name:    "valid",
			wantErr: false,
		},
		{
			name:    "invalid",
			charset: "in va lid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(testutil.Context(), http.MethodGet, "/?charset="+tt.charset, nil)

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
