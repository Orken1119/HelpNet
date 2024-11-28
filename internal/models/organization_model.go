package models

import (
	"context"
	"time"
)

type Organization struct {
	ID                       uint       `json:"id"`
	Email                    string     `json:"email"`
	PosterUrl                *string    `json:"poster_url"`
	Name                     *string    `json:"name,omitempty"`
	Password                 string     `json:"password,omitempty"`
	PhoneNumber              *string    `json:"phoneNumber"`
	RoleID                   uint       `json:"roleId"`
	CreatedAt                *time.Time `json:"created_at"`
	City                     *string    `json:"city,omitempty"`
	Information              *string    `json:"information"`
	Direction                *string    `json:"direction"`
	FinishedProjects         *[]Event   `json:"ready_ivents"`
	Projects                 *[]Event   `json:"events"`
	VolunteerExperienceYears *int       `json:"years_of_experience"`
}

type OrganizationProfile struct {
	ID                       uint     `json:"id"`
	Email                    string   `json:"email"`
	PosterUrl                *string  `json:"poster_url"`
	Name                     *string  `json:"name,omitempty"`
	PhoneNumber              *string  `json:"phoneNumber"`
	City                     *string  `json:"city,omitempty"`
	Information              *string  `json:"information"`
	Direction                *string  `json:"direction"`
	FinishedProjects         *[]Event `json:"finished_ivents"`
	Projects                 *[]Event `json:"events"`
	VolunteerExperienceYears *int     `json:"years_of_experience"`
}

type OrganizationProfileEditing struct {
	Email                    string  `json:"email"`
	PosterUrl                *string `json:"poster_url"`
	Name                     *string `json:"name,omitempty"`
	PhoneNumber              *string `json:"phoneNumber"`
	City                     *string `json:"city,omitempty"`
	Information              *string `json:"information"`
	Direction                *string `json:"direction"`
	VolunteerExperienceYears *int    `json:"years_of_experience"`
}

type OrganizationPreview struct {
	ID        uint       `json:"id"`
	Email     string     `json:"email"`
	RoleID    uint       `json:"roleId"`
	Password  string     `json:"password,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
}

type OrganizationRequest struct {
	Email    string   `json:"email"`
	Password Password `json:"password,omitempty"`
}

type OrganizationRepository interface {
	DeleteOrganization(c context.Context, id int) error
	EditOrganizationProfile(c context.Context, orgID int, organization *OrganizationProfileEditing) error
	GetOrganizationProfile(c context.Context, id int) (*OrganizationProfile, error)
	ChangePasswordForOrganizations(c context.Context, orgID int, password string) error
	GetAllOrganizations(c context.Context) ([]OrganizationProfile, error)
	DeleteMemberFromEvent(c context.Context, userID int, eventID int) error
}

func (u *OrganizationPreview) GetID() uint {
	return u.ID
}

func (u *OrganizationPreview) GetRoleID() uint {
	return u.RoleID
}
