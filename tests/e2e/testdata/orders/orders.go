package orders

import (
	"time"
)

type Order interface {
	isOrder()
}

type PriorityOrder struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Priority  int       `json:"priority"`
}

func (a PriorityOrder) isOrder() {}

type RegularOrder struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func (a RegularOrder) isOrder() {}
