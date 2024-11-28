package event_controller

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags        event
// @Accept      json
// @Produce     json
// @Param id path int true "id"
// @Security    Bearer
// @Success     200 {object} models.SuccessResponse
// @Failure     400 {object} models.ErrorResponse
// @Failure     500 {object} models.ErrorResponse
// @Router      /events/participate-event/{id} [post]
func (av *EventController) JoinEvent(c *gin.Context) {

	userID := c.GetUint("userID")
	idVal := c.Param("id")
	id, err := strconv.Atoi(idVal)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_IVENT",
					Message: "incorrect id format",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	err = av.EventRepository.ParticipateEvent(c, int(userID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "PARTICIPATE_EVENT_ERROR",
					Message: "Error adding the user to the event.",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: "you joined successfully"})
}
