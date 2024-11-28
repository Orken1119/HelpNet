package event_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags		event
// @Accept		json
// @Produce	json
// @Param directioin path string true "direction"
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=[]models.Event}
// @Failure	default	{object}	models.ErrorResponse
// @Router		/events/get-event-by-direction/{direction} [get]
func (av *EventController) GetEventsByDirection(c *gin.Context) {
	direction := c.Param("direction")

	ivent, err := av.EventRepository.GetEventsByDirection(c, direction)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "IVENT_ERROR",
					Message: "Error to get all ivents",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: &ivent})
}
