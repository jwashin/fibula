// model.go

package main

import (
	"context"
	"sort"

	"cloud.google.com/go/datastore"
)

// CompetitionEvent is data about the competition event
type CompetitionEvent struct {
	ID                *datastore.Key `json:"id" datastore:"__key__"`
	Active            bool           `json:"Active"`
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

func getEvent(db *datastore.Client, key *datastore.Key) (*CompetitionEvent, error) {
	context := context.Background()
	var e CompetitionEvent

	err := db.Get(context, key, e)
	// fmt.Printf("this error is from Get: %v", err)

	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (e *CompetitionEvent) updateEvent(db *datastore.Client) error {
	context := context.Background()
	if e.ID != nil {
		eventKey := e.ID
		_, err := db.Put(context, eventKey, e)
		if err != nil {
			return err
		}
	}
	return nil

}

func (e *CompetitionEvent) deleteEvent(db *datastore.Client) error {
	context := context.Background()

	eventKey := e.ID
	err := db.Delete(context, eventKey)
	if err != nil {
		return err

	}
	return nil
}

func (e *CompetitionEvent) createEvent(db *datastore.Client) error {
	context := context.Background()
	eventKey := datastore.IncompleteKey("CompetitionEvent", nil)
	_, err := db.Put(context, eventKey, e)
	if err != nil {
		return err
	}
	return nil

}

func (e *CompetitionEvent) createRegion(db *datastore.Client, region string) []string {

	e.Regions = append(e.Regions, region)
	sort.Strings(e.Regions)
	db.Put(context.Background(), e.ID, e)

	return e.Regions
}

func getEvents(db *datastore.Client, active bool) ([]CompetitionEvent, error) {

	context := context.Background()
	query := datastore.NewQuery("CompetitionEvent")
	if active != false {
		query.Filter("Active=", active)
	}

	var events []CompetitionEvent
	_, err := db.GetAll(context, query, events)

	if err != nil {
		if err == datastore.ErrInvalidEntityType {
			// We got "invalid entity type". return empty list
			return []CompetitionEvent{}, nil
		}
	}
	return events, err
}

/*Schools**********************************************************/

// School has the basic info about a School
type School struct {
	ID     *datastore.Key `json:"id" datastore:"__key__"`
	Name   string         `json:"name"`
	Region string         `json:"region"`
	Town   string         `json:"town"`
	Active bool           `json:"active"`
	Event  *datastore.Key
}

func getSchool(db *datastore.Client, key *datastore.Key) (*School, error) {
	context := context.Background()
	var s School

	// productKey := datastore.NameKey("School", s.ID, nil)
	err := db.Get(context, key, s)
	// t, err := json.Marshal(p)
	// fmt.Printf("this error is from Get: %v", err)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *School) updateSchool(db *datastore.Client) error {
	context := context.Background()
	if s.ID != nil {
		productKey := s.ID
		_, err := db.Put(context, productKey, s)
		if err != nil {
			return err
		}
	}
	return nil

}

func (s *School) deleteSchool(db *datastore.Client) error {
	context := context.Background()

	productKey := s.ID
	err := db.Delete(context, productKey)
	if err != nil {
		return err

	}
	return nil
}

func (s *School) createSchool(db *datastore.Client) error {
	context := context.Background()
	schoolKey := datastore.IncompleteKey("School", nil)
	_, err := db.Put(context, schoolKey, s)
	if err != nil {
		return err
	}
	return nil

}

func getSchools(db *datastore.Client, event *datastore.Key, region string, active bool) ([]School, error) {

	context := context.Background()
	query := datastore.NewQuery("School")
	if event != nil {
		query.Filter("Event=", event)
	}
	if region != "" {
		query.Filter("Region=", region)
	}
	if active != false {
		query.Filter("Active=", active)
	}
	var schools []School
	_, err := db.GetAll(context, query, schools)

	if err != nil {
		if err == datastore.ErrInvalidEntityType {
			// we got "invalid entity type". return empty list
			return []School{}, nil
		}
	}
	return schools, err
}
