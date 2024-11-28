package volunteer_controller

import (
	"net/http"
	"strconv"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure default {object} models.ErrorResponse
// @Router /user/delete-certificate/{id} [delete]
func (sc *UserController) DeleteCertificate(c *gin.Context) {
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

	err = sc.UserRepository.DeleteCertificate(c, userID)
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

	c.JSON(http.StatusOK, gin.H{"message": "certificate was deleted"})
}
