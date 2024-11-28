package organization

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags organization
// @Accept json
// @Produce json
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=[]models.OrganizationProfile}
// @Failure default {object} models.ErrorResponse
// @Router /organizations/get-all-organizations-profile [get]
func (av *OrganizationController) GetAllOrganizationsProfile(c *gin.Context) {

	profile, err := av.OrganizationRepository.GetAllOrganizations(c)
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
