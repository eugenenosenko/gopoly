package models

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/eugenenosenko/gopoly/internal/events"
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
	ID       string       `json:"id"`
	Owner    Owner        `json:"owner"`
	Property Property     `json:"property"`
	Type     string       `json:"kind"`
	Event    events.Event `json:"event"`
}

func (a SellAdvert) IsAdvertBase() {}

type RentAdvert struct {
	ID       string       `json:"id"`
	Owner    Owner        `json:"owner"`
	Property Property     `json:"property"`
	Type     string       `json:"kind"`
	Event    events.Event `json:"event"`
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

func UnmarshalAdvertBaseJSON(data []byte) (AdvertBase, error) {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil, nil
	}
	var probe struct {
		Discriminator string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("unmarshal AdvertBase type: %v", err)
	}
	switch probe.Discriminator {
	case "sell":
		var v SellAdvert
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("unmarshal SellAdvert: %v", err)
		}
		return v, nil
	case "rent":
		var v RentAdvert
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("unmarshal RentAdvert: %v", err)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("could not unmarshal AdvertBase JSON: unknown variant %q", probe.Discriminator)
	}
}

func UnmarshalPropertyJSON(data []byte) (Property, error) {
	var err error
	matches := make([]Property, 0)

	newStrictDecoder := func(data []byte) *json.Decoder {
		dec := json.NewDecoder(bytes.NewBuffer(data))
		dec.DisallowUnknownFields()
		return dec
	}

	houseProperty := &HouseProperty{}
	err = newStrictDecoder(data).Decode(&houseProperty)
	if err == nil {
		jsonRentAdvert, _ := json.Marshal(houseProperty)
		if string(jsonRentAdvert) == "{}" { // empty struct
			houseProperty = nil
		} else {
			matches = append(matches, houseProperty)
		}
	} else {
		houseProperty = nil
	}

	flatProperty := &FlatProperty{}
	err = newStrictDecoder(data).Decode(&flatProperty)
	if err == nil {
		jsonSellAdvert, _ := json.Marshal(flatProperty)
		if string(jsonSellAdvert) == "{}" { // empty struct
			flatProperty = nil
		} else {
			matches = append(matches, flatProperty)
		}
	} else {
		flatProperty = nil
	}

	if len(matches) > 1 { // more than 1 match
		return nil, fmt.Errorf("data matches more than one schema in oneOf(Property)")
	} else if len(matches) == 1 {
		return matches[0], nil // exactly one match
	} else { // no match
		return nil, fmt.Errorf("data failed to match schemas in oneOf(Property)")
	}
}
