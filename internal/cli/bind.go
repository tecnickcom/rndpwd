package cli

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Vonage/gosrvlib/pkg/bootstrap"
	"github.com/Vonage/gosrvlib/pkg/healthcheck"
	"github.com/Vonage/gosrvlib/pkg/httpclient"
	"github.com/Vonage/gosrvlib/pkg/httpserver"
	"github.com/Vonage/gosrvlib/pkg/httputil/jsendx"
	"github.com/Vonage/gosrvlib/pkg/ipify"
	"github.com/Vonage/gosrvlib/pkg/metrics"
	"github.com/Vonage/gosrvlib/pkg/traceid"
	"github.com/tecnickcom/rndpwd/internal/httphandler"
	instr "github.com/tecnickcom/rndpwd/internal/metrics"
	"github.com/tecnickcom/rndpwd/internal/password"
	"github.com/tecnickcom/rndpwd/internal/validator"
	"go.uber.org/zap"
)

// bind is the entry point of the service, this is where the wiring of all components happens.
//
//nolint:funlen
func bind(cfg *appConfig, appInfo *jsendx.AppInfo, mtr instr.Metrics, wg *sync.WaitGroup, sc chan struct{}) bootstrap.BindFunc {
	return func(ctx context.Context, _ *zap.Logger, m metrics.Client) error {
		// We assume the service is disabled and override the service binder if required
		serviceBinder := httpserver.NopBinder()
		statusHandler := jsendx.DefaultStatusHandler(appInfo)

		// common HTTP client options used for all outbound requests
		httpClientOpts := []httpclient.Option{
			httpclient.WithRoundTripper(m.InstrumentRoundTripper),
			httpclient.WithTraceIDHeaderName(traceid.DefaultHeader),
			httpclient.WithComponent(appInfo.ProgramName),
		}

		ipifyTimeout := time.Duration(cfg.Clients.Ipify.Timeout) * time.Second
		ipifyHTTPClient := httpclient.New(
			append(httpClientOpts, httpclient.WithTimeout(ipifyTimeout))...,
		)

		ipifyClient, err := ipify.New(
			ipify.WithHTTPClient(ipifyHTTPClient),
			ipify.WithTimeout(ipifyTimeout),
			ipify.WithURL(cfg.Clients.Ipify.Address),
		)
		if err != nil {
			return fmt.Errorf("failed to build ipify client: %w", err)
		}

		if cfg.Enabled {
			val, _ := validator.New("json")

			// wire the binder
			serviceBinder = httphandler.New(
				appInfo,
				mtr,
				val,
				password.New(
					cfg.Random.Charset,
					cfg.Random.Length,
					cfg.Random.Quantity,
				),
			)

			// override the default healthcheck handler
			healthCheckHandler := healthcheck.NewHandler(
				[]healthcheck.HealthCheck{},
				healthcheck.WithResultWriter(jsendx.HealthCheckResultWriter(appInfo)),
			)
			statusHandler = healthCheckHandler.ServeHTTP
		}

		middleware := func(args httpserver.MiddlewareArgs, next http.Handler) http.Handler {
			return m.InstrumentHandler(args.Path, next.ServeHTTP)
		}

		// start monitoring server
		httpMonitoringOpts := []httpserver.Option{
			httpserver.WithServerAddr(cfg.Servers.Monitoring.Address),
			httpserver.WithRequestTimeout(time.Duration(cfg.Servers.Monitoring.Timeout) * time.Second),
			httpserver.WithMetricsHandlerFunc(m.MetricsHandlerFunc()),
			httpserver.WithTraceIDHeaderName(traceid.DefaultHeader),
			httpserver.WithMiddlewareFn(middleware),
			httpserver.WithNotFoundHandlerFunc(jsendx.DefaultNotFoundHandlerFunc(appInfo)),
			httpserver.WithMethodNotAllowedHandlerFunc(jsendx.DefaultMethodNotAllowedHandlerFunc(appInfo)),
			httpserver.WithPanicHandlerFunc(jsendx.DefaultPanicHandlerFunc(appInfo)),
			httpserver.WithEnableAllDefaultRoutes(),
			httpserver.WithIndexHandlerFunc(jsendx.DefaultIndexHandler(appInfo)),
			httpserver.WithIPHandlerFunc(jsendx.DefaultIPHandler(appInfo, ipifyClient.GetPublicIP)),
			httpserver.WithPingHandlerFunc(jsendx.DefaultPingHandler(appInfo)),
			httpserver.WithStatusHandlerFunc(statusHandler),
			httpserver.WithShutdownWaitGroup(wg),
			httpserver.WithShutdownSignalChan(sc),
		}

		httpMonitoringServer, err := httpserver.New(ctx, httpserver.NopBinder(), httpMonitoringOpts...)
		if err != nil {
			return fmt.Errorf("error starting monitoring HTTP server: %w", err)
		}

		// example of custom metric
		mtr.IncExampleCounter("START")

		// start public server
		httpPublicOpts := []httpserver.Option{
			httpserver.WithServerAddr(cfg.Servers.Public.Address),
			httpserver.WithRequestTimeout(time.Duration(cfg.Servers.Public.Timeout) * time.Second),
			httpserver.WithMiddlewareFn(middleware),
			httpserver.WithTraceIDHeaderName(traceid.DefaultHeader),
			httpserver.WithEnableDefaultRoutes(httpserver.PingRoute),
			httpserver.WithShutdownWaitGroup(wg),
			httpserver.WithShutdownSignalChan(sc),
		}

		httpPublicServer, err := httpserver.New(ctx, serviceBinder, httpPublicOpts...)
		if err != nil {
			return fmt.Errorf("error starting public HTTP server: %w", err)
		}

		httpMonitoringServer.StartServer()
		httpPublicServer.StartServer()

		return nil
	}
}
