package module

import (
	"eticket-api/config"
	"eticket-api/internal/client"
	"net/http"
)

type ClientModule struct {
	TripayClient *client.TripayClient
}

func NewClientModule(http_client *http.Client, config *config.Configuration) *ClientModule {
	return &ClientModule{
		TripayClient: client.NewTripayClient(http_client, &config.Tripay),
	}
}
