package auth_controller

import (
	"net/http"

	"github.com/Orken1119/HelpNet/internal/controller/auth_controller/tokenutil"
	models "github.com/Orken1119/HelpNet/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	UserRepository models.UserRepository
}

// @Accept json
// @Produce json
// @Param request body models.SignUpRequest true "query params"
// @Success 200 {object} models.SuccessResponse
// @Failure default {object} models.ErrorResponse
// @Router /authentication/signup-as-volunteer [post]
func (vc AuthController) SignupAsVolunteer(c *gin.Context) {
	var request models.SignUpRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_BIND_JSON",
					Message: "Data does not match the struct of signup",
				},
			},
		})
		return
	}

	// Check if the user already exists

	user, _ := vc.UserRepository.GetVolunteerByEmail(c, request.Email)
	if user.ID > 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "USER_EXISTS",
					Message: "User with this email already exists",
				},
			},
		})
		return
	}

	// Validate password format
	if err := ValidatePassword(request.Password.Password); err != nil {
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

	// Confirm password match
	if err := ConfirmPassword(request.Password.Password, request.Password.ConfirmPassword); err != nil {
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

	// Encrypt password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_ENCRYPT_PASSWORD",
					Message: "Couldn't encrypt password",
				},
			},
		})
		return
	}
	request.Password.Password = string(encryptedPassword)

	request.Password.ConfirmPassword = ""

	// Create volunteer
	if _, err := vc.UserRepository.CreateUserVolunteer(c, &request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_CREATE_USERS",
					Message: "Couldn't create user",
					Metadata: models.Properties{
						Properties1: err.Error(),
					},
				},
			},
		})
		return
	}

	// Retrieve user to create access token
	user, err = vc.UserRepository.GetVolunteerByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "ERROR_GET_USER",
					Message: "User with this email wasn't found",
				},
			},
		})
		return
	}

	// Create access token
	accessToken, err := tokenutil.CreateAccessToken(user, `access-secret-key`, 50)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Result: []models.ErrorDetail{
				{
					Code:    "TOKEN_ERROR",
					Message: "Error creating access token",
				},
			},
		})
		return
	}

	// Return success response with the access token
	c.JSON(http.StatusOK, models.SuccessResponse{Result: accessToken})
}
