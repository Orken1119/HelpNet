package event_controller

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags		event
// @Accept		json
// @Produce	json
// @Param id path int true "id"
// @Security Bearer
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/events/get-event-by-id/{id} [get]
func (av *EventController) GetEventById(c *gin.Context) {
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

	ivent, err := av.EventRepository.GetEventById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "IVENT_ERROR",
					Message: "Error to get all ivents",
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: ivent})
}
