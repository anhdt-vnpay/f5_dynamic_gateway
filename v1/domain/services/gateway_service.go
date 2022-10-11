package services

import (
	"context"
	"net/http"
)

type GatewayService interface {
	Handle(ctx context.Context, pattern string) (*http.ServeMux, error)
}
