package types

import (
	"github.com/eugenenosenko/gopoly/internal/events"
)

type Typ interface {
	IsTyp()
}

type TypA struct {
	ID string `json:"id"`
}

func (a TypA) IsTyp() {}

type TypB struct {
	ID    string       `json:"id"`
	Event events.Event `json:"event"`
}

func (a TypB) IsTyp() {}
