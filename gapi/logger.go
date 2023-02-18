package gapi

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error()
	}

	logger.
		Str("protocol", "gRPC").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received grpc request")
	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (r *ResponseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}

// HTTPLogger HttpLogger Middleware responsible for logging http request and http response
func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recorder := &ResponseRecorder{w, http.StatusOK, []byte{}}
		next.ServeHTTP(recorder, r)
		duration := time.Since(startTime)

		logger := log.Info()
		if recorder.statusCode != http.StatusOK {
			logger = log.Error().Bytes("response_body", recorder.body)
		}

		logger.
			Str("protocol", "HTTP").
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status_code", recorder.statusCode).
			Str("status_text", http.StatusText(recorder.statusCode)).
			Dur("duration", duration).
			Msg("received http request")
	})
}
