package events

import (
	"time"
)

type Event interface {
	IsEvent()
}

type OwnerEvent interface {
	IsOwner()
}

type OwnerClosedEvent struct {
	ID   string
	Type string
}

func (e OwnerClosedEvent) IsOwner() {}

type OwnerOpenEvent struct {
	ID   string
	Type string
}

func (e OwnerOpenEvent) IsOwner() {}

func (e SoldEvent) IsEvent() {}

type SoldEvent struct {
	ID   string
	Type string
}

func (e CancelledEvent) IsEvent() {}

type CancelledEvent struct {
	ID          string
	Type        string
	CancelledOn time.Time
}
