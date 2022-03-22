package listapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Talis-Aspire-Developer-Community/examples/get-list-title-csv-input/golang/listapi"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assert := assert.New(t)
	m := mux.NewRouter()

	tenant := "broadminster"
	listID := "list_id"

	want := &listapi.GetResponse{
		Data: listapi.GetResponseData{
			Attr: listapi.GetResponseDataAttributes{
				Title: "Expected Title",
			},
		},
	}

	jsonResp := `
		{
			"data": {
				"attributes": {
					"title": "Expected Title"
				}
			}
		}
	`

	urlPath := fmt.Sprintf("/3/%s/lists/%s", tenant, listID)
	m.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResp))
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	a := &listapi.Client{
		BaseURL:    srv.URL,
		TenantCode: tenant,
		Client:     srv.Client(),
	}

	result, err := a.Get(listID)
	assert.NoError(err)
	assert.Equal(want, result)
}

func TestGetError(t *testing.T) {
	assert := assert.New(t)
	m := mux.NewRouter()

	tenant := "broadminster"
	listID := "list_id"

	jsonResp := `
		{
			"data": {
				"attributes": {
					"title": "Expected Title"
				}
			}
		}
	`

	urlPath := fmt.Sprintf("/3/%s/lists/%s", tenant, listID)
	m.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonResp))
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	a := &listapi.Client{
		BaseURL:    srv.URL,
		TenantCode: tenant,
		Client:     srv.Client(),
	}

	var wantResult *listapi.GetResponse
	wantErr := fmt.Errorf("status code was: 500 Internal Server Error")

	result, err := a.Get(listID)
	assert.Equal(wantErr, err)
	assert.Equal(wantResult, result)
}
