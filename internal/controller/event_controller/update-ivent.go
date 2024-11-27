package event_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags		ivent
// @Accept		json
// @Produce	json
// @Param request body models.Event true "query params"
// @Security Bearer
// @Success	200		{object}	models.SuccessResponse
// @Failure	default	{object}	models.ErrorResponse
// @Router		/ivents/update-ivent [put]
func (av *EventController) UpdateEvent(c *gin.Context) {
	var ivent models.Event

	err := c.ShouldBind(&ivent)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Datas dont match with struct of signin",
				},
			},
		})
		return
	}

	err = av.EventRepository.UpdateEvent(c, &ivent)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "IVENT_ERROR",
					Message: "Error to create ivent",
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: ivent})

}
