package user

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "github.com/Hubcher/project-management/contracts/gen/proto/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService interface {
	CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error)
	GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error)
	ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error)
	UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error)
	DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*emptypb.Empty, error)
}

func NewCreateUserHandler(log *slog.Logger, service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := service.CreateUser(r.Context(), toCreateUserPB(req))
		if err != nil {
			writeGRPCError(log, w, err, "cannot create user")
			return
		}

		writeJSON(w, http.StatusCreated, toUserResponse(user))
	}
}

func NewGetUserByIdHandler(log *slog.Logger, service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := extractUserID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		user, err := service.GetUser(r.Context(), &userpb.GetUserRequest{Id: id})
		if err != nil {
			writeGRPCError(log, w, err, "cannot get user")
			return
		}

		writeJSON(w, http.StatusOK, toUserResponse(user))
	}
}

func NewListUsersHandler(log *slog.Logger, service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.URL.Query().Get("role")

		resp, err := service.ListUsers(r.Context(), &userpb.ListUsersRequest{Role: role})
		if err != nil {
			writeGRPCError(log, w, err, "cannot list users")
			return
		}

		writeJSON(w, http.StatusOK, toListUsersResponse(resp.GetUsers()))
	}
}

func NewUpdateUserHandler(log *slog.Logger, service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := extractUserID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		var req UpdateUserRequest
		if err = decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := service.UpdateUser(r.Context(), toUpdateUserPB(id, req))
		if err != nil {
			writeGRPCError(log, w, err, "cannot update user")
			return
		}
		writeJSON(w, http.StatusOK, toUserResponse(user))
	}
}

func NewDeleteUserHandler(log *slog.Logger, service UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := extractUserID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		_, err = service.DeleteUser(r.Context(), &userpb.DeleteUserRequest{Id: id})
		if err != nil {
			writeGRPCError(log, w, err, "cannot delete user")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func decodeJSON(r *http.Request, dst any) error {
	defer func() {
		if err := r.Body.Close(); err != nil {
			return
		}
	}()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body id empty")
		}
		return err
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain only one JSON object")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, statuCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuCode)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}

func writeGRPCError(log *slog.Logger, w http.ResponseWriter, err error, msg string) {
	st, ok := status.FromError(err)
	if !ok {
		log.Error(msg, "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	log.Error(msg, "code", st.Code(), "message", st.Message())

	switch st.Code() {
	case codes.InvalidArgument:
		writeError(w, http.StatusBadRequest, st.Message())
	case codes.NotFound:
		writeError(w, http.StatusNotFound, st.Message())
	case codes.AlreadyExists:
		writeError(w, http.StatusConflict, st.Message())
	default:
		writeError(w, http.StatusInternalServerError, st.Message())
	}
}

func extractUserID(r *http.Request) (string, error) {
	id := r.PathValue("id")
	if id == "" {
		return "", errors.New("empty id")
	}
	return id, nil
}
