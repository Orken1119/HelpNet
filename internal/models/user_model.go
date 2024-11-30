package models

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	RoleID    uint      `json:"roleId"`
}

type Volunteer struct {
	ID           uint           `json:"id"`
	Email        string         `json:"email"`
	Name         *string        `json:"name,omitempty"`
	PhotoUrl     *string        `json:"photo_url,omitempty"`
	Password     *string        `json:"password,omitempty"`
	PhoneNumber  *string        `json:"phoneNumber"`
	RoleID       uint           `json:"roleId"`
	CreatedAt    time.Time      `json:"created_at"`
	Skills       *string        `json:"skills,omitempty"`
	City         *string        `json:"city,omitempty"`
	Age          *int           `json:"age,omitempty"`
	EventsNow    *[]Event       `json:"events_now"`
	Participated *[]Event       `json:"finished"`
	Grade        *int           `json:"grade"`
	Certificates *[]Certificate `json:"certificates"`
}

type Certificate struct {
	ID       int    `json:"id"`
	ImageUrl string `json:"certificate_url"`
}

type ChangePasswordRequest struct {
	Email    string   `json:"email"`
	Code     string   `json:"code"`
	Password Password `json:"password"`
}

type VolunteerRequest struct {
	Email       string   `json:"email"`
	Password    Password `json:"password,omitempty"`
	PhoneNumber string   `json:"phoneNumber"`
	Name        string   `json:"name"`
	Skills      string   `json:"skills"`
	City        string   `json:"city"`
	Age         int      `json:"age"`
}

type VolunteerSignUp struct {
	Email    string   `json:"email"`
	Password Password `json:"password,omitempty"`
}

type VolunteerProfile struct {
	ID           uint           `json:"id"`
	PhotoUrl     *string        `json:"photo_url,omitempty"`
	Email        string         `json:"email"`
	PhoneNumber  *string        `json:"phoneNumber"`
	Name         *string        `json:"name"`
	Skills       *string        `json:"skills"`
	City         *string        `json:"city"`
	Age          *int           `json:"age"`
	Direction    *string        `json:"direction"`
	Grade        *int           `json:"grade"`
	EventsNow    *[]Event       `json:"events_now"`
	Participated *[]Event       `json:"finished"`
	Certificates *[]Certificate `json:"certificates"`
}

type VolunteerProfileEditing struct {
	PhotoUrl    *string `json:"photo_url,omitempty"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	Name        *string `json:"name"`
	Skills      *string `json:"skills"`
	City        *string `json:"city"`
	Age         *int    `json:"age"`
	Direction   *string `json:"direction"`
}

type VolunteerMainInfo struct {
	ID          uint    `json:"id"`
	PhotoUrl    *string `json:"photo_url,omitempty"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	Name        *string `json:"name"`
	Skills      *string `json:"skills"`
	City        *string `json:"city"`
	Age         *int    `json:"age"`
	Direction   *string `json:"direction"`
	Grade       *int    `json:"grade"`
}

var ErrEmailAlreadyExists = errors.New("email already exists")

type LoginRequest struct {
	Email    string   `json:"email"`
	Password Password `json:"password"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Email    string   `json:"email"`
	Password Password `json:"password"`
}

type Password struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type UserRepository interface {
	CreateOrganization(c context.Context, request *OrganizationRequest) (int, error)
	CreateUserVolunteer(c context.Context, request *SignUpRequest) (int, error)
	CreateUserOrganization(c context.Context, request *SignUpRequest) (int, error)
	GetOrganizationByEmail(c context.Context, email string) (User, error)
	GetVolunteerProfile(c context.Context, userID int) (VolunteerProfile, error)
	GetCodeByEmail(c context.Context, email string) (string, error)
	CreateVolunteer(c context.Context, user VolunteerRequest) (int, error)
	CreatePasswordResetCode(c context.Context, email string, code string) error
	GetVolunteerByEmail(c context.Context, email string) (User, error)
	ChangeForgottenVolunteersPassword(c context.Context, code string, email string, newPassword string) error
	ChangeForgottenOrgPassword(c context.Context, code string, email string, newPassword string) error
	ChangePassword(c context.Context, userID int, password string) error
	ChangePasswordForOrg(c context.Context, orgID int, password string) error
	EditVolunteerProfile(c context.Context, userID int, volunteer VolunteerProfileEditing) error
	GetAllVolunteers(c context.Context) ([]VolunteerProfile, error)

	AddCertificate(c context.Context, imageUrl string, userID int) error
	DeleteCertificate(c context.Context, id int) error

	SearchEvent(c context.Context, name string) (*[]Event, error)
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) GetRoleID() uint {
	return u.RoleID
}
