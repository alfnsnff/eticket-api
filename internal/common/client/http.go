package client

import (
	"eticket-api/config"
	"net/http"
)

func NewHttp(cfg *config.Config) *http.Client {
	return &http.Client{}
}
