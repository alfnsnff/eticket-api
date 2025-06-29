package httpclient

import (
	"eticket-api/config"
	"net/http"
	"time"
)

type HTTP struct {
	*http.Client
}

func NewHTTPClient(cfg *config.Config) *HTTP {
	return &HTTP{
		Client: &http.Client{
			Timeout: 45 * time.Second, // contoh timeout 10 detik
		},
	}
}
