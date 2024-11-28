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

	// Query basic volunteer profile details
	profileQuery := `SELECT id, email, photo_url, phone_number, name, skills, city, age, grade FROM volunteers WHERE id = $1`
	row := ur.db.QueryRow(c, profileQuery, userID)
	err := row.Scan(&user.ID, &user.Email, &user.PhotoUrl, &user.PhoneNumber, &user.Name, &user.Skills, &user.City, &user.Age, &user.Grade)
	if err != nil {
		return user, err
	}

	// Query current events the volunteer is participating in
	currentEventsQuery := `
		SELECT e.id, e.event_name, e.information, e.organization_id, e.poster_url, 
		       e.preview_url, e.skill_direction, e.address, e.start_date, e.end_date,
		       e.necessary_people_count, e.members_count, e.finished
		FROM events e
		JOIN volunteer_events ve ON ve.event_id = e.id
		WHERE ve.volunteer_id = $1 AND e.finished = false`
	currentRows, err := ur.db.Query(c, currentEventsQuery, userID)
	if err != nil {
		return user, err
	}
	defer currentRows.Close()

	var currentEvents []models.Event
	for currentRows.Next() {
		var event models.Event
		if err := currentRows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
			&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
			&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
			&event.HowManyPeopleAccepted, &event.Finished); err != nil {
			return user, err
		}
		currentEvents = append(currentEvents, event)
	}
	user.EventsNow = &currentEvents

	// Query finished events the volunteer participated in
	finishedEventsQuery := `
		SELECT e.id, e.event_name, e.information, e.organization_id, e.poster_url, 
		       e.preview_url, e.skill_direction, e.address, e.start_date, e.end_date,
		       e.necessary_people_count, e.members_count, e.finished
		FROM events e
		JOIN volunteer_events ve ON ve.event_id = e.id
		WHERE ve.volunteer_id = $1 AND e.finished = true`
	finishedRows, err := ur.db.Query(c, finishedEventsQuery, userID)
	if err != nil {
		return user, err
	}
	defer finishedRows.Close()

	var finishedEvents []models.Event
	for finishedRows.Next() {
		var event models.Event
		if err := finishedRows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
			&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
			&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
			&event.HowManyPeopleAccepted, &event.Finished); err != nil {
			return user, err
		}
		finishedEvents = append(finishedEvents, event)
	}
	user.Participated = &finishedEvents
	user.Certificates, err = ur.getCertificates(c, userID)
	if err != nil {
		return user, nil
	}

	return user, nil
}

func (ur *UserRepository) GetAllVolunteers(c context.Context) ([]models.VolunteerProfile, error) {
	var volunteers []models.VolunteerProfile

	// Query all basic volunteer profiles
	profileQuery := `
		SELECT id, email, photo_url, phone_number, name, skills, city, age, grade
		FROM volunteers`
	rows, err := ur.db.Query(c, profileQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.VolunteerProfile

		// Scan basic volunteer details
		err := rows.Scan(&user.ID, &user.Email, &user.PhotoUrl, &user.PhoneNumber,
			&user.Name, &user.Skills, &user.City, &user.Age, &user.Grade)
		if err != nil {
			return nil, err
		}

		// Query and populate current events
		currentEventsQuery := `
			SELECT e.id, e.event_name, e.information, e.organization_id, e.poster_url, 
			       e.preview_url, e.skill_direction, e.address, e.start_date, e.end_date,
			       e.necessary_people_count, e.members_count, e.finished
			FROM events e
			JOIN volunteer_events ve ON ve.event_id = e.id
			WHERE ve.volunteer_id = $1 AND e.finished = false`
		currentRows, err := ur.db.Query(c, currentEventsQuery, user.ID)
		if err != nil {
			return nil, err
		}

		var currentEvents []models.Event
		for currentRows.Next() {
			var event models.Event
			if err := currentRows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
				&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
				&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
				&event.HowManyPeopleAccepted, &event.Finished); err != nil {
				return nil, err
			}
			currentEvents = append(currentEvents, event)
		}
		currentRows.Close()
		user.EventsNow = &currentEvents

		// Query and populate finished events
		finishedEventsQuery := `
			SELECT e.id, e.event_name, e.information, e.organization_id, e.poster_url, 
			       e.preview_url, e.skill_direction, e.address, e.start_date, e.end_date,
			       e.necessary_people_count, e.members_count, e.finished
			FROM events e
			JOIN volunteer_events ve ON ve.event_id = e.id
			WHERE ve.volunteer_id = $1 AND e.finished = true`
		finishedRows, err := ur.db.Query(c, finishedEventsQuery, user.ID)
		if err != nil {
			return nil, err
		}

		var finishedEvents []models.Event
		for finishedRows.Next() {
			var event models.Event
			if err := finishedRows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
				&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
				&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
				&event.HowManyPeopleAccepted, &event.Finished); err != nil {
				return nil, err
			}
			finishedEvents = append(finishedEvents, event)
		}
		finishedRows.Close()
		user.Participated = &finishedEvents

		// Query and populate certificates
		user.Certificates, err = ur.getCertificates(c, int(user.ID))
		if err != nil {
			return nil, err
		}

		// Add the user to the list of volunteers
		volunteers = append(volunteers, user)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return volunteers, nil
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
		name, 
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
	query := `INSERT INTO passwordresetcode(email, code, createdAt)
	VALUES ($1, $2, $3);`

	_, err := ur.db.Exec(c, query, email, code, currentTime)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetVolunteerByEmail(c context.Context, email string) (models.User, error) {
	var user models.User

	query := `SELECT id, email, password, created_at, role_id FROM volunteers where email = $1`
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

func (ur *UserRepository) EditVolunteerProfile(c context.Context, userID int, volunteer models.VolunteerProfileEditing) error {
	query := `UPDATE volunteers
	SET 
		email = $1,
		phone_number = $2,
		name = $3,
		skills = $4,
		city = $5,
		age = $6,
		photo_url =$7,
		direction = $8
	WHERE id = $9;
	`
	_, err := ur.db.Exec(
		c,
		query,
		volunteer.Email,
		volunteer.PhoneNumber,
		volunteer.Name,
		volunteer.Skills,
		volunteer.City,
		volunteer.Age,
		volunteer.PhotoUrl,
		volunteer.Direction,
		userID)
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
		return org, err
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

func (ur *UserRepository) AddCertificate(c context.Context, imageUrl string, userID int) error {
	query1 := `INSERT INTO volunteer_certificates (
		volunteer_id, image_url
	)
	VALUES ($1, $2)`

	_, err := ur.db.Exec(c, query1, userID, imageUrl)
	if err != nil {
		return err
	}
	return nil

}

func (ur *UserRepository) DeleteCertificate(c context.Context, id int) error {
	query := `DELETE FROM volunteer_certificates WHERE id = $1`

	_, err := ur.db.Exec(c, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) getCertificates(c context.Context, userID int) (*[]models.Certificate, error) {
	var certificates []models.Certificate

	query := `
	SELECT id, image_url FROM volunteer_certificates WHERE volunteer_id = $1`

	rows, err := ur.db.Query(c, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var certificate models.Certificate
		if err = rows.Scan(&certificate.ID, &certificate.ImageUrl); err != nil {
			return nil, err
		}
		certificates = append(certificates, certificate)
	}

	// Check for errors that may have occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &certificates, nil
}
