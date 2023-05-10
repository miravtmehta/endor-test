package main

import (
	"reflect"
	"strings"
	"time"
)

type Object interface {
	// GetKind returns the type of the object.
	GetKind() string
	// GetID returns a unique UUID for the object.
	GetID() string
	// GetName returns the name of the object. Names are not unique.
	GetName() string
	// SetID sets the ID of the object.
	SetID(string)
	// SetName sets the name of the object.
	SetName(string)
}

type Person struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	LastName  string    `json:"last_name,omitempty"`
	Birthday  string    `json:"birthday,omitempty"`
	BirthDate time.Time `json:"birthdate,omitempty"`
}

type Animal struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Type    string `json:"type"`
	OwnerID string `json:"owner_id,omitempty"`
}

var _ Object = (*Person)(nil)
var _ Object = (*Animal)(nil)

//---------------
// Person receiver functions
//---------------

func (p *Person) GetKind() string {
	return strings.Split(reflect.TypeOf(p).String(), ".")[1]
}

func (p *Person) GetID() string {
	return p.ID
}

func (p *Person) GetName() string {
	return p.Name
}

func (p *Person) SetID(s string) {
	p.ID = s
}

func (p *Person) SetName(s string) {
	p.Name = s
}

//---------------
// Animal receiver functions
//---------------

func (a *Animal) GetKind() string {
	return strings.Split(reflect.TypeOf(a).String(), ".")[1]
}

func (a *Animal) GetID() string {
	return a.ID
}

func (a *Animal) GetName() string {
	return a.Name
}

func (a *Animal) SetID(s string) {
	a.ID = s
}

func (a *Animal) SetName(s string) {
	a.Name = s
}
