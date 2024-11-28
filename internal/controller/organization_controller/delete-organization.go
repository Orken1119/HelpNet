package organization

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	OrganizationRepository models.OrganizationRepository
}

// @Tags		organization
// @Accept		json
// @Produce	json
// @Param id path int true "id"
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=string}
// @Failure	default	{object}	models.ErrorResponse
// @Router		/organizations/delete-organizations/{id} [delete]
func (av *OrganizationController) DeleteOrganization(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.Atoi(idVal)
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

	err = av.OrganizationRepository.DeleteOrganization(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "DELETE_ORGANIZATION_ERROR",
					Message: "Error to delete organization",
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: "Organization was deleted successfully"})

}
