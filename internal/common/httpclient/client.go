package httpclient

import (
	"eticket-api/config"
	"net/http"
)

type HTTP struct {
	*http.Client
}

func NewHTTPClient(cfg *config.Config) *HTTP {
	return &HTTP{
		Client: &http.Client{},
	}
}
