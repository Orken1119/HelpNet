package organization

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags organization
// @Summary	exist permission for organization and for volunteer
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=models.OrganizationProfile}
// @Failure default {object} models.ErrorResponse
// @Router /organizations/profile/{id} [get]
func (av *OrganizationController) GetOrganizationProfile(c *gin.Context) {
	idVal := c.Param("id")
	orgID, err := strconv.Atoi(idVal)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_ORGANIZATION",
					Message: "incorrect id format",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	profile, err := av.OrganizationRepository.GetOrganizationProfile(c, orgID)
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
