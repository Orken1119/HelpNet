package auth_controller

import (
	"net/http"

	"github.com/Orken1119/HelpNet/internal/controller/auth_controller/tokenutil"
	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Accept json
// @Produce json
// @Param request body models.SigninRequest true "body json"
// @Success 200 {object} models.SuccessResponse
// @Failure default {object} models.ErrorResponse
// @Router /signin-as-organization [post]
func (lc *AuthController) Signin(c *gin.Context) {
	var loginRequest models.SigninRequest

	err := c.ShouldBind(&loginRequest)
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

	if loginRequest.Email == "" || loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "EMPTY_VALUES",
					Message: "Empty values are written in the form",
				},
			},
		})
		return
	}

	org, err := lc.UserRepository.GetOrganizationByEmail(c, loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_ORGANIZATION",
					Message: "Organization with this email wasn't found",
				},
			},
		})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(org.Password), []byte(loginRequest.Password)) != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "PASSWORD_INCORRECT",
					Message: "Password doesn't match",
				},
			},
		})
		return
	}
	accessToken, err := tokenutil.CreateAccessToken(org, `access-secret-key`, 50) //
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "TOKEN_ERROR",
					Message: "Error to create access token",
				},
			},
		})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Result: accessToken})
}
