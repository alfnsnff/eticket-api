package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// ShipDTO represents a Ship.
type PriceRead struct {
	ID          uint      `json:"id"`
	RouteID     uint      `json:"route_id"`
	ShipClassID uint      `json:"ship_class_id"`
	Price       float32   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ShipDTO represents a Ship.
type PriceCreate struct {
	RouteID     uint    `json:"route_id"`
	ShipClassID uint    `json:"ship_class_id"`
	Price       float32 `json:"price"`
}

type PriceRequest struct {
	RouteID     uint `json:"route_id"`
	ShipClassID uint `json:"ship_class_id"`
}

func ToPriceDTO(price *entities.Price) PriceRead {
	var priceResponse PriceRead
	copier.Copy(&priceResponse, &price) // Automatically maps matching fields
	return priceResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToPriceDTOs(prices []*entities.Price) []PriceRead {
	var priceResponses []PriceRead
	for _, price := range prices {
		priceResponses = append(priceResponses, ToPriceDTO(price))
	}
	return priceResponses
}

func ToPriceEntity(priceCreate *PriceCreate) entities.Price {
	var price entities.Price
	copier.Copy(&price, &priceCreate)
	return price
}
