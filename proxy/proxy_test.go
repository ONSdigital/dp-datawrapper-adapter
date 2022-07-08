package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func MustParse(rawurl string) *url.URL {
	url, err := url.Parse(rawurl)
	if err != nil {
		panic("invalid url: " + err.Error())
	}
	return url
}
func NewRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic("invalid request: " + err.Error())
	}
	return req
}

func TestDirector(t *testing.T) {
	Convey("Director correctly trims the router path", t, func() {
		director := director("/api", MustParse("https://api.datawrapper.de"))

		req := NewRequest("GET", "http://1.2.3.4/api/v3/auth", nil)
		director(req)
		So(req.URL.String(), ShouldEqual, "https://api.datawrapper.de/v3/auth")

		req = NewRequest("GET", "http://1.2.3.4/v3/auth", nil)
		director(req)
		So(req.URL.String(), ShouldEqual, "https://api.datawrapper.de/v3/auth")

		req = NewRequest("GET", "http://1.2.3.4/v3/auth/api", nil)
		director(req)
		So(req.URL.String(), ShouldEqual, "https://api.datawrapper.de/v3/auth/api")
	})

	Convey("Director doesn't trim the path if set up under /", t, func() {
		director := director("/", MustParse("https://api.datawrapper.de"))

		req := NewRequest("GET", "http://1.2.3.4/api/v3/auth", nil)
		director(req)
		So(req.URL.String(), ShouldEqual, "https://api.datawrapper.de/api/v3/auth")

		req = NewRequest("GET", "http://1.2.3.4/v3/auth", nil)
		director(req)
		So(req.URL.String(), ShouldEqual, "https://api.datawrapper.de/v3/auth")

		req = NewRequest("GET", "http://1.2.3.4/v3/auth/api", nil)
		director(req)
		So(req.URL.String(), ShouldEqual, "https://api.datawrapper.de/v3/auth/api")
	})

	Convey("Director doesn't trim the path if set up under empty router path", t, func() {
		director := director("", MustParse("https://api.datawrapper.de"))

		req := NewRequest("GET", "http://1.2.3.4/api/v3/auth", nil)
		director(req)
		So(req.URL.String(), ShouldEqual, "https://api.datawrapper.de/api/v3/auth")

		req = NewRequest("GET", "http://1.2.3.4/v3/auth", nil)
		director(req)
		So(req.URL.String(), ShouldEqual, "https://api.datawrapper.de/v3/auth")

		req = NewRequest("GET", "http://1.2.3.4/v3/auth/api", nil)
		director(req)
		So(req.URL.String(), ShouldEqual, "https://api.datawrapper.de/v3/auth/api")
	})
}

func TestForwarding(t *testing.T) {
	Convey("Proxy correctly forwards the requests", t, func() {
		const backendResponse = "I am the backend"
		const backendStatus = 201
		backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Convey("URL path is trimmed as expected", t, func() {
				So(r.URL.Path, ShouldEqual, "/v3/chart")
			})
			w.WriteHeader(backendStatus)
			_, _ = w.Write([]byte(backendResponse))
		}))
		defer backend.Close()

		proxyHandler, err := New("/api", backend.URL)
		So(err, ShouldBeNil)
		frontend := httptest.NewServer(proxyHandler)
		defer frontend.Close()

		getReq, _ := http.NewRequest("GET", frontend.URL+"/api/v3/chart", nil)
		getReq.Host = "some-name"
		getReq.Header.Set("Connection", "close")
		getReq.Close = true
		res, err := frontend.Client().Do(getReq)
		So(err, ShouldBeNil)
		So(res.StatusCode, ShouldEqual, backendStatus)

		bodyBytes, err := io.ReadAll(res.Body)
		So(err, ShouldBeNil)
		So(string(bodyBytes), ShouldEqual, backendResponse)
	})
}
