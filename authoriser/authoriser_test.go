package authoriser

// import (
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/ONSdigital/dp-authorisation/v2/permissions"
// 	"github.com/ONSdigital/dp-datawrapper-adapter/authoriser/mocks"
// 	. "github.com/smartystreets/goconvey/convey"
// )

// var (
// 	tpMock = &mocks.TokenParserMock{
// 		ParseFunc: func(tokenString string) (*permissions.EntityData, error) { return &permissions.EntityData{}, nil },
// 	}
// 	csMock = &mocks.ChartStoreMock{
// 		GetCollectionIDFunc: func(chartID string) (string, error) { return "coll01", nil },
// 	}
// 	pcAllowMock = &mocks.PermissionsCheckerMock{
// 		HasPermissionFunc: func(ctx context.Context, entityData permissions.EntityData, permission string, attributes map[string]string) (bool, error) {
// 			return true, nil
// 		},
// // 	}
// // 	pcDenyMock = &mocks.PermissionsCheckerMock{
// // 		HasPermissionFunc: func(ctx context.Context, entityData permissions.EntityData, permission string, attributes map[string]string) (bool, error) {
// // 			return false, nil
// // 		},
// // 	}
// // )

// func TestMiddleware(t *testing.T) {
// 	Convey("Middleware", t, func() {
// 		successHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			http.Error(w, "success", http.StatusOK)
// 		})

// 		Convey("returns forbidden status if the request is not chart related", func() {
// 			a := New(nil, nil, nil)
// 			handler := a.Middleware(successHandler)
// 			r := httptest.NewRequest("GET", "http://example.com/foo", nil)
// 			w := httptest.NewRecorder()

// 			handler.ServeHTTP(w, r)

// 			So(w.Code, ShouldEqual, http.StatusForbidden)
// 			So(w.Body.String(), ShouldContainSubstring, "forbidden")
// 		})
// 		Convey("returns unauthorised status if the auth token  is not provided", func() {
// 			a := New(nil, nil, nil)
// 			handler := a.Middleware(successHandler)
// 			r := httptest.NewRequest("GET", "http://example.com/v3/charts/abcde", nil)
// 			w := httptest.NewRecorder()

// 			handler.ServeHTTP(w, r)

// 			So(w.Code, ShouldEqual, http.StatusUnauthorized)
// 			So(w.Body.String(), ShouldContainSubstring, "no authorisation provided")
// 		})
// 		Convey("returns unauthorised status if chart access is not allowed", func() {
// 			a := New(pcDenyMock, tpMock, csMock)
// 			handler := a.Middleware(successHandler)
// 			r := httptest.NewRequest("GET", "http://example.com/v3/charts/abcde", nil)
// 			r.Header.Add("Authorization", "Bearer abc")
// 			w := httptest.NewRecorder()

// 			handler.ServeHTTP(w, r)

// 			So(w.Code, ShouldEqual, http.StatusUnauthorized)
// 			So(w.Body.String(), ShouldContainSubstring, "unauthorised")
// 		})
// 		Convey("forwards to the next handler if chart access is allowed", func() {
// 			a := New(pcAllowMock, tpMock, csMock)
// 			handler := a.Middleware(successHandler)
// 			r := httptest.NewRequest("GET", "http://example.com/v3/charts/abcde", nil)
// 			r.Header.Add("Authorization", "Bearer abc")
// 			w := httptest.NewRecorder()

// 			handler.ServeHTTP(w, r)

// 			So(w.Code, ShouldEqual, http.StatusOK)
// 			So(w.Body.String(), ShouldContainSubstring, "success")
// 		})
// 	})
// }

// func TestChartAccessAllowed(t *testing.T) {
// 	Convey("chartAccessAllowed", t, func() {
// 		Convey("denies access if incorrect token provided", func() {
// 			tpErrorMock := &mocks.TokenParserMock{
// 				ParseFunc: func(tokenString string) (*permissions.EntityData, error) { return nil, errors.New("parse error") },
// 			}
// 			a := New(pcAllowMock, tpErrorMock, csMock)

// 			allowed := a.chartAccessAllowed(context.Background(), "token01", "chart01")

// 			So(allowed, ShouldBeFalse)
// 			So(tpErrorMock.ParseCalls(), ShouldHaveLength, 1)
// 			So(tpErrorMock.ParseCalls()[0].TokenString, ShouldEqual, "token01")
// 		})
// 		Convey("denies access on collection lookup error", func() {
// 			csErrorMock := &mocks.ChartStoreMock{
// 				GetCollectionIDFunc: func(chartID string) (string, error) { return "", errors.New("parse error") },
// 			}
// 			a := New(pcAllowMock, tpMock, csErrorMock)

// 			allowed := a.chartAccessAllowed(context.Background(), "token01", "chart01")

// 			So(allowed, ShouldBeFalse)
// 			So(csErrorMock.GetCollectionIDCalls(), ShouldHaveLength, 1)
// 			So(csErrorMock.GetCollectionIDCalls()[0].ChartID, ShouldEqual, "chart01")
// 		})
// 		Convey("denies access on permission lookup error", func() {
// 			pcErrorMock := &mocks.PermissionsCheckerMock{
// 				HasPermissionFunc: func(ctx context.Context, entityData permissions.EntityData, permission string, attributes map[string]string) (bool, error) {
// 					return false, errors.New("lookup error")
// 				},
// 			}
// 			a := New(pcErrorMock, tpMock, csMock)

// 			allowed := a.chartAccessAllowed(context.Background(), "token01", "chart01")

// 			So(allowed, ShouldBeFalse)
// 			So(pcErrorMock.HasPermissionCalls(), ShouldHaveLength, 1)
// 		})
// 		Convey("allows access on permission", func() {
// 			a := New(pcAllowMock, tpMock, csMock)

// 			allowed := a.chartAccessAllowed(context.Background(), "token01", "chart01")

// 			So(allowed, ShouldBeTrue)
// 		})
// 	})
// }

// func TestGetToken(t *testing.T) {
// 	Convey("GetToken", t, func() {
// 		Convey("returns empty string if no auth header is provided", func() {
// 			a := New(nil, nil, nil)
// 			r := httptest.NewRequest("GET", "http://example.com/v3/charts/abcde", nil)

// 			token := a.getToken(r)

// 			So(token, ShouldEqual, "")

// 		})
// 		Convey("returns token if provided in the auth header", func() {
// 			a := New(nil, nil, nil)
// 			r := httptest.NewRequest("GET", "http://example.com/v3/charts/abcde", nil)
// 			r.Header.Add("Authorization", "abc")

// 			token := a.getToken(r)

// 			So(token, ShouldEqual, "abc")
// 		})
// 		Convey("strips bearer prefix if provided in the auth header", func() {
// 			a := New(nil, nil, nil)
// 			r := httptest.NewRequest("GET", "http://example.com/v3/charts/abcde", nil)
// 			r.Header.Add("Authorization", "Bearer def")

// 			token := a.getToken(r)

// 			So(token, ShouldEqual, "def")
// 		})
// 	})
// }
// func TestExtractChartID(t *testing.T) {
// 	Convey("Decode parses the URL value from a string", t, func() {
// 		tests := []struct {
// 			description string
// 			url         string
// 			expected    string
// 		}{
// 			{
// 				description: "API get chart",
// 				url:         "https://api.datawrapper.de/v3/charts/abc01",
// 				expected:    "abc01",
// 			},
// 			{
// 				description: "API export chart",
// 				url:         "https://api.datawrapper.de/v3/charts/ABCDE/export/pdf",
// 				expected:    "ABCDE",
// 			},
// 			{
// 				description: "API get chart - too long ID",
// 				url:         "https://api.datawrapper.de/v3/charts/abcde01",
// 				expected:    "",
// 			},
// 			{
// 				description: "API export chart - too long ID",
// 				url:         "https://api.datawrapper.de/v3/charts/ABCDEFGH/export/pdf",
// 				expected:    "",
// 			},
// 			{
// 				description: "UI chart preview",
// 				url:         "https://app.datawrapper.de/preview/abc01",
// 				expected:    "abc01",
// 			},
// 			{
// 				description: "UI chart preview - too long ID",
// 				url:         "https://app.datawrapper.de/preview/abc012345",
// 				expected:    "",
// 			},
// 			{
// 				description: "UI chart preview - extra path segment",
// 				url:         "https://app.datawrapper.de/preview/abc01/pdf",
// 				expected:    "",
// 			},
// 			{
// 				description: "no match",
// 				url:         "https://app.datawrapper.de/abc01",
// 				expected:    "",
// 			},
// 		}
// 		for _, tt := range tests {
// 			Convey(tt.description, func() {
// 				url, err := url.Parse(tt.url)
// 				So(err, ShouldBeNil)
// 				So(extractChartID(url), ShouldEqual, tt.expected)
// 			})
// 		}
// 	})
// }
