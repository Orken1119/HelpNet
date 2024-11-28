package organization

import (
	"net/http"

	"github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags		organization
// @Summary	exist permission for organization
// @Accept		json
// @Produce	json
// @Param request body models.OrganizationProfileEditing true "query params"
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=models.OrganizationProfileEditing}
// @Failure	default	{object}	models.ErrorResponse
// @Router		/organizations/edit-organization-profile [put]
func (av *OrganizationController) EditOrganization(c *gin.Context) {
	var org models.OrganizationProfileEditing

	orgID := c.GetUint("userID")

	err := c.ShouldBind(&org)
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

	err = av.OrganizationRepository.EditOrganizationProfile(c, int(orgID), &org)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "IVENT_ERROR",
					Message: "Error to create ivent",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: org})

}
