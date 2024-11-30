package repository

import (
	"context"
	"time"

	"github.com/Orken1119/HelpNet/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
)

type OrganizationRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) models.OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (op *OrganizationRepository) DeleteMemberFromEvent(c context.Context, userID int, eventID int) error {
	query := `
        DELETE FROM volunteer_events
        WHERE volunteer_id = $1 AND event_id = $2
    `

	_, err := op.db.Exec(c, query, userID, eventID)
	if err != nil {
		return err
	}

	updateEventQuery := `
        UPDATE events
        SET members_count = members_count - 1
        WHERE id = $1 AND members_count > 0
    `
	_, err = op.db.Exec(c, updateEventQuery, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (op *OrganizationRepository) GetAllOrganizations(c context.Context) ([]models.OrganizationProfile, error) {
	var organizations []models.OrganizationProfile

	query := `
		SELECT id, email, poster_url, name, phone_number, city, information, direction, volunteer_experience_years
		FROM organizations`
	rows, err := op.db.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var org models.OrganizationProfile

		err := rows.Scan(
			&org.ID,
			&org.Email,
			&org.PosterUrl,
			&org.Name,
			&org.PhoneNumber,
			&org.City,
			&org.Information,
			&org.Direction,
			&org.VolunteerExperienceYears,
		)
		if err != nil {
			return nil, err
		}

		org.FinishedProjects, err = op.getProjectsByStatus(c, int(org.ID), true)
		if err != nil {
			return nil, err
		}

		org.Projects, err = op.getProjectsByStatus(c, int(org.ID), false)
		if err != nil {
			return nil, err
		}

		organizations = append(organizations, org)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return organizations, nil
}

func (op *OrganizationRepository) DeleteOrganization(c context.Context, id int) error {
	query := `DELETE FROM organizations WHERE id = $1`
	_, err := op.db.Exec(c, query, id)
	return err
}
func (op *OrganizationRepository) EditOrganizationProfile(c context.Context, orgID int, organization *models.OrganizationProfileEditing) error {
	query := `UPDATE organizations
	SET 
		email = $1,
		poster_url = $2,
		name = $3,
		phone_number = $4,
		city = $5,
		information = $6,
		direction = $7,
		volunteer_experience_years = $8
	WHERE id = $9;
	`
	_, err := op.db.Exec(
		c,
		query,
		organization.Email,
		organization.PosterUrl,
		organization.Name,
		organization.PhoneNumber,
		organization.City,
		organization.Information,
		organization.Direction,
		organization.VolunteerExperienceYears,
		orgID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return models.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}
func (op *OrganizationRepository) GetOrganizationProfile(c context.Context, id int) (*models.OrganizationProfile, error) {
	var org models.OrganizationProfile

	query := `
		SELECT id, email, poster_url, name, phone_number, city, information, direction, volunteer_experience_years
		FROM organizations
		WHERE id = $1`
	row := op.db.QueryRow(c, query, id)

	err := row.Scan(
		&org.ID,
		&org.Email,
		&org.PosterUrl,
		&org.Name,
		&org.PhoneNumber,
		&org.City,
		&org.Information,
		&org.Direction,
		&org.VolunteerExperienceYears,
	)
	if err != nil {
		return nil, err
	}

	org.FinishedProjects, err = op.getProjectsByStatus(c, id, true)
	if err != nil {
		return nil, err
	}

	org.Projects, err = op.getProjectsByStatus(c, id, false)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

func (op *OrganizationRepository) getProjectsByStatus(c context.Context, organizationID int, finished bool) (*[]models.Event, error) {
	query := `
		SELECT e.id, e.event_name, e.information, e.organization_id, e.poster_url, e.preview_url, 
		       e.skill_direction, e.address, e.start_date, e.end_date, 
		       e.necessary_people_count, e.members_count, e.finished
		FROM events e
		WHERE e.organization_id = $1 AND e.finished = $2`
	rows, err := op.db.Query(c, query, organizationID, finished)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Information,
			&event.OrganizationID,
			&event.PosterUrl,
			&event.PreviewUrl,
			&event.SkillsDirection,
			&event.Address,
			&event.StartingDate,
			&event.EndDate,
			&event.NecCountOfPeople,
			&event.HowManyPeopleAccepted,
			&event.Finished,
		)
		if err != nil {
			return nil, err
		}

		event.Members, err = op.getEventMembers(c, event.ID)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return &events, nil
}

func (op *OrganizationRepository) getEventMembers(c context.Context, eventID int) (*[]models.VolunteerMainInfo, error) {
	query := `
		SELECT v.id, v.name, v.photo_url, v.phone_number, v.city, v.skills, v.age, v.direction, v.grade
		FROM volunteers v
		JOIN volunteer_events ve ON v.id = ve.volunteer_id
		WHERE ve.event_id = $1`
	rows, err := op.db.Query(c, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.VolunteerMainInfo
	for rows.Next() {
		var member models.VolunteerMainInfo
		err := rows.Scan(
			&member.ID,
			&member.Name,
			&member.PhotoUrl,
			&member.PhoneNumber,
			&member.City,
			&member.Skills,
			&member.Age,
			&member.Direction,
			&member.Grade,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return &members, nil
}

func (op *OrganizationRepository) ChangePasswordForOrganizations(c context.Context, orgID int, password string) error {

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `UPDATE organizations
	SET
	password = $1,
	created_at = $2
	where
	id = $3`
	_, err := op.db.Exec(c, query, password, currentTime, orgID)
	if err != nil {
		return err
	}

	return nil

}
