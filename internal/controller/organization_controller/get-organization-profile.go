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
// @Success 200 {object} models.SuccessResponse
// @Failure default {object} models.ErrorResponse
// @Router /organizations/profile [get]
func (av *OrganizationController) GetOrganizationProfile(c *gin.Context) {
	orgID := c.GetUint("userID")

	profile, err := av.OrganizationRepository.GetOrganizationProfile(c, int(orgID))
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
