package tfeapi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/DataDog/jsonapi"
	"github.com/tofutf/tofutf/internal"
)

var codes = map[error]int{
	internal.ErrResourceNotFound:        http.StatusNotFound,
	internal.ErrAccessNotPermitted:      http.StatusForbidden,
	internal.ErrInvalidTerraformVersion: http.StatusUnprocessableEntity,
	internal.ErrResourceAlreadyExists:   http.StatusConflict,
	internal.ErrConflict:                http.StatusConflict,
}

func lookupHTTPCode(err error) int {
	for errType, code := range codes {
		if errors.Is(err, errType) {
			return code
		}
	}
	return http.StatusInternalServerError
}

// Error writes an HTTP response with a JSON-API encoded error.
func Error(w http.ResponseWriter, err error) {
	var (
		httpError *internal.HTTPError
		missing   *internal.MissingParameterError
		code      int
	)
	// If error is type internal.HTTPError then extract its status code
	if errors.As(err, &httpError) {
		code = httpError.Code
	} else if errors.As(err, &missing) {
		// report missing parameter errors as a 422
		code = http.StatusUnprocessableEntity
	} else {
		code = lookupHTTPCode(err)
	}
	b, marshalErr := jsonapi.Marshal(&jsonapi.Error{
		Status: &code,
		Title:  http.StatusText(code),
		Detail: err.Error(),
	})
	if marshalErr != nil {
		panic(marshalErr)
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf("failed to handle tfeapi request: %s", err.Error()), pcs[0])
	_ = slog.Default().Handler().Handle(context.Background(), r)
	w.Header().Set("Content-type", mediaType)
	w.WriteHeader(code)
	w.Write(b) //nolint:errcheck
}
