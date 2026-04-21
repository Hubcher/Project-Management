package httpx

import (
    "encoding/json"
    "errors"
    "io"
    "net/http"

    "github.com/Hubcher/project-management/gateway/internal/core"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func DecodeJSON(r *http.Request, dst any) error {
    defer func() {
        _ = r.Body.Close()
    }()

    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    if err := dec.Decode(dst); err != nil {
        if errors.Is(err, io.EOF) {
            return errors.New("request body is empty")
        }
        return err
    }

    if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
        return errors.New("request body must contain only one JSON object")
    }
    return nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    _ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
    WriteJSON(w, statusCode, map[string]string{"error": message})
}

func WriteAnyError(w http.ResponseWriter, err error) {
    var statusErr *core.StatusError
    if errors.As(err, &statusErr) {
        WriteError(w, statusErr.Code, statusErr.Message)
        return
    }

    if st, ok := status.FromError(err); ok {
        switch st.Code() {
        case codes.InvalidArgument:
            WriteError(w, http.StatusBadRequest, st.Message())
        case codes.NotFound:
            WriteError(w, http.StatusNotFound, st.Message())
        case codes.AlreadyExists:
            WriteError(w, http.StatusConflict, st.Message())
        case codes.Unauthenticated:
            WriteError(w, http.StatusUnauthorized, st.Message())
        case codes.PermissionDenied:
            WriteError(w, http.StatusForbidden, st.Message())
        default:
            WriteError(w, http.StatusBadGateway, st.Message())
        }
        return
    }

    WriteError(w, http.StatusInternalServerError, "internal server error")
}
