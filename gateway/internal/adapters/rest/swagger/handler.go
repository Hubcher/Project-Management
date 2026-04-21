package swagger

import (
	"net/http"
	"os"
	"path/filepath"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewUIHandler() http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL("/openapi/swagger.json"),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DefaultModelsExpandDepth(-1),
	)
}

func NewSpecHandler() http.Handler {
	return http.StripPrefix("/openapi", http.FileServer(http.Dir(resolveOpenAPIDir())))
}

func resolveOpenAPIDir() string {
	candidates := []string{
		"api",
		filepath.Join("gateway", "api"),
		filepath.Join(string(filepath.Separator), "app", "api"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(filepath.Join(candidate, "swagger.json")); err == nil {
			return candidate
		}
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	return "api"
}
