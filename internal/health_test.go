package internal

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHealthCheck(t *testing.T) {
	RegisterHealthCheck()
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/health", nil)
	assert.NoError(t, err)
	http.DefaultServeMux.ServeHTTP(w, r)
	result := w.Result()
	assert.Equal(t, result.StatusCode, 200)
}
