package models

import (
	"context"
)

type Event struct {
	ID                    int                  `json:"id"`
	Name                  string               `json:"event_name"`
	Information           string               `json:"information"`
	OrganizationID        int                  `json:"organization_id"`
	PosterUrl             string               `json:"poster_url"`
	PreviewUrl            string               `json:"preview_url"`
	SkillsDirection       string               `json:"skill_direction"`
	Address               string               `json:"address"`
	StartingDate          string               `json:"start_data"`
	EndDate               string               `json:"end_data"`
	Members               *[]VolunteerMainInfo `json:"members_info"`
	NecCountOfPeople      int                  `json:"neccessary_people_count"`
	HowManyPeopleAccepted int                  `json:"members_count"`
	Finished              *bool                `json:"finished"`
}

type EventForCreating struct {
	Name             string `json:"event_name"`
	Information      string `json:"information"`
	OrganizationID   int    `json:"organization_id"`
	PosterUrl        string `json:"poster_url"`
	PreviewUrl       string `json:"preview_url"`
	SkillsDirection  string `json:"skill_direction"`
	Address          string `json:"address"`
	StartingDate     string `json:"start_data"`
	EndDate          string `json:"end_data"`
	NecCountOfPeople int    `json:"neccessary_people_count"`
}

type EventForEditing struct {
	Name             string `json:"event_name"`
	Information      string `json:"information"`
	PosterUrl        string `json:"poster_url"`
	PreviewUrl       string `json:"preview_url"`
	SkillsDirection  string `json:"skill_direction"`
	Address          string `json:"address"`
	StartingDate     string `json:"start_data"`
	EndDate          string `json:"end_data"`
	NecCountOfPeople int    `json:"neccessary_people_count"`
}

type EventRepository interface {
	CreateEvent(c context.Context, event *EventForCreating) (*Event, error)
	DeleteEvent(c context.Context, id int) error
	UpdateEvent(c context.Context, event *EventForEditing, eventID int) error
	GetOrganizationsInProcessEvents(c context.Context, id int) (*[]Event, error)
	GetFinishedEventsByOrganization(c context.Context, id int) (*[]Event, error)
	GetAllEvent(c context.Context) (*[]Event, error)
	GetEventById(c context.Context, id int) (*Event, error)
	FinishEvent(c context.Context, id int) error
	ParticipateEvent(c context.Context, userID int, eventID int) error
	GetVolunteerParticipatingEvents(c context.Context, userID int) (*[]Event, error)
	GetVolunteerFinishedEvents(c context.Context, userID int) (*[]Event, error)
	GetVolunteersForEvent(c context.Context, eventID int) (*[]VolunteerMainInfo, error)
}
