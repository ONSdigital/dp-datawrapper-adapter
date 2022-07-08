package routes

import (
	"context"
	"net/http"
	"net/http/httputil"

	"github.com/ONSdigital/dp-authorisation/v2/jwt"
	"github.com/ONSdigital/dp-authorisation/v2/permissions"
	"github.com/ONSdigital/dp-datawrapper-adapter/authoriser"
	"github.com/ONSdigital/dp-datawrapper-adapter/charts"
	"github.com/ONSdigital/dp-datawrapper-adapter/config"
	"github.com/ONSdigital/dp-datawrapper-adapter/datawrapper"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	HealthCheckHandler func(w http.ResponseWriter, req *http.Request)
	Datawrapper        *datawrapper.Client
	APIProxy           *httputil.ReverseProxy
	UIProxy            *httputil.ReverseProxy
	PermissionsChecker *permissions.Checker
	TokenParser        *jwt.CognitoRSAParser
	ChartStore         *charts.MongoStore
}

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, c Clients) {
	authoriser := authoriser.New(c.PermissionsChecker, c.TokenParser, c.ChartStore)
	authoriserMiddleware := authoriser.Middleware()

	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(c.HealthCheckHandler)
	r.StrictSlash(true).PathPrefix("/api").Handler(authoriserMiddleware(c.APIProxy))
	r.StrictSlash(true).PathPrefix("").Handler(authoriserMiddleware(c.UIProxy))
}
