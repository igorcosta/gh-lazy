package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key":"value"}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token")
	require.NoError(t, err)
	client.client.baseURL = server.URL

	var response map[string]string
	err = client.Get(context.Background(), "/test", &response)
	require.NoError(t, err)
	assert.Equal(t, "value", response["key"])
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":1}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token")
	require.NoError(t, err)
	client.client.baseURL = server.URL

	var response struct {
		ID int `json:"id"`
	}
	err = client.Post(context.Background(), "/test", map[string]string{"key": "value"}, &response)
	require.NoError(t, err)
	assert.Equal(t, 1, response.ID)
}

func TestClient_Patch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"updated":true}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token")
	require.NoError(t, err)
	client.client.baseURL = server.URL

	var response struct {
		Updated bool `json:"updated"`
	}
	err = client.Patch(context.Background(), "/test", map[string]string{"key": "value"}, &response)
	require.NoError(t, err)
	assert.True(t, response.Updated)
}
