package event_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags		event
// @Summary	exist permission for volunteer
// @Accept		json
// @Produce	json
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=[]models.Event}
// @Failure	default	{object}	models.ErrorResponse
// @Router		/events/get-events [get]
func (av *EventController) GetAllEvent(c *gin.Context) {
	ivents, err := av.EventRepository.GetAllEvent(c)
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

	c.JSON(http.StatusOK, models.SuccessResponse{Result: ivents})
}
