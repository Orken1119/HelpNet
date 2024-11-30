package event_controller

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags		event
// @Summary	exist permission for organization
// @Description	Update an event's details by its ID
// @Accept		json
// @Produce	json
// @Param		id path int true "id"
// @Param		request body models.EventForEditing true "Event details"
// @Security	Bearer
// @Success     200 {object} models.SuccessResponse{result=models.Event} "Event created successfully"
// @Failure	400 {object} models.ErrorResponse "Bad Request"
// @Failure	404 {object} models.ErrorResponse "Event not found"
// @Router		/events/update-event/{id} [put]
func (av *EventController) UpdateEvent(c *gin.Context) {
	var ivent models.EventForEditing

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

	err = c.ShouldBind(&ivent)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Datas dont match with expected structure",
				},
			},
		})
		return
	}

	err = av.EventRepository.UpdateEvent(c, &ivent, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "IVENT_ERROR",
					Message: "Error to update ivent",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: ivent})

}
