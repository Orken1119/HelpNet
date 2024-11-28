package volunteer_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags user
// @Summary	exist permission for organization and volunteers
// @Accept json
// @Produce json
// @Param name path int true "name"
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=[]models.Event}
// @Failure default {object} models.ErrorResponse
// @Router /user/search-event/{name} [get]
func (sc *UserController) SearchEvent(c *gin.Context) {
	name := c.Param("name")

	profile, err := sc.UserRepository.SearchEvent(c, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_SEARCH_EVENT",
					Message: "Can't get event from db",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Result: profile})
}
