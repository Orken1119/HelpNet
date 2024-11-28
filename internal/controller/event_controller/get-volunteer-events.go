package event_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags        event
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Success     200 {object} models.SuccessResponse
// @Failure     400 {object} models.ErrorResponse
// @Failure     500 {object} models.ErrorResponse
// @Router      /events/get-user-participating [get]
func (av *EventController) GetVolEvents(c *gin.Context) {
	userID := c.GetUint("userID")

	events, err := av.EventRepository.GetVolunteerParticipatingEvents(c, int(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "FETCH_EVENTS_ERROR",
					Message: "Error fetching participating events for the volunteer.",
				},
			},
		})
		return
	}

	if events == nil || len(*events) == 0 {
		c.JSON(http.StatusOK, models.SuccessResponse{
			Result: "you didn't join in any event",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Result: events,
	})
}
