package event_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	EventRepository models.EventRepository
}

// @Tags		ivent
// @Accept		json
// @Produce	json
// @Param request body models.Event true "query params"
// @Security Bearer
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/ivents/create-ivent [post]
func (av *EventController) CreateEvent(c *gin.Context) {
	var ivent models.Event

	id := c.GetUint("orgID")

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

	readyivent, err := av.EventRepository.CreateEvent(c, &ivent, int(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "IVENT_ERROR",
					Message: "Error creating event",
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: readyivent})
}
