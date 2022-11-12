package events

type Event interface {
	IsEvent()
}

func (e SoldEvent) IsEvent() {}

type SoldEvent struct {
	ID   string
	Type string
}

func (e CancelledEvent) IsEvent() {}

type CancelledEvent struct {
	ID   string
	Type string
}
