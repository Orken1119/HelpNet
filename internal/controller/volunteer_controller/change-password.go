package volunteer_controller

import (
	"fmt"
	"net/http"
	"unicode"

	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Tags user
// @Summary	exist permission for volunteer
// @Accept json
// @Produce json
// @Param request body models.Password true "query params"
// @Security Bearer
// @Success     200 {object} map[string]string "Password successfully changed"
// @Failure default {object} models.ErrorResponse
// @Router /user/change-password [put]
func (sc *UserController) ChangePassword(c *gin.Context) {
	var request models.Password

	userID := c.GetUint("userID")

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

	err = ValidatePassword(request.Password)
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
	// Подтверждение пароля
	err = ConfirmPassword(request.Password, request.ConfirmPassword)
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
	//

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
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
	request.Password = string(encryptedPassword)

	err = sc.UserRepository.ChangePassword(c, int(userID), request.Password)
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

	c.JSON(http.StatusOK, gin.H{"message": "Password was changed successfully"})

}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var (
		hasUpper, hasLower, hasDigit bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter and one digit")
	}

	return nil
}

// Функция подтверждения
func ConfirmPassword(password string, confirm string) error {
	if password != confirm {
		return fmt.Errorf("doesn't match the password")
	}
	return nil
}
