package event_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

func (av *EventController) GetOrgEvents(c *gin.Context) {
	orgID := c.GetUint("orgID")

	events, err := av.EventRepository.GetOrganizationsInProcessEvents(c, int(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "FETCH_EVENTS_ERROR",
					Message: "Error fetching in-process events for the organization.",
				},
			},
		})
		return
	}

	if events == nil || len(*events) == 0 {
		c.JSON(http.StatusOK, models.SuccessResponse{
			Result: "No in-process events found for the organization.",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Result: events,
	})
}
