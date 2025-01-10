package domain

import "time"

type Ticket struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}


type TicketRepository interface {
    Create(ticket *Ticket) error
    // GetByID(id int) (*Ticket, error)
    GetAll() ([]*Ticket, error)
    // Update(ticket *Ticket) error
    // Delete(id int) error
}
