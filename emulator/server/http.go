package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	LivenessPath = "/live"
	MetricsPath  = "/metrics"
)

type HTTPHeader struct {
	Key   string
	Value string
}

type HTTPServer struct {
	httpServer *http.Server
}

func NewHTTPServer(
	grpcServer *GRPCServer,
	liveness *LivenessTicker,
	port int,
	headers []HTTPHeader,
) *HTTPServer {
	wrappedServer := grpcweb.WrapServer(
		grpcServer.Server(),
		// TODO: is this needed?
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	mux := http.NewServeMux()

	// register metrics handler
	mux.Handle(MetricsPath, promhttp.Handler())

	// register liveness handler
	mux.Handle(LivenessPath, liveness.Handler())

	// register gRPC HTTP proxy
	mux.Handle("/", wrappedHandler(wrappedServer, headers))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &HTTPServer{
		httpServer: httpServer,
	}
}

func (h *HTTPServer) Start() error {
	err := h.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (h *HTTPServer) Stop() {
	_ = h.httpServer.Shutdown(context.Background())
}

func wrappedHandler(wrappedServer *grpcweb.WrappedGrpcServer, headers []HTTPHeader) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		setResponseHeaders(&res, headers)

		if (*req).Method == "OPTIONS" {
			return
		}

		wrappedServer.ServeHTTP(res, req)
	}
}

func setResponseHeaders(w *http.ResponseWriter, headers []HTTPHeader) {
	for _, header := range headers {
		(*w).Header().Set(header.Key, header.Value)
	}
}
