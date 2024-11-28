package volunteer_controller

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserRepository models.UserRepository
}

// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=models.VolunteerProfile}
// @Failure default {object} models.ErrorResponse
// @Router /user/profile/{id} [get]
func (sc *UserController) GetProfile(c *gin.Context) {
	idVal := c.Param("id")
	userID, err := strconv.Atoi(idVal)
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

	profile, err := sc.UserRepository.GetVolunteerProfile(c, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_USER_PROFILE",
					Message: "Can't get profile from db",
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
