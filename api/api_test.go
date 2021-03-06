package api

import (
	"encoding/json"
	"github.com/elferink/swarm/cluster"
	"github.com/elferink/swarm/version"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func serveRequest(c cluster.Cluster, w http.ResponseWriter, req *http.Request) error {
	context := &context{
		cluster: c,
	}

	r := createRouter(context, false)
	r.ServeHTTP(w, req)
	return nil
}

func TestGetVersion(t *testing.T) {
	r := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/version", nil)
	assert.NoError(t, err)

	assert.NoError(t, serveRequest(nil, r, req))
	assert.Equal(t, r.Code, http.StatusOK)

	v := struct {
		Version string
	}{}

	json.NewDecoder(r.Body).Decode(&v)
	assert.Equal(t, v.Version, "swarm/"+version.VERSION)
}
