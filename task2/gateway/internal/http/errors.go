package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

type errorResponse struct {
	Message string `json:"message"`
}

func requestErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func responseErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	statusCode, message := mapErrorToHTTP(err)
	writeJSONError(w, statusCode, message)
}

func mapErrorToHTTP(err error) (int, string) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout, "request timed out"
	case errors.Is(err, context.Canceled):
		return http.StatusRequestTimeout, "request canceled"
	}

	st, ok := grpcStatus.FromError(err)
	if !ok {
		return http.StatusInternalServerError, "internal server error"
	}

	message := st.Message()
	if message == "" {
		message = "internal server error"
	}

	switch st.Code() {
	case grpcCodes.InvalidArgument:
		return http.StatusBadRequest, message
	case grpcCodes.Unauthenticated:
		return http.StatusUnauthorized, message
	case grpcCodes.PermissionDenied:
		return http.StatusForbidden, message
	case grpcCodes.NotFound:
		return http.StatusNotFound, message
	case grpcCodes.ResourceExhausted:
		return http.StatusTooManyRequests, message
	case grpcCodes.DeadlineExceeded:
		return http.StatusGatewayTimeout, message
	case grpcCodes.Unavailable:
		return http.StatusBadGateway, message
	default:
		return http.StatusInternalServerError, message
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Message: message})
}