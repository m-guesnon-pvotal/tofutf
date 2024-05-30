package api

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

type LoggingRoundTripper struct {
	logger              *slog.Logger
	defaultRoundTripper http.RoundTripper
}

func (t LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqStr string
	var respStr string
	var statusCode int

	if req != nil && t.logger.Enabled(req.Context(), slog.LevelDebug) {
		defer func() {
			// Do work after the response is received
			t.logger.DebugContext(req.Context(), "Sent client request", "url", req.URL.String(), "method", req.Method, "status", statusCode, "payload", reqStr, "response", respStr)
		}()

		if req.Body != nil {
			buf, err := io.ReadAll(req.Body)
			if err != nil {
				t.logger.ErrorContext(req.Context(), "Failed reading http client request body", "err", err)
			} else {
				reqStr = string(buf)
				req.Body = io.NopCloser(bytes.NewBuffer(buf))
			}
		}
	}

	// Do work before the request is sent
	resp, err := t.defaultRoundTripper.RoundTrip(req)

	if req != nil && t.logger.Enabled(req.Context(), slog.LevelDebug) {
		if resp != nil {
			statusCode = resp.StatusCode
			if resp.Body != nil {
				buf, err := io.ReadAll(resp.Body)
				if err != nil {
					t.logger.ErrorContext(req.Context(), "Failed reading http client response body", "err", err)
				} else {
					resp.Body = io.NopCloser(bytes.NewBuffer(buf))
					respStr = string(buf)
				}
			}
		}
	}

	return resp, err
}
