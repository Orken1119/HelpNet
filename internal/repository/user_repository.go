package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Orken1119/HelpNet/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) models.UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetVolunteerProfile(c context.Context, userID int) (models.VolunteerProfile, error) {
	var user models.VolunteerProfile

	query := `SELECT id, email, phone_number, volunteer_name, skills, city, age FROM volunteers where id = $1`
	row := ur.db.QueryRow(c, query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.PhoneNumber, &user.Name, &user.Skills, &user.City, &user.Age)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (ur *UserRepository) GetCodeByEmail(c context.Context, email string) (string, error) {
	var code string

	query := `SELECT code FROM PasswordResetCode where email = $1`
	row := ur.db.QueryRow(c, query, email)
	err := row.Scan(&code)

	if err != nil {
		return "", err
	}
	return code, nil

}

func (ur *UserRepository) CreateVolunteer(c context.Context, volunteer models.VolunteerRequest) (int, error) {
	var userID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO volunteers(
		email, 
		password, 
		phone_number, 
		volunteer_name, 
		created_at, 
		skills, 
		city, 
		age, 
		role_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;`

	err := ur.db.QueryRow(c,
		userQuery,
		volunteer.Email,
		volunteer.Password.Password,
		volunteer.PhoneNumber,
		volunteer.Name,
		currentTime,
		volunteer.Skills,
		volunteer.City,
		volunteer.Age,
		2).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (ur *UserRepository) ChangeForgottenVolunteersPassword(c context.Context, code string, email string, newPassword string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	var storedCode string

	codeQuery := `SELECT code FROM PasswordResetCode WHERE email = $1`

	err := ur.db.QueryRow(c, codeQuery, email).Scan(&storedCode)

	if err != nil {
		return fmt.Errorf("failed to retrieve reset code: %v", err)
	}

	if storedCode != code {
		return fmt.Errorf("invalid reset code")
	}

	updateQuery := `UPDATE volunteers 
	SET password = $1, created_at = $2 
	WHERE email = $3 
	RETURNING email;`

	var updatedEmail string
	err = ur.db.QueryRow(c, updateQuery, newPassword, currentTime, email).Scan(&updatedEmail)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) CreatePasswordResetCode(c context.Context, email string, code string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO passwordresetcode(email, code, created_at)
	VALUES ($1, $2, $3);`

	_, err := ur.db.Exec(c, query, email, code, currentTime)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetVolunteerByEmail(c context.Context, email string) (models.User, error) {
	var user models.User

	query := `SELECT id, email, password, created_at, name, role_id FROM volunteers where email = $1`
	row := ur.db.QueryRow(c, query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.RoleID)

	if err != nil {
		return user, nil
	}
	return user, nil

}

func (ur *UserRepository) ChangePassword(c context.Context, userID int, password string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `UPDATE volunteers
	SET
	password = $1,
	created_at = $2
	where
	id = $3`
	_, err := ur.db.Exec(c, query, password, currentTime, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) EditVolunteerProfile(c context.Context, userID int, volunteer models.VolunteerProfile) error {
	query := `UPDATE volunteers
	SET 
		email = $1,
		phone_number = $2,
		volunteer_name = $3,
		skills = $4,
		city = $5,
		age = $6
	WHERE id = $7 ;
	`
	_, err := ur.db.Exec(c, query, volunteer.Email, volunteer.PhoneNumber, volunteer.Name, volunteer.Skills, volunteer.City, volunteer.Age, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) CreateOrganization(c context.Context, request *models.OrganizationRequest) (int, error) {
	var userID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO organizations(
		email, 
		password,  
		created_at,
		role_id)
		VALUES ($1, $2, $3, $4) returning id;`

	err := ur.db.QueryRow(c,
		query,
		request.Email,
		request.Password.Password,
		currentTime,
		1,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (ur *UserRepository) GetOrganizationByEmail(c context.Context, email string) (models.User, error) {
	var org models.User

	query := `SELECT id, email, password, created_at, role_id FROM organizations where email = $1`
	row := ur.db.QueryRow(c, query, email)
	err := row.Scan(&org.ID, &org.Email, &org.Password, &org.CreatedAt, &org.RoleID)

	if err != nil {
		return org, nil
	}
	return org, nil

}

func (ur *UserRepository) ChangeForgottenOrgPassword(c context.Context, code string, email string, newPassword string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	var storedCode string

	codeQuery := `SELECT code FROM PasswordResetCode WHERE email = $1`

	err := ur.db.QueryRow(c, codeQuery, email).Scan(&storedCode)

	if err != nil {
		return fmt.Errorf("failed to retrieve reset code: %v", err)
	}

	if storedCode != code {
		return fmt.Errorf("invalid reset code")
	}

	updateQuery := `UPDATE organizations 
	SET password = $1, created_at = $2 
	WHERE email = $3 
	RETURNING email;`

	var updatedEmail string
	err = ur.db.QueryRow(c, updateQuery, newPassword, currentTime, email).Scan(&updatedEmail)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) ChangePasswordForOrg(c context.Context, orgID int, password string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `UPDATE organizations
	SET
	password = $1,
	created_at = $2
	where
	id = $3`
	_, err := ur.db.Exec(c, query, password, currentTime, orgID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) CreateUserOrganization(c context.Context, request *models.SignUpRequest) (int, error) {
	var userID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO organizations(
		email, 
		password,  
		created_at,
		role_id)
		VALUES ($1, $2, $3, $4) returning id;`

	err := ur.db.QueryRow(c,
		query,
		request.Email,
		request.Password.Password,
		currentTime,
		1,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (ur *UserRepository) CreateUserVolunteer(c context.Context, request *models.SignUpRequest) (int, error) {
	var userID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO volunteers(
		email, 
		password,  
		created_at,
		role_id)
		VALUES ($1, $2, $3, $4) returning id;`

	err := ur.db.QueryRow(c,
		query,
		request.Email,
		request.Password.Password,
		currentTime,
		1,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
