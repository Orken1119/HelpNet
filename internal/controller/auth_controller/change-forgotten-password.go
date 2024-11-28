package auth_controller

import (
	"net/http"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "query params"
// @Success 200 {object} map[string]string
// @Failure default {object} models.ErrorResponse
// @Router /authentication/change-forgotten-password [post]
func (fp *AuthController) ChangeForgottenPassword(c *gin.Context) {
	var request models.ChangePasswordRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Datas dont match with struct of changePassword",
				},
			},
		})
		return
	}

	_, err = fp.UserRepository.GetVolunteerByEmail(c, request.Email)
	if err != nil {
		_, err := fp.UserRepository.GetVolunteerByEmail(c, request.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Result: []models.ErrorDetail{
					{
						Code:    "ERROR_GET_USER",
						Message: "User with this email doesn't found",
					},
				},
			})
			return
		}
	}

	err = ValidatePassword(request.Password.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_PASSWORD_FORMAT",
					Message: err.Error(),
				},
			},
		})
		return
	}

	err = ConfirmPassword(request.Password.Password, request.Password.ConfirmPassword)
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

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_ENCRYPTE_PASSWORD",
					Message: "Couldn't encrypte password",
				},
			},
		})
		return
	}
	request.Password.Password = string(encryptedPassword)

	code, err := fp.UserRepository.GetCodeByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_RECEIVE_CODE",
					Message: "Couldn't receive code",
				},
			},
		})
		return
	}

	err = fp.UserRepository.ChangeForgottenVolunteersPassword(c, code, request.Email, request.Password.Password)
	if err != nil {
		err := fp.UserRepository.ChangeForgottenOrgPassword(c, code, request.Email, request.Password.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Result: []models.ErrorDetail{
					{
						Code:    "ERROR_CHANGE_RASSWORD",
						Message: "Couldn't change password",
						Metadata: models.Properties{
							Properties1: err.Error(),
						},
					},
				},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password was changed"})

}
