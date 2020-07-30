// event.go

package main

import (
	"context"

	"cloud.google.com/go/datastore"
)

// Event is data about the competition event
type Event struct {
	ID                *datastore.Key `datastore:"__key__" json:"id"`
	Active            bool           `json:"active"`
	State             string         `json:"state"`
	Level             string         `json:"level"`
	Organization      string         `json:"organization"`
	Year              string         `json:"year"`
	Title             string         `json:"title"`
	AdminEmail        string         `json:"admin_email"`
	InfoURL           string         `json:"info_url"`
	RegistrationStart string         `json:"registration_start"`
	CompetitionStart  string         `json:"competition_start"`
	RegistrationEnd   string         `json:"registration_end"`
	CompetitionEnd    string         `json:"competition_end"`
	Mode              string         `json:"mode"`
	Regions           []string       `json:"regions"`
}

func (e *Event) keyFromInt(db *datastore.Client, id int) *datastore.Key {
	key := datastore.IDKey("Event", int64(id), nil)
	return key
}

func (e *Event) getEvent(db *datastore.Client, id int) error {
	context := context.Background()
	key := e.keyFromInt(db, id)
	err := db.Get(context, key, e)
	// fmt.Printf("getEvent yielded %v", e)
	// fmt.Printf("this error is from Get: %v", err)

	return err
}

func (e *Event) updateEvent(db *datastore.Client, i int) error {
	context := context.Background()

	eventKey := e.keyFromInt(db, i)
	id, err := db.Put(context, eventKey, e)
	e.ID = id
	if err != nil {
		return err
	}

	return nil

}

func (e *Event) deleteEvent(db *datastore.Client, i int) error {
	context := context.Background()

	eventKey := e.keyFromInt(db, i)
	err := db.Delete(context, eventKey)
	if err != nil {
		return err

	}
	return nil
}

func (e *Event) createEvent(db *datastore.Client) error {
	context := context.Background()
	eventKey := datastore.IncompleteKey("Event", nil)
	_, err := db.Put(context, eventKey, e)
	if err != nil {
		return err
	}
	return nil
}

func getEvents(db *datastore.Client) ([]Event, error) {
	context := context.Background()

	q := datastore.NewQuery("Event").Filter("Active=", true)

	var events []Event

	_, err := db.GetAll(context, q, events)
	// fmt.Printf("this error is from Get: %v", err)

	if err != nil {
		if err == datastore.ErrInvalidEntityType {
			// We got "invalid entity type". return empty list
			return []Event{}, nil
		}
	}
	return events, nil

}
