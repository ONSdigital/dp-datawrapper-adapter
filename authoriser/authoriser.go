package authoriser

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/ONSdigital/dp-authorisation/v2/permissions"
	"github.com/ONSdigital/log.go/v2/log"
)

type Authoriser struct {
	permissionsChecker PermissionsChecker
	tokenParser        TokenParser
	chartStore         ChartStore
}

//go:generate moq -skip-ensure -out mocks/checker.go -pkg mocks . PermissionsChecker
type PermissionsChecker interface {
	HasPermission(ctx context.Context, entityData permissions.EntityData, permission string, attributes map[string]string) (bool, error)
}

//go:generate moq -skip-ensure -out mocks/parser.go -pkg mocks . TokenParser
type TokenParser interface {
	Parse(tokenString string) (*permissions.EntityData, error)
}

//go:generate moq -skip-ensure -out mocks/store.go -pkg mocks . ChartStore
type ChartStore interface {
	GetCollectionID(chartID string) (string, error)
}

func New(permissionsChecker PermissionsChecker, tokenParser TokenParser, chartStore ChartStore) *Authoriser {
	return &Authoriser{
		permissionsChecker: permissionsChecker,
		tokenParser:        tokenParser,
		chartStore:         chartStore,
	}
}

func extractChartID(url *url.URL) string {
	regex := regexp.MustCompile(`^/v3/charts/(?P<id>[a-zA-Z0-9]{5})($|/)`)
	match := regex.FindStringSubmatch(url.Path)
	if len(match) > 1 {
		return match[1]
	}
	regex = regexp.MustCompile(`^/preview/(?P<id>[a-zA-Z0-9]{5})$`)
	match = regex.FindStringSubmatch(url.Path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func (a *Authoriser) getToken(r *http.Request) string {
	authToken := r.Header.Get("Authorization")
	return strings.TrimPrefix(authToken, "Bearer ")
}

func (a *Authoriser) chartAccessAllowed(ctx context.Context, token string, chartID string) bool {
	entityData, err := a.tokenParser.Parse(token)
	if err != nil {
		log.Error(ctx, "chart access check: jwt parse error", err)
		return false
	}
	collectionID, err := a.chartStore.GetCollectionID(chartID)
	if err != nil {
		log.Error(ctx, "chart access check: collection lookup error", err)
		return false
	}

	permission := "legacy.read"
	attributes := map[string]string{"collection_id": collectionID}

	hasPermission, err := a.permissionsChecker.HasPermission(ctx, *entityData, permission, attributes)
	if err != nil {
		log.Error(ctx, "chart access check: permissions lookup error", err)
		return false
	}

	if !hasPermission {
		log.Info(ctx, "chart access check: request not permitted")
		return false
	}

	return true
}

func (a *Authoriser) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chartID := extractChartID(r.URL)
		if chartID == "" {
			log.Info(r.Context(), "authorisation failed: non chart related request")
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		token := a.getToken(r)
		if len(token) == 0 {
			log.Info(r.Context(), "authorisation failed: no authorisation header in request")
			http.Error(w, "no authorisation provided", http.StatusUnauthorized)
			return
		}
		if !a.chartAccessAllowed(r.Context(), token, chartID) {
			log.Info(r.Context(), "authorisation failed: chart access not allowed")
			http.Error(w, "unauthorised", http.StatusUnauthorized)
			return
		}
		log.Info(r.Context(), "authorisation success: chart access allowed")
		next.ServeHTTP(w, r)
	})
}

func (a *Authoriser) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return a.handler(next)
	}
}
