package models

import (
	e "github.com/eugenenosenko/gopoly/internal/events"
	"github.com/eugenenosenko/gopoly/internal/types"
)

type AdvertBase interface {
	IsAdvertBase()
}

type Owner interface {
	IsOwner()
}

type Property interface {
	IsProperty()
}

type Contact interface {
	IsContact()
}

type SellAdvert struct {
	ID            string               `json:"id"`
	Owners        []Owner              `json:"owners"`
	Property      Property             `json:"property"`
	PropertiesMap map[string]Property  `json:"properties_map"`
	Typ           types.Typ            `json:"typ"`
	TypMap        map[string]types.Typ `json:"typ_map"`
	Type          string               `json:"kind"`
	Events        []e.Event            `json:"events"`
	EventsMap     map[string]e.Event   `json:"events_map"`
	OwnerEvent    e.OwnerEvent         `json:"owner_event"`
}

func (a SellAdvert) IsAdvertBase() {}

type RentAdvert struct {
	ID       string   `json:"id"`
	Owner    Owner    `json:"owner"`
	Property Property `json:"property"`
	Type     string   `json:"kind"`
	Event    e.Event  `json:"event"`
}

func (a RentAdvert) IsAdvertBase() {}

type IndividualOwner struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Contact Contact `json:"contact"`
}

func (o IndividualOwner) IsOwner() {}

type DeveloperOwner struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Contact []Contact `json:"contacts"`
}

func (o DeveloperOwner) IsOwner() {}

type AgencyOwner struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Contact []Contact `json:"contacts"`
}

func (o AgencyOwner) IsOwner() {}

type HouseProperty struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Contact Contact `json:"contact"`
}

func (p HouseProperty) IsProperty() {}

type FlatProperty struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (f FlatProperty) IsProperty() {}

type BusinessContact struct {
	ID           string `json:"id"`
	BusinessName string `json:"business_name"`
}

func (c BusinessContact) IsContact() {}

type FullName struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type PrivateContact struct {
	ID       string   `json:"id"`
	FullName FullName `json:"full_name"`
}

type PolymorphicField struct {
}

type Type struct {
	Name       string
	PolyFields []*PolymorphicField
}

func (c PrivateContact) IsContact() {}
