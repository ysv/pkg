package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeJSONBody(t *testing.T) {
	var fooBar struct {
		Foo string `json:"foo"`
	}

	tests := []struct {
		name    string
		request *http.Request
		dst     interface{}
		wantErr *APIError
	}{
		{
			name:    "ok",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`)),
			dst:     &fooBar,
			wantErr: nil,
		},
		{
			name:    "ok_empty_json",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`)),
			dst:     &fooBar,
			wantErr: nil,
		},
		{
			name:    "empty_body",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(``)),
			dst:     &fooBar,
			wantErr: NewAPIError(http.StatusBadRequest, "body.empty"),
		},
		{
			name:    "multiple_json",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{},{}`)),
			dst:     &fooBar,
			wantErr: NewAPIError(http.StatusBadRequest, "body.invalid_json"),
		},
		{
			name:    "invalid_json",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{`)),
			dst:     &fooBar,
			wantErr: NewAPIError(http.StatusBadRequest, "body.invalid_json"),
		},
		{
			name:    "invalid_json",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":"}`)),
			dst:     &fooBar,
			wantErr: NewAPIError(http.StatusBadRequest, "body.invalid_json"),
		},
		{
			name:    "invalid_type",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":{}}`)),
			dst:     &fooBar,
			wantErr: NewAPIError(http.StatusBadRequest, "foo.invalid_type"),
		},
		{
			name:    "unknown_field",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"foo":"bar", "bar": 1}`)),
			dst:     &fooBar,
			wantErr: NewAPIError(http.StatusBadRequest, "bar.unknown_field"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			err := DecodeJSONBody(rr, tt.request, tt.dst)

			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tt.wantErr, err)
			}
		})
	}

	t.Run("unsupported_content_type", func(t *testing.T) {
		rr := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "text/html")

		err := DecodeJSONBody(rr, req, fooBar)
		require.NotNil(t, err)
		require.Equal(t, NewAPIError(http.StatusUnsupportedMediaType, "content_type.unsupported"), err)
	})
}
