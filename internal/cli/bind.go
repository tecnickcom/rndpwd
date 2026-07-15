package cli

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/tecnickcom/nurago/pkg/bootstrap"
	"github.com/tecnickcom/nurago/pkg/healthcheck"
	"github.com/tecnickcom/nurago/pkg/httpclient"
	"github.com/tecnickcom/nurago/pkg/httpserver"
	"github.com/tecnickcom/nurago/pkg/httputil"
	"github.com/tecnickcom/nurago/pkg/httputil/jsendx"
	"github.com/tecnickcom/nurago/pkg/ipify"
	"github.com/tecnickcom/nurago/pkg/metrics"
	"github.com/tecnickcom/nurago/pkg/redact"
	"github.com/tecnickcom/nurago/pkg/traceid"
	"github.com/tecnickcom/rndpwd/internal/httphandler"
	instr "github.com/tecnickcom/rndpwd/internal/metrics"
	"github.com/tecnickcom/rndpwd/internal/password"
	"github.com/tecnickcom/rndpwd/internal/validator"
)

// bind is the entry point of the service, this is where the wiring of all
// components happens. The wiring is split into small, named helpers
// (newIpifyClient, bindServiceHandlers) so each concern stays readable and
// independently testable.
func bind(cfg *appConfig, appInfo *jsendx.AppInfo, mtr instr.Metrics, wg *sync.WaitGroup, sc chan struct{}) bootstrap.BindFunc {
	return func(ctx context.Context, l *slog.Logger, m metrics.Client) error {
		jsx := jsendx.NewJSXResp(httputil.NewHTTPResp(l))

		// Sanitizes the HTTP dumps, query strings, and error URLs before they are
		// logged. One instance is shared by every client and server below.
		logRedactor := newLogRedactor()

		// Common outbound HTTP client options shared by every external client
		// (structured logging, trace propagation, and metrics instrumentation).
		// This base is built once and reused: each client constructor appends its
		// own timeout on top of it (see newIpifyClient).
		httpClientOpts := []httpclient.Option{
			httpclient.WithLogger(l),
			httpclient.WithRoundTripper(m.InstrumentRoundTripper),
			httpclient.WithTraceIDHeaderName(traceid.DefaultHeader),
			httpclient.WithComponent(appInfo.ProgramName),
			httpclient.WithRedactFn(logRedactor.BytesToString),
		}

		// ipify is used only as a diagnostic (the monitoring /ip route); it is
		// intentionally not part of the health checks.
		ipifyClient, err := newIpifyClient(cfg, httpClientOpts)
		if err != nil {
			return err
		}

		serviceBinder, statusHandler := bindServiceHandlers(cfg, appInfo, jsx, l, mtr)

		middleware := func(args httpserver.MiddlewareArgs, next http.Handler) http.Handler {
			return m.InstrumentHandler(args.Path, next.ServeHTTP)
		}

		// start monitoring server
		httpMonitoringOpts := []httpserver.Option{
			httpserver.WithLogger(l),
			httpserver.WithServerAddr(cfg.Servers.Monitoring.Address),
			httpserver.WithRequestTimeout(time.Duration(cfg.Servers.Monitoring.Timeout) * time.Second),
			httpserver.WithMetricsHandlerFunc(m.MetricsHandlerFunc()),
			httpserver.WithTraceIDHeaderName(traceid.DefaultHeader),
			httpserver.WithMiddlewareFn(middleware),
			httpserver.WithNotFoundHandlerFunc(jsx.DefaultNotFoundHandlerFunc(appInfo)),
			httpserver.WithMethodNotAllowedHandlerFunc(jsx.DefaultMethodNotAllowedHandlerFunc(appInfo)),
			httpserver.WithPanicHandlerFunc(jsx.DefaultPanicHandlerFunc(appInfo)),
			httpserver.WithEnableAllDefaultRoutes(),
			httpserver.WithIndexHandlerFunc(jsx.DefaultIndexHandler(appInfo)),
			httpserver.WithIPHandlerFunc(jsx.DefaultIPHandler(appInfo, ipifyClient.GetPublicIP)),
			httpserver.WithPingHandlerFunc(jsx.DefaultPingHandler(appInfo)),
			httpserver.WithStatusHandlerFunc(statusHandler),
			httpserver.WithRedactFn(logRedactor.BytesToString),
			httpserver.WithShutdownWaitGroup(wg),
			httpserver.WithShutdownSignalChan(sc),
		}

		httpMonitoringServer, err := httpserver.New(ctx, httpserver.NopBinder(), httpMonitoringOpts...)
		if err != nil {
			return fmt.Errorf("error creating monitoring HTTP server: %w", err)
		}

		// example of custom metric
		mtr.IncExampleCounter("START")

		// start public server
		httpPublicOpts := []httpserver.Option{
			httpserver.WithLogger(l),
			httpserver.WithServerAddr(cfg.Servers.Public.Address),
			httpserver.WithRequestTimeout(time.Duration(cfg.Servers.Public.Timeout) * time.Second),
			httpserver.WithMiddlewareFn(middleware),
			httpserver.WithTraceIDHeaderName(traceid.DefaultHeader),
			httpserver.WithEnableDefaultRoutes(httpserver.PingRoute),
			httpserver.WithRedactFn(logRedactor.BytesToString),
			httpserver.WithShutdownWaitGroup(wg),
			httpserver.WithShutdownSignalChan(sc),
		}

		httpPublicServer, err := httpserver.New(ctx, serviceBinder, httpPublicOpts...)
		if err != nil {
			return fmt.Errorf("error creating public HTTP server: %w", err)
		}

		httpMonitoringServer.StartServer()
		httpPublicServer.StartServer()

		return nil
	}
}

// newLogRedactor builds the redactor applied to the HTTP request and response
// dumps, query strings, and error URLs written to the logs.
//
// Every rule class is disabled here, so the redactor passes its input through
// unchanged: this service's logs are NOT sanitized. This departs from the
// library default: httpclient, httpserver, and httpreverseproxy redact with the
// package-level default redactor ([redact.HTTPDataString], all rules on) when the
// WithRedactFn option is omitted. Anything sensitive that reaches a logged dump,
// query string, or error URL is therefore logged verbatim.
//
// Re-enable a rule class by removing it from the list below. The other knobs are
// [redact.WithMarker] (replace the `***` placeholder), [redact.WithExtraTokens]
// and [redact.WithoutTokens] (adjust the sensitive key names), and
// [redact.WithLuhnCheck] (gate card detection on the Luhn checksum). Deleting
// this function and the WithRedactFn options that reference it restores the
// fully redacting default.
//
// A [redact.Redactor] is immutable after construction and safe for concurrent
// use, so a single instance is shared by every client and server.
func newLogRedactor() *redact.Redactor {
	return redact.New(
		redact.WithoutRules(
			redact.RuleHeaders,
			redact.RuleJSON,
			redact.RuleURLEncoded,
			redact.RuleXML,
			redact.RuleUserinfo,
			redact.RuleJWT,
			redact.RuleVendorTokens,
			redact.RulePEM,
			redact.RuleCards,
		),
	)
}

// newIpifyClient builds the ipify client used by the monitoring server's /ip
// diagnostic route.
//
// It takes the shared base HTTP client options and appends the ipify-specific
// timeout, keeping every outbound client consistent (logging, tracing, metrics)
// while each keeps its own timeout. The base slice is left untouched: it has no
// spare capacity, so the append always copies rather than mutating the caller's
// slice, which keeps it safe to reuse for the next client.
func newIpifyClient(cfg *appConfig, baseHTTPClientOpts []httpclient.Option) (*ipify.Client, error) {
	ipifyTimeout := time.Duration(cfg.Clients.Ipify.Timeout) * time.Second

	ipifyHTTPClient := httpclient.New(append(
		baseHTTPClientOpts,
		httpclient.WithTimeout(ipifyTimeout),
	)...)

	ipifyClient, err := ipify.New(
		ipify.WithHTTPClient(ipifyHTTPClient),
		ipify.WithTimeout(ipifyTimeout),
		ipify.WithURL(cfg.Clients.Ipify.Address),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build ipify client: %w", err)
	}

	return ipifyClient, nil
}

// bindServiceHandlers wires the service binder together with the status handler.
//
// When the service is disabled it returns a no-op binder and the default status
// handler. When enabled it attaches the real password-generator handler and
// upgrades the status handler to a health check.
func bindServiceHandlers(
	cfg *appConfig,
	appInfo *jsendx.AppInfo,
	jsx *jsendx.JSXResp,
	l *slog.Logger,
	mtr instr.Metrics,
) (httpserver.Binder, http.HandlerFunc) {
	if !cfg.Enabled {
		return httpserver.NopBinder(), jsx.DefaultStatusHandler(appInfo)
	}

	// The validation options are static and already proven valid, so New cannot
	// fail here; the error is intentionally discarded.
	val, _ := validator.New("json")

	serviceBinder := httphandler.New(
		l,
		appInfo,
		mtr,
		val,
		password.New(
			cfg.Random.Charset,
			cfg.Random.Length,
			cfg.Random.Quantity,
		),
	)

	// override the default status handler with a health check
	healthCheckHandler := healthcheck.NewHandler(
		[]healthcheck.HealthCheck{},
		healthcheck.WithLogger(l),
		healthcheck.WithResultWriter(jsx.HealthCheckResultWriter(appInfo)),
	)

	return serviceBinder, healthCheckHandler.ServeHTTP
}
