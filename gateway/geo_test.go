package gateway_test

import (
	"encoding/json"
	"github.com/semirm-dev/findhotel/datastore"
	"github.com/semirm-dev/findhotel/gateway"
	"github.com/semirm-dev/findhotel/geo"
	"github.com/semirm-dev/findhotel/internal/web"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetGeoLocation_IpExists(t *testing.T) {
	searchApi := datastore.NewInMemory()
	stored, err := searchApi.Store([]*geo.Geo{
		{
			Ip:           "1.1.1.1",
			CountryCode:  "cc1",
			Country:      "c1",
			City:         "city1",
			Latitude:     123.123,
			Longitude:    123.123,
			MysteryValue: 123,
		},
		{
			Ip:           "2.2.2.2",
			CountryCode:  "cc2",
			Country:      "c2",
			City:         "city2",
			Latitude:     123.123,
			Longitude:    123.123,
			MysteryValue: 123,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, stored)

	router := web.NewRouter()
	router.GET("geo", gateway.GetGeoLocation(searchApi))

	req, _ := http.NewRequest("GET", "/geo?ip=1.1.1.1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	respBody, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fail()
	}

	var resp *geo.Geo
	if err = json.Unmarshal(respBody, &resp); err != nil {
		t.Fail()
	}
	assert.NotNil(t, resp)

	assert.Equal(t, "1.1.1.1", resp.Ip)
	assert.Equal(t, "cc1", resp.CountryCode)
}

func TestGetGeoLocation_IpNotExists(t *testing.T) {
	searchApi := datastore.NewInMemory()

	router := web.NewRouter()
	router.GET("geo", gateway.GetGeoLocation(searchApi))

	req, _ := http.NewRequest("GET", "/geo?ip=5.5.5.5", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	respBody, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fail()
	}

	var resp *geo.Geo
	if err = json.Unmarshal(respBody, &resp); err != nil {
		t.Fail()
	}
	assert.Nil(t, resp)
}
