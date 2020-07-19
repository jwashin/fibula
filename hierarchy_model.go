// model.go

package main

import (
	"context"

	"cloud.google.com/go/datastore"
)

// Event is data about the competition event
type Event struct {
	ID                *datastore.Key `json:"id" datastore:"__key__"`
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

func getEvent(db *datastore.Client, key *datastore.Key) (*Event, error) {
	context := context.Background()
	var e Event

	err := db.Get(context, key, e)
	// fmt.Printf("this error is from Get: %v", err)

	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *School) updateEvent(db *datastore.Client) error {
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

func (s *School) deleteEvent(db *datastore.Client) error {
	context := context.Background()

	productKey := s.ID
	err := db.Delete(context, productKey)
	if err != nil {
		return err

	}
	return nil
}

func (s *School) createEvent(db *datastore.Client) error {
	context := context.Background()
	schoolKey := datastore.IncompleteKey("School", nil)
	_, err := db.Put(context, schoolKey, s)
	if err != nil {
		return err
	}
	return nil

}

func getEvents(db *datastore.Client, event *datastore.Key, region string) ([]School, error) {

	context := context.Background()
	query := datastore.NewQuery("School")
	if event != nil {
		query.Filter("Event=", event)
	}
	if region != "" {
		query.Filter("Region=", region)
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

/*Schools**********************************************************/

// School has the basic info about a School
type School struct {
	ID     *datastore.Key `json:"id" datastore:"__key__"`
	Name   string         `json:"name"`
	Region string         `json:"region"`
	Town   string         `json:"town"`
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

func getSchools(db *datastore.Client, event *datastore.Key, region string) ([]School, error) {

	context := context.Background()
	query := datastore.NewQuery("School")
	if event != nil {
		query.Filter("Event=", event)
	}
	if region != "" {
		query.Filter("Region=", region)
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
