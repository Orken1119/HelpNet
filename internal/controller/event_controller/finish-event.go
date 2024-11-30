package event_controller

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags        event
// @Summary	exist permission for organization
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Param       id  path      int  true  "Event ID"
// @Success     200 {object} models.SuccessResponse{result=string} ""organization finished event successfully""
// @Failure     400 {object} models.ErrorResponse
// @Failure     500 {object} models.ErrorResponse
// @Router      /events/finish/{id} [put]
func (av *EventController) FinishEvent(c *gin.Context) {
	eventID := c.Param("id")

	id, err := strconv.Atoi(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "INVALID_EVENT_ID",
					Message: "Event ID must be a valid integer.",
				},
			},
		})
		return
	}

	err = av.EventRepository.FinishEvent(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "FINISH_EVENT_ERROR",
					Message: "Error in finishing the event",
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: "organization finished event successfully"})

}
