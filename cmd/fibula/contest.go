// contest.go

package main

import (
	"context"
	"time"

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

// Contest is data about a competitive event.
type Contest struct {
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
	// Centers      []byte         `json:"centers"` // we may actially provide this magically
	RegionEvents []*RegionEvent `json:"-"`
	SchoolEvents []*SchoolEvent `json:"-"`
}



// RegionEvent is an event that changes a region's state
type RegionEvent struct {
	ID          string      // Unique ID for the region.
	Registered  Registered  `datastore:"noindex" json:"registered,omitempty"`
	Renamed     Renamed     `datastore:"noindex" json:"renamed,omitempty"`
	Activated   Activated   `datastore:"noindex" json:"activated,omitempty"`
	Deactivated Deactivated `datastore:"noindex" json:"deactivated,omitempty"`
	TimeStamp   time.Time   `json:"timestamp,omitempty"`
}

// SchoolEvent is an event that changes a school's state
type SchoolEvent struct {
	ID          string      // Unique ID for the school.
	Registered  Registered  `datastore:"noindex" json:"registered,omitempty"`
	Renamed     Renamed     `datastore:"noindex" json:"renamed,omitempty"`
	Activated   Activated   `datastore:"noindex" json:"activated,omitempty"`
	Deactivated Deactivated `datastore:"noindex" json:"deactivated,omitempty"`
	Placed      Placed      `datastore:"noindex" json:"placed,omitempty"`
	TimeStamp   time.Time   `json:"timestamp,omitempty"`
}

// Registered is the initial name and parent for reporting center
type Registered struct {
	Name   string
	Parent string
}

// Renamed is when something's name attribute changes
type Renamed struct {
	Name string
}

// Activated is when something shows up in the current listings
type Activated struct {
}

// Deactivated is when something isn't active anymore
type Deactivated struct {
}

// Placed is when a school is placed in a town for name disambiguation
type Placed struct {
	Town string
}

// Region, School, State are defined in protobuffer hier.proto (as hier.pb.go)

// School represents the thing that contains seats for testing
// type School struct{

// 	Seats 	[]Seat  `json:"seats"`
// }
// Seat represents a thing that holds students and test papers

func (e *Contest) keyFromInt(db *datastore.Client, id int) *datastore.Key {
	key := datastore.IDKey("Contest", int64(id), nil)
	return key
}

func (e *Contest) getContest(db *datastore.Client, id string) error {
	context := context.Background()
	key := datastore.NameKey("Contest", id, nil)
	err := db.Get(context, key, e)
	// fmt.Printf("getEvent yielded %v", e)
	// fmt.Printf("this error is from Get: %v", err)

	return err
}

func (e *Contest) updateContest(db *datastore.Client, i string) error {
	context := context.Background()

	// eventKey := e.keyFromInt(db, i)
	eventKey := datastore.NameKey("Contest", i, nil)
	_, err := db.Put(context, eventKey, e)
	if err != nil {
		return err
	}

	return nil

}

func (e *Contest) deleteContest(db *datastore.Client, i string) error {
	context := context.Background()

	// eventKey := e.keyFromInt(db, i)
	eventKey := datastore.NameKey("Contest", i, nil)
	err := db.Delete(context, eventKey)
	if err != nil {
		return err
	}
	return nil
}

func (e *Contest) createContest(db *datastore.Client) error {
	context := context.Background()
	if e.ID == "" {
		e.ID = makeXID()
	}
	eventKey := datastore.NameKey("Contest", e.ID, nil)
	_, err := db.Put(context, eventKey, e)
	if err != nil {
		return err
	}
	return nil
}

func getContests(db *datastore.Client) ([]*Contest, error) {
	context := context.Background()

	q := datastore.NewQuery("Contest").Filter("Active=", true)

	var contests []*Contest

	_, err := db.GetAll(context, q, &contests)
	// fmt.Printf("this error is from GetEvents: %v", err)

	if err != nil {
		if err == datastore.ErrInvalidEntityType {
			// We got "invalid entity type". return empty list
			return contests, nil
		}
	}
	return contests, nil

}
