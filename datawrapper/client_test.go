package datawrapper

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	. "github.com/smartystreets/goconvey/convey"
)

type MockedHTTPResponse struct {
	StatusCode int
	Body       string
}

const (
	apiName  = "datawrapper"
	apiToken = "token01"
)

var ctx = context.Background()

func getMockAPI(expectRequest http.Request, mockedHTTPResponse MockedHTTPResponse) *Client {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectRequest.Method {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		expectedAuth, ok := expectRequest.Header["Authorization"]
		if ok && r.Header.Get("Authorization") != expectedAuth[0] {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(mockedHTTPResponse.StatusCode)
		fmt.Fprintln(w, mockedHTTPResponse.Body)
	}))

	api := NewClient(ts.URL, apiToken)

	return api
}

func TestClient_Checker(t *testing.T) {
	initialTime := time.Now().UTC()

	Convey("When DW /me endpoint returns status OK", t, func() {
		mockedAPI := getMockAPI(
			http.Request{Method: "GET", Header: map[string][]string{
				"Authorization": {"Bearer " + apiToken},
			}},
			MockedHTTPResponse{StatusCode: 200, Body: "{\"status\": \"OK\"}"},
		)

		check := health.NewCheckState(apiName)

		err := mockedAPI.Checker(ctx, check)
		So(check.Name(), ShouldEqual, apiName)
		So(check.StatusCode(), ShouldEqual, 200)
		So(check.Status(), ShouldEqual, health.StatusOK)
		So(check.Message(), ShouldEqual, "datawrapper is ok")
		So(*check.LastChecked(), ShouldHappenAfter, initialTime)
		So(check.LastFailure(), ShouldBeNil)
		So(*check.LastSuccess(), ShouldHappenAfter, initialTime)
		So(err, ShouldBeNil)
	})

	Convey("When DW /me endpoint returns status Forbidden", t, func() {
		mockedAPI := getMockAPI(
			http.Request{Method: "GET", Header: map[string][]string{
				"Authorization": {"INVALID AUTH"},
			}},
			MockedHTTPResponse{},
		)

		check := health.NewCheckState(apiName)

		err := mockedAPI.Checker(ctx, check)
		So(check.Name(), ShouldEqual, apiName)
		So(check.StatusCode(), ShouldEqual, 403)
		So(check.Status(), ShouldEqual, health.StatusWarning)
		So(check.Message(), ShouldEqual, "datawrapper is degraded, but at least partially functioning")
		So(*check.LastChecked(), ShouldHappenAfter, initialTime)
		So(*check.LastFailure(), ShouldHappenAfter, initialTime)
		So(check.LastSuccess(), ShouldBeNil)
		So(err, ShouldBeNil)
	})

	Convey("When DW /me endpoint returns unexpected status", t, func() {
		mockedAPI := getMockAPI(
			http.Request{Method: "GET", Header: map[string][]string{
				"Authorization": {"Bearer " + apiToken},
			}},
			MockedHTTPResponse{StatusCode: 500},
		)

		check := health.NewCheckState(apiName)

		err := mockedAPI.Checker(ctx, check)
		So(check.Name(), ShouldEqual, apiName)
		So(check.StatusCode(), ShouldEqual, 500)
		So(check.Status(), ShouldEqual, health.StatusCritical)
		So(check.Message(), ShouldEqual, "datawrapper functionality is unavailable or non-functioning")
		So(*check.LastChecked(), ShouldHappenAfter, initialTime)
		So(*check.LastFailure(), ShouldHappenAfter, initialTime)
		So(check.LastSuccess(), ShouldBeNil)
		So(err, ShouldBeNil)
	})

	Convey("When DW /me endpoint doesn't respond", t, func() {
		client := NewClient("", "")

		check := health.NewCheckState(apiName)

		err := client.Checker(ctx, check)
		So(check.Name(), ShouldEqual, apiName)
		So(check.StatusCode(), ShouldEqual, 0)
		So(check.Status(), ShouldEqual, health.StatusCritical)
		So(check.Message(), ShouldEqual, `Get "/v3/me": unsupported protocol scheme ""`)
		So(*check.LastChecked(), ShouldHappenAfter, initialTime)
		So(*check.LastFailure(), ShouldHappenAfter, initialTime)
		So(check.LastSuccess(), ShouldBeNil)
		So(err, ShouldBeNil)
	})
}
