package volunteer_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=[]models.VolunteerProfile}
// @Failure default {object} models.ErrorResponse
// @Router /user/get-all-volunteers-profile [get]
func (sc *UserController) GetAllVolunteersProfile(c *gin.Context) {

	profile, err := sc.UserRepository.GetAllVolunteers(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_ORGANIZATION_PROFILE",
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
