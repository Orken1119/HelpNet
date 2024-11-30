package organization

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags        organization
// @Summary     permission for organization
// @Accept      json
// @Produce     json
// @Param       memberID path int true "Member ID"
// @Param       eventID path int true "Event ID"
// @Security    Bearer
// @Success     200 {object} models.SuccessResponse{result=string} "Member successfully removed"
// @Failure     400 {object} models.ErrorResponse "Invalid ID format"
// @Failure     500 {object} models.ErrorResponse "Server error"
// @Router      /organizations/delete-member/{memberID}/event-id/{eventID} [delete]
func (av *OrganizationController) DeleteMember(c *gin.Context) {
	idVal := c.Param("memberID")
	memberID, err := strconv.Atoi(idVal)
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

	eventIdVal := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIdVal)
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

	err = av.OrganizationRepository.DeleteMemberFromEvent(c, memberID, eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "DELETE_ORGANIZATION_ERROR",
					Message: "Error to delete member",
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Result: "Member was deleted successfully"})

}
