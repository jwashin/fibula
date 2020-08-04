// event.go

package main

import (
	"context"

	"github.com/sethvargo/go-password/password"

	"cloud.google.com/go/datastore"
)

func makeXID() string {
	a2, _ := password.Generate(2, 0, 0, true, false)
	d2, _ := password.Generate(2, 2, 0, false, false)
	a2b, _ := password.Generate(2, 0, 0, true, false)
	d2b, _ := password.Generate(2, 2, 0, false, false)

	return a2 + d2 + a2b + d2b
}

// Event is data about the competition event
type Event struct {
	ID                string `json:"id"`
	Active            bool   `json:"active"`
	State             string `json:"state"`
	Level             string `json:"level"`
	Organization      string `json:"organization"`
	Year              string `json:"year"`
	Title             string `json:"title"`
	AdminEmail        string `json:"admin_email"`
	InfoURL           string `json:"info_url"`
	RegistrationStart string `json:"registration_start"`
	CompetitionStart  string `json:"competition_start"`
	RegistrationEnd   string `json:"registration_end"`
	CompetitionEnd    string `json:"competition_end"`
	Mode              string `json:"mode"`
	// Centers is compressed JSON of the hierarchy of regions and schools.
	Centers []byte `json:"centers"`
}

// Region represents the thing that contains Schools
type Region struct {
	Schools []School `json:"schools"`
}

// School represents the thing that contains seats for testing
// type School struct{

// 	Seats 	[]Seat  `json:"seats"`
// }
// Seat represents a thing that holds students and test papers

func (e *Event) keyFromInt(db *datastore.Client, id int) *datastore.Key {
	key := datastore.IDKey("Event", int64(id), nil)
	return key
}

func (e *Event) getEvent(db *datastore.Client, id string) error {
	context := context.Background()
	key := datastore.NameKey("Event", id, nil)
	err := db.Get(context, key, e)
	// fmt.Printf("getEvent yielded %v", e)
	// fmt.Printf("this error is from Get: %v", err)

	return err
}

func (e *Event) updateEvent(db *datastore.Client, i string) error {
	context := context.Background()

	// eventKey := e.keyFromInt(db, i)
	eventKey := datastore.NameKey("Event", i, nil)
	_, err := db.Put(context, eventKey, e)
	if err != nil {
		return err
	}

	return nil

}

func (e *Event) deleteEvent(db *datastore.Client, i string) error {
	context := context.Background()

	// eventKey := e.keyFromInt(db, i)
	eventKey := datastore.NameKey("Event", i, nil)
	err := db.Delete(context, eventKey)
	if err != nil {
		return err
	}
	return nil
}

func (e *Event) createEvent(db *datastore.Client) error {
	context := context.Background()
	if e.ID == "" {
		e.ID = makeXID()
	}
	eventKey := datastore.NameKey("Event", e.ID, nil)
	_, err := db.Put(context, eventKey, e)
	if err != nil {
		return err
	}
	return nil
}

func getEvents(db *datastore.Client) ([]*Event, error) {
	context := context.Background()

	q := datastore.NewQuery("Event").Filter("Active=", true)

	var events []*Event

	_, err := db.GetAll(context, q, &events)
	// fmt.Printf("this error is from GetEvents: %v", err)

	if err != nil {
		if err == datastore.ErrInvalidEntityType {
			// We got "invalid entity type". return empty list
			return events, nil
		}
	}
	return events, nil

}
