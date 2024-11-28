package repository

import (
	"context"

	"github.com/Orken1119/HelpNet/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) models.EventRepository {
	return &EventRepository{db: db}
}

// Event Repository methods
func (or *EventRepository) CreateEvent(c context.Context, event *models.EventForCreating) (*models.Event, error) {
	query := `INSERT INTO events (event_name, information, organization_id, poster_url, preview_url, 
                               skill_direction, address, start_date, end_date, necessary_people_count) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	var id int
	err := or.db.QueryRow(c, query, event.Name, event.Information, event.OrganizationID, event.PosterUrl, event.PreviewUrl,
		event.SkillsDirection, event.Address, event.StartingDate, event.EndDate, event.NecCountOfPeople).Scan(&id)
	if err != nil {
		return nil, err
	}

	finished := false
	return &models.Event{
		ID:                    id,
		Name:                  event.Name,
		Information:           event.Information,
		OrganizationID:        event.OrganizationID,
		PosterUrl:             event.PosterUrl,
		PreviewUrl:            event.PreviewUrl,
		SkillsDirection:       event.SkillsDirection,
		Address:               event.Address,
		StartingDate:          event.StartingDate,
		EndDate:               event.EndDate,
		NecCountOfPeople:      event.NecCountOfPeople,
		HowManyPeopleAccepted: 0,
		Finished:              &finished,
	}, nil
}

func (or *EventRepository) DeleteEvent(c context.Context, id int) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := or.db.Exec(c, query, id)
	return err
}
func (or *EventRepository) GetEventsByDirection(c context.Context, direction string) (*[]models.Event, error) {
	query := `
        SELECT 
            id, event_name, information, organization_id, poster_url, 
            preview_url, skill_direction, address, start_date, end_date, 
            necessary_people_count, members_count, finished
        FROM events 
        WHERE skill_direction = $1 AND finished = false
    `
	rows, err := or.db.Query(c, query, direction)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
			&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
			&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
			&event.HowManyPeopleAccepted, &event.Finished); err != nil {
			return nil, err
		}

		// Fetch the list of volunteers for the current event
		members, err := or.GetVolunteersForEvent(c, event.ID)
		if err != nil {
			return nil, err
		}
		event.Members = members

		events = append(events, event)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &events, nil
}

func (or *EventRepository) UpdateEvent(c context.Context, event *models.EventForEditing, eventID int) error {
	query := `UPDATE events SET event_name = $1, information = $2, poster_url = $3, preview_url = $4, 
              skill_direction = $5, address = $6, start_date = $7, end_date = $8, 
              necessary_people_count = $9
              WHERE id = $10`

	_, err := or.db.Exec(c, query, event.Name, event.Information, event.PosterUrl, event.PreviewUrl,
		event.SkillsDirection, event.Address, event.StartingDate, event.EndDate,
		event.NecCountOfPeople, eventID)
	return err
}

func (or *EventRepository) GetOrganizationsInProcessEvents(c context.Context, organizationID int) (*[]models.Event, error) {
	query := `
        SELECT 
            e.id, e.event_name, e.information, e.organization_id, e.poster_url, 
            e.preview_url, e.skill_direction, e.address, e.start_date, e.end_date, 
            e.necessary_people_count, e.members_count, e.finished
        FROM events e
        WHERE e.organization_id = $1 AND e.finished = false
    `

	rows, err := or.db.Query(c, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
			&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
			&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
			&event.HowManyPeopleAccepted, &event.Finished); err != nil {
			return nil, err
		}

		// Получаем список волонтеров для текущего события
		members, err := or.GetVolunteersForEvent(c, event.ID)
		if err != nil {
			return nil, err
		}
		event.Members = members

		events = append(events, event)
	}

	return &events, nil
}

func (or *EventRepository) GetAllEvent(c context.Context) (*[]models.Event, error) {
	query := `
        SELECT 
            id, event_name, information, organization_id, poster_url, 
            preview_url, skill_direction, address, start_date, end_date, 
            necessary_people_count, members_count, finished
        FROM events 
        WHERE finished = false
    `
	rows, err := or.db.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
			&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
			&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
			&event.HowManyPeopleAccepted, &event.Finished); err != nil {
			return nil, err
		}

		// Получаем список волонтеров для текущего события
		members, err := or.GetVolunteersForEvent(c, event.ID)
		if err != nil {
			return nil, err
		}
		event.Members = members

		events = append(events, event)
	}

	return &events, nil
}

func (or *EventRepository) GetEventById(c context.Context, id int) (*models.Event, error) {
	query := `
        SELECT 
            id, event_name, information, organization_id, poster_url, 
            preview_url, skill_direction, address, start_date, end_date, 
            necessary_people_count, members_count, finished
        FROM events 
        WHERE id = $1
    `

	var event models.Event
	err := or.db.QueryRow(c, query, id).Scan(
		&event.ID, &event.Name, &event.Information, &event.OrganizationID,
		&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection,
		&event.Address, &event.StartingDate, &event.EndDate,
		&event.NecCountOfPeople, &event.HowManyPeopleAccepted, &event.Finished,
	)
	if err != nil {
		return nil, err
	}

	// Получение списка волонтеров для данного события
	members, err := or.GetVolunteersForEvent(c, event.ID)
	if err != nil {
		return nil, err
	}
	event.Members = members

	return &event, nil
}

// FinishEvent marks the event as finished
func (or *EventRepository) FinishEvent(c context.Context, id int) error {
	query := `UPDATE events SET finished = true WHERE id = $1`
	_, err := or.db.Exec(c, query, id)
	return err
}

// ParticipateEvent allows a volunteer to participate in an event
func (or *EventRepository) ParticipateEvent(c context.Context, userID int, eventID int) error {
	// Добавление участника в таблицу volunteer_events
	participateQuery := `
		INSERT INTO volunteer_events (volunteer_id, event_id) 
		VALUES ($1, $2);`
	_, err := or.db.Exec(c, participateQuery, userID, eventID)
	if err != nil {
		return err
	}

	// Увеличение количества участников в таблице events
	updateMembersQuery := `
		UPDATE events 
		SET members_count = members_count + 1 
		WHERE id = $1 AND members_count < necessary_people_count`
	_, err = or.db.Exec(c, updateMembersQuery, eventID)
	if err != nil {
		return err
	}

	return nil
}

// GetVolunteerParticipatingEvents retrieves events a volunteer is currently participating in
func (or *EventRepository) GetVolunteerParticipatingEvents(c context.Context, userID int) (*[]models.Event, error) {
	query := `SELECT e.* FROM events e 
              JOIN volunteer_events ve ON ve.event_id = e.id 
              WHERE ve.volunteer_id = $1 AND ve.in_process = true`
	rows, err := or.db.Query(c, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
			&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
			&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
			&event.HowManyPeopleAccepted, &event.Finished); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return &events, nil
}

// GetVolunteerFinishedEvents retrieves events a volunteer has finished
func (or *EventRepository) GetVolunteerFinishedEvents(c context.Context, userID int) (*[]models.Event, error) {
	// Основной запрос для получения завершенных событий волонтера
	query := `SELECT e.id, e.event_name, e.information, e.organization_id, 
                     e.poster_url, e.preview_url, e.skill_direction, e.address, 
                     e.start_date, e.end_date, e.necessary_people_count, 
                     e.members_count, e.finished 
              FROM events e 
              JOIN volunteer_events ve ON ve.event_id = e.id 
              WHERE ve.volunteer_id = $1 AND e.finished = true`
	rows, err := or.db.Query(c, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		// Считываем данные события
		if err := rows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
			&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
			&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
			&event.HowManyPeopleAccepted, &event.Finished); err != nil {
			return nil, err
		}

		// Подзапрос для получения участников события
		membersQuery := `SELECT v.id, v.name, v.email, v.photo_url, v.skills, v.city, v.age, v.direction, v.grade 
                         FROM volunteers v 
                         JOIN volunteer_events ve ON ve.volunteer_id = v.id 
                         WHERE ve.event_id = $1`
		memberRows, err := or.db.Query(c, membersQuery, event.ID)
		if err != nil {
			return nil, err
		}

		var members []models.VolunteerMainInfo
		for memberRows.Next() {
			var member models.VolunteerMainInfo
			if err := memberRows.Scan(&member.ID, &member.Name, &member.Email, &member.PhotoUrl,
				&member.Skills, &member.City, &member.Age, &member.Direction, &member.Grade); err != nil {
				memberRows.Close()
				return nil, err
			}
			members = append(members, member)
		}
		memberRows.Close()

		// Присваиваем участников событию
		event.Members = &members

		// Добавляем событие в результат
		events = append(events, event)
	}

	return &events, nil
}

func (or *EventRepository) GetVolunteersForEvent(c context.Context, eventID int) (*[]models.VolunteerMainInfo, error) {
	query := `
        SELECT 
            v.id, v.email, v.name, v.photo_url, v.phone_number, v.skills, 
            v.city, v.age, v.grade, v.direction
        FROM volunteers v
        INNER JOIN volunteer_events ve ON ve.volunteer_id = v.id
        WHERE ve.event_id = $1
    `

	rows, err := or.db.Query(c, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.VolunteerMainInfo
	for rows.Next() {
		var member models.VolunteerMainInfo
		if err := rows.Scan(&member.ID, &member.Email, &member.Name, &member.PhotoUrl,
			&member.PhoneNumber, &member.Skills, &member.City, &member.Age, &member.Grade, &member.Direction); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return &members, nil
}

func (or *EventRepository) GetFinishedEventsByOrganization(c context.Context, id int) (*[]models.Event, error) {
	query := `
        SELECT 
            e.id, e.event_name, e.information, e.organization_id, e.poster_url, 
            e.preview_url, e.skill_direction, e.address, e.start_date, e.end_date, 
            e.necessary_people_count, e.members_count, e.finished
        FROM events e
        WHERE e.organization_id = $1 AND e.finished = true
    `

	rows, err := or.db.Query(c, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Information, &event.OrganizationID,
			&event.PosterUrl, &event.PreviewUrl, &event.SkillsDirection, &event.Address,
			&event.StartingDate, &event.EndDate, &event.NecCountOfPeople,
			&event.HowManyPeopleAccepted, &event.Finished); err != nil {
			return nil, err
		}

		// Получаем список волонтеров для текущего события
		members, err := or.GetVolunteersForEvent(c, event.ID)
		if err != nil {
			return nil, err
		}
		event.Members = members

		events = append(events, event)
	}

	return &events, nil
}
