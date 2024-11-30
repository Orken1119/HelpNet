package event_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	EventRepository models.EventRepository
}

// @Tags		event
// @Summary	exist permission for organization
// @Accept		json
// @Produce	json
// @Param request body models.EventForCreating true "query params"
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=models.Event}
// @Failure	default	{object}	models.ErrorResponse
// @Router		/events/create-event [post]
func (av *EventController) CreateEvent(c *gin.Context) {
	var ivent models.EventForCreating

	err := c.ShouldBind(&ivent)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Data doesn't match the expected structure",
				},
			},
		})
		return
	}

	readyivent, err := av.EventRepository.CreateEvent(c, &ivent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "IVENT_ERROR",
					Message: "Error on creating event",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: readyivent})
}
