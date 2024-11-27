package event_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags        ivent
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Success     200 {object} models.SuccessResponse
// @Failure     400 {object} models.ErrorResponse
// @Failure     500 {object} models.ErrorResponse
// @Router      /ivents/finished-events-by-organization [get]
func (av *EventController) GetOrgFinishedEvents(c *gin.Context) {
	orgID := c.GetUint("orgID")

	events, err := av.EventRepository.GetFinishedEventsByOrganization(c, int(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "FETCH_EVENTS_ERROR",
					Message: "Error fetching finished events for the organization.",
				},
			},
		})
		return
	}

	if events == nil || len(*events) == 0 {

		c.JSON(http.StatusOK, models.SuccessResponse{
			Result: "No finished events found for the organization.",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Result: events,
	})
}
