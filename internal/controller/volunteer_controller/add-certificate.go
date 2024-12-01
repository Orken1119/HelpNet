package volunteer_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
)

// @Tags user
// @Summary	exist permission for volunteer
// @Accept json
// @Produce json
// @Param request body models.AddingSertificate true "query params"
// @Security Bearer
// @Success     200 {object} models.SuccessResponse{result=string}
// @Failure default {object} models.ErrorResponse
// @Router /user/add-certificate [post]
func (sc *UserController) AddCertificate(c *gin.Context) {
	var request models.AddingSertificate
	userID := c.GetUint("userID")

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Datas dont match with struct of certificate",
				},
			},
		})
		return
	}

	err = sc.UserRepository.AddCertificate(c, request.ImageUrl, int(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_PASSWORD_MISMATCH",
					Message: err.Error(),
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certificate was added"})
}
