package event_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags		ivent
// @Accept		json
// @Produce	json
// @Security Bearer
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/ivents/get-ivents [get]
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
