package httpx

type ErrorResponse struct {
	Error string `json:"error" example:"access denied"`
}