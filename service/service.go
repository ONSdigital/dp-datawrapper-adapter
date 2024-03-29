package service

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-authorisation/v2/permissions"
	"github.com/ONSdigital/dp-datawrapper-adapter/config"
	"github.com/ONSdigital/dp-datawrapper-adapter/datawrapper"
	"github.com/ONSdigital/dp-datawrapper-adapter/proxy"
	"github.com/ONSdigital/dp-datawrapper-adapter/routes"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

// Service contains the healthcheck, server and serviceList for the controller
type Service struct {
	Config      *config.Config
	HealthCheck HealthChecker
	Server      HTTPServer
	ServiceList *ExternalServiceList
}

// New creates a new service
func New() *Service {
	return &Service{}
}

// Init initialises all the service dependencies, including healthcheck with checkers, api and middleware
func (svc *Service) Init(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList) (err error) {
	log.Info(ctx, "initialising service")

	svc.Config = cfg
	svc.ServiceList = serviceList

	// Initialise clients
	// parser, _ := jwt.NewCognitoRSAParser(cfg.JWTVerificationPublicKeys)
	// if err != nil {
	// 	log.Fatal(ctx, "failed to create cognito parser", err)
	// 	return err
	// }
	apiProxy, err := proxy.New("/api", cfg.DatawrapperAPIURL)
	if err != nil {
		log.Fatal(ctx, "failed to create api proxy", err)
		return err
	}
	uiProxy, err := proxy.New("", cfg.DatawrapperUIURL)
	if err != nil {
		log.Fatal(ctx, "failed to create ui proxy", err)
		return err
	}
	permissionsChecker := permissions.NewChecker(ctx, cfg.PermissionsAPIHost, cfg.PermissionsCacheUpdateInterval, cfg.PermissionsMaxCacheTime)
	clients := routes.Clients{
		Datawrapper:        datawrapper.NewClient(cfg.DatawrapperAPIURL, cfg.DatawrapperAPIToken),
		APIProxy:           apiProxy,
		UIProxy:            uiProxy,
		PermissionsChecker: permissionsChecker,
		// Authoriser: authoriser.New(
		// 	permissionsChecker,
		// 	parser,
		// 	&charts.Stub{},
		// ),
	}

	// Get healthcheck with checkers
	svc.HealthCheck, err = serviceList.GetHealthCheck(cfg, BuildTime, GitCommit, Version)
	if err != nil {
		log.Fatal(ctx, "failed to create health check", err)
		return err
	}
	if err = svc.registerCheckers(ctx, clients); err != nil {
		log.Error(ctx, "failed to register checkers", err)
		return err
	}
	clients.HealthCheckHandler = svc.HealthCheck.Handler

	// Initialise router
	r := mux.NewRouter()
	routes.Setup(ctx, r, cfg, clients)
	svc.Server = serviceList.GetHTTPServer(cfg.BindAddr, r)

	return nil
}

// Run starts an initialised service
func (svc *Service) Run(ctx context.Context, svcErrors chan error) {
	log.Info(ctx, "Starting service", log.Data{"config": svc.Config})

	// Start healthcheck
	svc.HealthCheck.Start(ctx)

	// Start HTTP server
	log.Info(ctx, "Starting server")
	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			log.Fatal(ctx, "failed to start http listen and serve", err)
			svcErrors <- err
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	log.Info(ctx, "commencing graceful shutdown")
	ctx, cancel := context.WithTimeout(ctx, svc.Config.GracefulShutdownTimeout)
	hasShutdownError := false

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		log.Info(ctx, "stop health checkers")
		svc.HealthCheck.Stop()

		// TODO: close any backing services here, e.g. client connections to databases

		// stop any incoming requests
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (svc *Service) registerCheckers(ctx context.Context, c routes.Clients) (err error) {
	hasErrors := false

	if err = svc.HealthCheck.AddCheck("datawrapper", c.Datawrapper.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "failed to add datawrapper health checker", err)
	}

	if err = svc.HealthCheck.AddCheck("permission checker", c.PermissionsChecker.HealthCheck); err != nil {
		hasErrors = true
		log.Error(ctx, "failed to add permission healthchecker", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}

	return nil
}
