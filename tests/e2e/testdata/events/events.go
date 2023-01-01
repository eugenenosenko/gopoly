package events

import (
	"time"

	"github.com/eugenenosenko/gopoly/tests/e2e/testdata/orders"
	u "github.com/eugenenosenko/gopoly/tests/e2e/testdata/users"
)

type UserEvent interface {
	IsUserEvent()
}

type UserDeletedEvent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	User u.User `json:"user"`
}

func (e UserDeletedEvent) IsUserEvent() {}

type UserCreatedEvent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	User u.User `json:"user"`
}

func (e UserCreatedEvent) IsUserEvent() {}

type OrderEvent interface {
	IsOrderEvent()
}

type OrderCompletedEvent struct {
	ID    string       `json:"id"`
	Type  string       `json:"type"`
	Order orders.Order `json:"order"`
}

func (e OrderCompletedEvent) IsOrderEvent() {}

type OrderCancelledEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	CancelledOn time.Time `json:"canceled_on"`
}

func (e OrderCancelledEvent) IsOrderEvent() {}
